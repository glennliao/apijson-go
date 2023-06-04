package query

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

type queryNode struct {
	node *Node
}

func newQueryNode(n *Node) *queryNode {
	return &queryNode{node: n}
}

func (q *queryNode) parse() {
	n := q.node

	accessConfig, err := n.queryContext.queryConfig.GetAccessConfig(n.Key, n.queryContext.NoAccessVerify)
	if err != nil {
		n.err = err
		return
	}

	n.executorConfig = config.NewExecutorConfig(accessConfig, http.MethodGet, n.queryContext.NoAccessVerify)
	n.executorConfig.DbFieldStyle = n.queryContext.DbFieldStyle
	n.executorConfig.JsonFieldStyle = n.queryContext.JsonFieldStyle
	n.executorConfig.DBMeta = n.queryContext.DbMeta

	if n.isList {
		fieldsGet := n.executorConfig.GetFieldsGetByRole()
		if *fieldsGet.MaxCount != 0 {
			if n.page.Count > *fieldsGet.MaxCount {
				n.err = gerror.New(" > maxCount: " + n.Path)
				return
			}
		}
	}

	var accessWhereCondition model.MapStrAny

	setNodeRole(n, n.Key, n.role)
	n.executorConfig.SetRole(n.role)

	if n.role == consts.DENY {
		n.err = gerror.New("deny node: " + n.Path)
		return
	}

	if n.queryContext.NoAccessVerify == false {
		has, condition, err := hasAccess(n)
		if err != nil {
			n.err = err
			return
		}

		if !has {
			n.err = gerror.New("无权限访问:" + n.Key + " by " + n.role)
			return
		}

		accessWhereCondition = condition.Where()
	}

	queryExecutor, err := executor.NewQueryExecutor(n.executorConfig.Executor(), n.ctx, n.executorConfig)
	if err != nil {
		n.err = err
		return
	}

	n.executor = queryExecutor

	// 查询条件
	refKeyMap, conditionMap, ctrlMap := parseQueryNodeReq(n.req, n.isList)

	n.executor.ParseCtrl(ctrlMap)

	if v, exists := ctrlMap["@column"]; exists {
		var exp = regexp.MustCompile(`^[\s\w][\w()]+`) // 匹配 field, COUNT(field)

		fieldStr := strings.ReplaceAll(gconv.String(v), ";", ",")

		fieldList := strings.Split(fieldStr, ",")

		for i, item := range fieldList {
			fieldList[i] = exp.ReplaceAllStringFunc(item, func(field string) string {
				return field
			})
		}

		for _, item := range fieldList {
			if strings.Contains(item, ":") {
				n.Column[item] = item
			} else {
				n.Column[item] = item
			}
		}

		fmt.Println(fieldList)
	}

	err = n.executor.ParseCondition(conditionMap, true)
	if err != nil {
		n.err = err
		return
	}

	err = n.executor.ParseCondition(accessWhereCondition, false)
	if err != nil {
		n.err = err
		return
	}

	n.primaryTableKey = filepath.Base(n.Path)

	if len(refKeyMap) > 0 { // 需要引用别处
		n.refKeyMap = make(map[string]NodeRef)
		hasRefBrother := false // 是否引用兄弟节点, 列表中的主表不能依赖兄弟节点

		for refKey, refStr := range refKeyMap {
			if strings.HasPrefix(refStr, "/") { // 这里/开头是相对同级
				refStr = filepath.Dir(n.Path) + refStr
			}

			refPath, refCol := util.ParseRefCol(refStr)

			if !hasRefBrother {
				if filepath.Dir(n.Path) == filepath.Dir(refPath) {
					hasRefBrother = true
				}
			}

			if refPath == n.Path { // 不能依赖自身
				n.err = gerror.Newf("node cannot ref self: (%s) {%s:%s}", refPath, refKey, refStr)
				return
			}

			refNode := n.queryContext.pathNodes[refPath]
			if refNode == nil {
				n.err = gerror.Newf(" node %s is nil, but ref by %s", refPath, n.Path)
				return
			}

			if refNode.err != nil {
				n.err = refNode.err
				return
			}

			for _, _refN := range refNode.refKeyMap {
				if _refN.node.Path == n.Path {
					n.err = gerror.Newf("circle ref %s & %s", refNode.Path, n.Path)
					return
				}
			}

			n.refKeyMap[refKey] = NodeRef{
				column: refCol,
				node:   refNode,
			}

		}

		if hasRefBrother {
			n.primaryTableKey = ""
		}
	}
}

func (q *queryNode) fetch() {
	n := q.node
	for refK, refNode := range n.refKeyMap {
		ret, err := refNode.node.Result()
		if err != nil {
			n.err = err
			return
		}

		if refNode.node.isList {
			list := ret.([]model.Map)

			valList := getColList(list, refNode.column)
			if len(valList) == 0 { // 未查询到主表, 故当前不再查询
				n.executor.SetEmptyResult()
				break
			}

			err = n.executor.ParseCondition(model.MapStrAny{
				refK + consts.OpIn: valList, //  @ 与 {}&等的结合 id{}@的处理
			}, false)

			if err != nil {
				n.err = err
				return
			}

		} else {

			if ret == nil { // 未查询到主表, 故当前不再查询
				n.executor.SetEmptyResult()
				break
			}

			item := ret.(model.Map)

			refVal := item[refNode.column]

			var refConditionMap = model.MapStrAny{
				refK: refVal,
			}
			err = n.executor.ParseCondition(refConditionMap, false)
			if err != nil {
				n.err = err
				return
			}
		}
	}

	if n.isList {

		page := n.page.Page
		count := n.page.Count

		if n.primaryTableKey == "" { // 主查询表 才分页
			page = 0
			count = 0
		}

		n.total, n.err = n.executor.Count()
		if n.total > 0 {
			n.ret, n.err = n.executor.List(page, count)
		}

	} else {
		n.ret, n.err = n.executor.One()
	}
	if n.err != nil {
		return
	}

	queryConfig := n.queryContext.queryConfig

	// 需优化调整
	for k, v := range n.req {
		if !strings.HasSuffix(k, consts.FunctionsKeySuffix) {
			continue
		}

		k = k[0 : len(k)-2]

		functionName, paramKeys := util.ParseFunctionsStr(v.(string))
		_func := queryConfig.Func(functionName)

		if n.isList {
			for i, item := range n.ret.([]model.Map) {

				param := model.Map{}
				for paramI, paramItem := range _func.ParamList {
					if paramItem.Name == consts.FunctionOriReqParam {
						param[paramItem.Name] = item
					} else {
						param[paramItem.Name] = item[paramKeys[paramI]]
					}
				}

				val, err := _func.Handler(n.ctx, param)
				if err != nil {
					n.err = err
					return
				}
				n.ret.([]model.Map)[i][k] = val
			}
		} else {
			param := model.Map{}
			for paramI, paramItem := range _func.ParamList {
				if paramItem.Name == consts.FunctionOriReqParam {
					param[paramItem.Name] = n.ret.(model.Map)
				} else {
					param[paramItem.Name] = n.ret.(model.Map)[paramKeys[paramI]]
				}
			}

			val, err := _func.Handler(n.ctx, param)
			if err != nil {
				n.err = err
				return
			}
			n.ret.(model.Map)[k] = val
		}
	}

}

func (q *queryNode) result() {
	n := q.node
	if n.isList {
		if n.ret == nil || n.ret.([]model.Map) == nil {
			n.ret = []model.Map{}
		}
	} else {
		if n.ret == nil || n.ret.(model.Map) == nil {
			n.ret = nil
		}
	}
}

func (q *queryNode) nodeType() int {
	return NodeTypeQuery
}
