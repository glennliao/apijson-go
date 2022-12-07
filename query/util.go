package query

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/db"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"path/filepath"
	"strings"
)

func isFirstUp(str string) bool {
	firstLetter := str[0]
	return firstLetter >= 'A' && firstLetter <= 'Z'
}

// hasFirstUpKey 用户判断是否存在查询节点
func hasFirstUpKey(m g.Map) bool {
	for k := range m {
		if isFirstUp(k) {
			return true
		}
	}
	return false
}

// parseTableKey 解析表名 //todo 增加一个通用解析key的方法, 避免到处出现判断key后缀以及截取操作
func parseTableKey(k string, p string) (tableName string) {
	tableName = k

	if strings.HasSuffix(k, "[]") {
		tableName = k[0 : len(k)-2]
	} else if strings.Contains(p, "[]") {
		tableName = k
	}
	return tableName
}

// parseQueryNodeReq 解析节点请求内容
func parseQueryNodeReq(reqMap g.Map, isList bool) (refMap g.MapStrStr, where g.Map, ctrlMap g.Map) {
	refMap = g.MapStrStr{}
	ctrlMap = g.Map{}
	where = g.Map{}
	for k, v := range reqMap {

		if strings.HasSuffix(k, consts.FunctionsKeySuffix) {
			continue
		}

		if strings.HasSuffix(k, "@") { //引用
			refMap[k[0:len(k)-1]] = gconv.String(v)
		} else if strings.HasPrefix(k, "@") { // @column等ctrl字段
			ctrlMap[k] = v
		} else {
			if isList {
				switch k {
				case "page", "count", "query": // todo 调整常量
					// 分页字段不传递到sqlExecutor
					continue
				}
			}

			where[k] = v
		}
	}
	return
}

func parseRefCol(refStr string) (refPath string, refCol string) {
	// "id@":"[]/User/userId"
	refCol = filepath.Base(refStr)                  // userId
	refPath = refStr[0 : len(refStr)-len(refCol)-1] // []/User
	return refPath, refCol
}

func hasAccess(node *Node, table string) (hasAccess bool, accessWhere g.Map, err error) {
	accessRoles, tableName, err := db.GetAccessRole(table, consts.MethodGet)
	if err != nil {
		return false, nil, err
	}

	if !lo.Contains(accessRoles, node.role) {
		g.Log().Debug(node.ctx, table, "role:", node.role, "accessRole", accessRoles, " -> deny")
		return false, nil, err
	}

	accessWhere, err = node.queryContext.AccessCondition(node.ctx, config.AccessConditionReq{
		Table:               tableName,
		TableAccessRoleList: accessRoles,
		Method:              consts.MethodGet,
		NodeReq:             node.req,
		NodeRole:            node.role,
	})

	return true, accessWhere, err

}

func getColList(list []g.Map, col string) []any {

	set := gset.New()
	for _, item := range list {
		set.Add(gconv.String(item[col]))
	}
	return set.Slice()
}

func setNeedTotal(node *Node) {
	node.needTotal = true
	if node.Type == NodeTypeStruct {
		setNeedTotal(node.children[node.primaryTableKey])
	}
}

// setNodeRole 设置节点的@role, 根据 config.DefaultRoleFunc 获取节点最终的@role
func setNodeRole(node *Node, tableName string, parenNodeRole string) {

	role, ok := node.req[consts.Role]

	if node.Type != NodeTypeQuery {
		if !ok {
			node.role = parenNodeRole
		} else {
			node.role = gconv.String(role)
		}
	} else {
		if ok {
			node.role, _ = config.DefaultRoleFunc(node.ctx, config.RoleReq{
				Table:    tableName,
				NodeRole: gconv.String(role),
			})
		} else {
			node.role, _ = config.DefaultRoleFunc(node.ctx, config.RoleReq{
				Table:    tableName,
				NodeRole: parenNodeRole,
			})
		}
	}
}

// analysisRef 分析依赖, 将依赖关系保存到prerequisites中
func analysisRef(p *Node, prerequisites *[][]string) {

	// 分析依赖关系, 让无依赖的先执行， 然后在执行后续的
	for _, node := range p.children {
		for _, refNode := range node.refKeyMap {
			*prerequisites = append(*prerequisites, []string{node.Path, refNode.node.Path})
		}
		analysisRef(node, prerequisites)
	}
}

// analysisOrder 使用拓扑排序 分析节点fetch优先级
func analysisOrder(prerequisites [][]string) ([]string, error) {

	var pointMap = make(map[string]bool)
	for _, prerequisite := range prerequisites {
		pointMap[prerequisite[0]] = true
		pointMap[prerequisite[1]] = true
	}

	var pointNum = len(pointMap)
	var edgesMap = make(map[string][]string)
	var inDeg = make(map[string]int)
	var result []string

	for _, prerequisite := range prerequisites {
		edgesMap[prerequisite[1]] = append(edgesMap[prerequisite[1]], prerequisite[0])
		inDeg[prerequisite[0]]++
	}

	var queue []string

	for point, _ := range pointMap {
		if inDeg[point] == 0 {
			queue = append(queue, point)
		}
	}

	for len(queue) > 0 {
		var first string
		first, queue = queue[0], queue[1:]
		result = append(result, first)
		for _, point := range edgesMap[first] {
			inDeg[point]--
			if inDeg[point] == 0 {
				queue = append(queue, point)
			}
		}
	}

	if len(result) != pointNum {
		return nil, gerror.New("依赖循环, 请检查请求")
	}

	return result, nil

}
