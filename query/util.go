package query

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

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
func parseQueryNodeReq(reqMap model.Map, isList bool) (refMap model.MapStrStr, where model.MapStrAny, ctrlMap model.Map) {
	refMap = model.MapStrStr{}
	ctrlMap = model.Map{}
	where = model.MapStrAny{}
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

func hasAccess(node *Node) (hasAccess bool, accessWhere *config.ConditionRet, err error) {
	accessRoles := node.executorConfig.AccessRoles()
	if err != nil {
		return false, nil, err
	}

	if !lo.Contains(accessRoles, node.role) {
		g.Log().Debug(node.ctx, node.Key, "role:", node.role, "accessRole", accessRoles, " -> deny")
		return false, nil, err
	}

	accessWhere, err = node.queryContext.AccessCondition(node.ctx, config.ConditionReq{
		AccessName:          node.Key,
		TableAccessRoleList: accessRoles,
		Method:              http.MethodGet,
		NodeReq:             node.req,
		NodeRole:            node.role,
	})

	return true, accessWhere, err

}

func getColList(list []model.Map, col string) []any {

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
			node.role, _ = node.queryContext.queryConfig.DefaultRoleFunc()(node.ctx, config.RoleReq{
				AccessName: tableName,
				NodeRole:   gconv.String(role),
			})
		} else {
			node.role, _ = node.queryContext.queryConfig.DefaultRoleFunc()(node.ctx, config.RoleReq{
				AccessName: tableName,
				NodeRole:   parenNodeRole,
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
