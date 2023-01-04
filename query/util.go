package query

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/db"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
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

// parseTableKey 解析表名
func parseTableKey(k string, p string) (tableName string) {
	tableName = k

	if strings.HasSuffix(k, consts.ListKeySuffix) {
		tableName = k[0 : len(k)-len(consts.ListKeySuffix)]
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

		if strings.HasSuffix(k, consts.RefKeySuffix) { //引用
			refMap[k[0:len(k)-1]] = gconv.String(v)
		} else if strings.HasPrefix(k, "@") { // @column等ctrl字段
			ctrlMap[k] = v
		} else {
			if isList {
				switch k {
				case consts.Page, consts.Count, consts.Query:
					// 分页字段不传递到sqlExecutor
					continue
				}
			}

			where[k] = v
		}
	}
	return
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
