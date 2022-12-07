package query

import (
	"context"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/db"
	"github.com/glennliao/apijson-go/functions"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"path/filepath"
	"strings"
	"time"
)

const (
	NodeTypeStruct = iota // 结构节点
	NodeTypeQuery         // 查询节点
	NodeTypeRef           // 引用节点
	NodeTypeFunc          // functions 节点
)

type Node struct {
	ctx          context.Context
	queryContext *Query

	// 当前节点key Todos
	Key string
	// 当前节点path -> []/Todos
	Path string
	// 节点类型
	Type int

	// 是否为列表节点
	isList bool
	page   g.Map // 分页参数

	// 访问当前节点的角色
	role string

	// 节点的请求数据
	req          g.Map
	simpleReqVal string //非对象结构

	// 节点数据执行器
	sqlExecutor *db.SqlExecutor

	startAt time.Time
	endAt   time.Time

	// 执行完毕
	finish bool

	ret any
	err error

	children map[string]*Node

	refKeyMap map[string]NodeRef // 关联字段

	primaryTableKey string // 主查询表

	total     int64 // 数据总条数
	needTotal bool
}

// NodeRef 节点依赖引用
type NodeRef struct {
	column string
	node   *Node
}

/**
node 生命周期
new -> buildChild -> parse -> fetch -> result
*/

func newNode(query *Query, key string, path string, nodeReq any) *Node {

	g.Log().Debugf(query.ctx, "【node】(%s) <new> ", path)

	node := &Node{
		ctx:          query.ctx,
		queryContext: query,
		Key:          key,
		Path:         path,
		startAt:      time.Now(),
		finish:       false,
	}

	// 节点类型判断
	if key != "" {
		if isFirstUp(key) { // 大写开头, 为查询节点(对应数据库)
			node.Type = NodeTypeQuery
		} else if strings.HasSuffix(key, "@") {
			node.Type = NodeTypeRef
		} else if strings.HasSuffix(key, consts.FunctionsKeySuffix) {
			node.Type = NodeTypeFunc
		} else {
			node.Type = NodeTypeStruct
			// 结构节点下应该必须存在查询节点
		}

		if strings.HasSuffix(key, consts.ListKeySuffix) || strings.HasSuffix(filepath.Dir(path), consts.ListKeySuffix) {
			node.isList = true
		}
	}

	if req, ok := nodeReq.(g.Map); ok {
		node.req = req

	} else {
		node.simpleReqVal = gconv.String(nodeReq)
	}

	return node
}

func (n *Node) buildChild() error {

	if n.Type == NodeTypeQuery && !hasFirstUpKey(n.req) { // 查询节点嵌套查询节点, 目前不支持
		return nil
	}

	// 最大深度检查
	if len(strings.Split(n.Path, "/")) > consts.MaxTreeDeep {
		return gerror.Newf("deep(%s) > %d", n.Path, consts.MaxTreeDeep)
	}

	children := make(map[string]*Node)

	for key, v := range n.req {

		if strings.HasPrefix(key, "@") {
			continue
		}

		if n.Type == NodeTypeQuery && !isFirstUp(key) { // 查询节点嵌套查询节点, 目前不支持
			continue
		}

		if n.isList {
			if lo.Contains([]string{"total", "page"}, key) {
				continue
			}
		}

		path := n.Path
		if path != "" { // 根节点时不带/
			path += "/"
		}
		node := newNode(n.queryContext, key, path+key, v)

		if n.Type != NodeTypeQuery { // 非查询节点role主要的功能是传递角色(设置该节点下子节点的角色)
			setNodeRole(node, "", n.role)
		}

		err := node.buildChild()
		if err != nil {
			return err
		}
		children[key] = node
	}

	if len(children) > 0 {

		// 最大宽度检查, 目前为某节点的宽度, 应该计算为整棵树的最大宽度
		if len(children) > consts.MaxTreeWidth {
			path := n.Path
			if path == "" {
				path = "root"
			}
			return gerror.Newf("width(%s) > %d", path, consts.MaxTreeWidth)
		}

		n.children = children

		for _, node := range children {
			n.queryContext.pathNodes[node.Path] = node
		}
	}

	return nil
}

func (n *Node) parse() {

	g.Log().Debugf(n.ctx, "【node】(%s) <parse> ", n.Path)

	switch n.Type {
	case NodeTypeQuery:
		tableKey := parseTableKey(n.Key, n.Path)

		access, err := db.GetAccess(tableKey, n.queryContext.AccessVerify)
		if err != nil {
			n.err = err
			return
		}

		var accessWhereCondition g.Map

		setNodeRole(n, access.Name, n.role)

		if n.role == consts.DENY {
			n.err = gerror.New("deny node: " + n.Path)
			return
		}

		if n.queryContext.AccessVerify {
			has, condition, err := hasAccess(n, tableKey)
			if err != nil {
				n.err = err
				return
			}

			if !has {
				n.err = gerror.New("无权限访问:" + tableKey + " by " + n.role)
				return
			}

			accessWhereCondition = condition
		}

		executor, err := db.NewSqlExecutor(n.ctx, n.queryContext.AccessVerify, n.role, access)
		if err != nil {
			n.err = err
			return
		}

		n.sqlExecutor = executor

		// 查询条件
		refKeyMap, conditionMap, ctrlMap := parseQueryNodeReq(n.req, n.isList)

		n.sqlExecutor.ParseCtrl(ctrlMap)

		err = n.sqlExecutor.ParseCondition(conditionMap, true)
		if err != nil {
			n.err = err
			return
		}

		err = n.sqlExecutor.ParseCondition(accessWhereCondition, false)
		if err != nil {
			n.err = err
			return
		}

		n.primaryTableKey = n.Key

		if len(refKeyMap) > 0 { // 需要应用别处
			n.refKeyMap = make(map[string]NodeRef)
			hasRefBrother := false // 是否引用兄弟节点, 列表中的主表不能依赖兄弟节点

			for refKey, refStr := range refKeyMap {
				if strings.HasPrefix(refStr, "/") { // 这里/开头是相对同级
					refStr = filepath.Dir(n.Path) + refStr
				}

				refPath, refCol := parseRefCol(refStr)

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

	case NodeTypeRef:

		refStr := n.simpleReqVal
		if strings.HasPrefix(refStr, "/") { // 这里/开头是相对同级
			refStr = filepath.Dir(n.Path) + refStr
		}
		refPath, refCol := parseRefCol(refStr)
		if refPath == n.Path { // 不能依赖自身
			panic(gerror.Newf("node cannot ref self: (%s)", refPath))
		}

		refNode := n.queryContext.pathNodes[refPath]
		if refNode == nil {
			panic(gerror.Newf(" node %s is nil, but ref by %s", refPath, n.Path))
		}

		n.refKeyMap = make(map[string]NodeRef)

		if strings.HasSuffix(n.simpleReqVal, "[]/total") {
			setNeedTotal(refNode)
		}

		n.refKeyMap[n.Key] = NodeRef{
			column: refCol,
			node:   refNode,
		}

	case NodeTypeStruct:

		for _, childNode := range n.children {
			childNode.parse()
		}

		if n.isList { // []节点

			page := g.Map{}
			if v, exists := n.req["page"]; exists {
				page["page"] = gconv.Int(v)
			}
			if v, exists := n.req["count"]; exists {
				page["count"] = gconv.Int(v)
			}

			hasPrimary := false // 是否存在主查询表
			for _, child := range n.children {

				if child.err != nil {
					n.err = child.err
					return
				}

				if child.primaryTableKey != "" {

					if hasPrimary {
						panic(gerror.Newf("node must only one primary table: (%s)", n.Path))
					}

					hasPrimary = true
					n.primaryTableKey = child.Key
					child.page = page

				}
			}

			if n.Key == "[]" && !hasPrimary {
				panic(gerror.Newf("node must have  primary table: (%s)", n.Path))
			}
		}

	case NodeTypeFunc:
		functionName, paramKeys := functions.ParseFunctionsStr(n.simpleReqVal)
		n.simpleReqVal = functionName
		for _, key := range paramKeys {
			g.Dump(key)
		}
	}
	g.Log().Debugf(n.ctx, "【node】(%s) <parse-endAt> ", n.Path)

}

func (n *Node) fetch() {

	defer func() {
		n.finish = true
		n.endAt = time.Now()
		g.Log().Debugf(n.ctx, "【node】(%s) <fetch-endAt> ", n.Path)
	}()

	g.Log().Debugf(n.ctx, "【node】(%s) <fetch> hasFinish: 【%v】", n.Path, n.finish)

	if n.finish {
		g.Log().Error(n.ctx, "再次执行", n.Path)
		return
	}

	if n.err != nil {
		return
	}

	switch n.Type {
	case NodeTypeQuery:

		for refK, refNode := range n.refKeyMap {
			ret, err := refNode.node.Result()
			if err != nil {
				//g.Log().Error(n.ctx, "", err)
				n.err = err
				return
			}

			if refNode.node.isList {
				list := ret.([]g.Map)

				valList := getColList(list, refNode.column)
				if len(valList) == 0 { // 未查询到主表, 故当前不再查询
					n.sqlExecutor.WithEmptyResult = true
					break
				}

				err := n.sqlExecutor.ParseCondition(g.Map{
					refK + "{}": valList, //  @ 与 {}&等的结合 id{}@的处理
				}, false)

				if err != nil {
					n.err = err
					return
				}

			} else {

				if ret == nil { // 未查询到主表, 故当前不再查询
					n.sqlExecutor.WithEmptyResult = true
					break
				}

				item := ret.(g.Map)

				refVal := item[refNode.column]

				var refConditionMap = g.Map{
					refK: refVal,
				}
				err := n.sqlExecutor.ParseCondition(refConditionMap, false)
				if err != nil {
					n.err = err
					return
				}
			}
		}

		if n.isList {

			page := 1
			count := 10

			for k, v := range n.page {
				switch k {
				case "page":
					page = gconv.Int(v)

				case "count":
					count = gconv.Int(v)
				}
			}

			for k, v := range n.req {
				switch k {
				case "page":
					page = gconv.Int(v)

				case "count":
					count = gconv.Int(v)
				case "query":
					switch gconv.String(v) {
					case "1", "2":
						n.needTotal = true
					}
				}
			}

			if n.primaryTableKey == "" { // 主查询表 才分页
				page = 0
				count = 0
			}

			n.ret, n.total, n.err = n.sqlExecutor.List(page, count, n.needTotal)
		} else {
			n.ret, n.err = n.sqlExecutor.One()
		}
		if n.err != nil {
			return
		}

		// 需优化调整
		for k, v := range n.req {
			if !strings.HasSuffix(k, consts.FunctionsKeySuffix) {
				continue
			}

			k := k[0 : len(k)-2]

			functionName, paramKeys := functions.ParseFunctionsStr(v.(string))

			if n.isList {
				for i, item := range n.ret.([]g.Map) {
					var param = g.Map{}
					for _, key := range paramKeys {
						if key == "$req" {
							param[key] = item
						} else {
							param[key] = item[key]
						}
					}
					var err error
					n.ret.([]g.Map)[i][k], err = functions.Call(n.ctx, functionName, param)
					if err != nil {
						panic(err) // todo
					}
				}
			} else {
				var param = g.Map{}
				for _, key := range paramKeys {
					if key == "$req" {
						param[key] = n.ret.(g.Map)
					} else {
						param[key] = n.ret.(g.Map)[key]
					}

				}
				var err error
				n.ret.(g.Map)[k], err = functions.Call(n.ctx, functionName, param)
				if err != nil {
					panic(err) // todo
				}
			}
		}

	case NodeTypeRef:
		for _, refNode := range n.refKeyMap {
			if strings.HasSuffix(refNode.column, "total") && strings.HasSuffix(refNode.node.Path, "[]") {
				n.total = refNode.node.total
			}
		}

	case NodeTypeStruct:
		// 目前结构节点组装数据在result, 如果被依赖的是组装后的, 则无法查询。 如遇到此情况再看
		if n.isList && n.needTotal {
			n.total = n.children[n.primaryTableKey].total
		}
	case NodeTypeFunc:
		param := g.Map{}
		n.ret, n.err = functions.Call(n.ctx, n.simpleReqVal, param)
	}

}

func (n *Node) Result() (any, error) {

	if n.err != nil {
		return nil, n.err
	}

	switch n.Type {
	case NodeTypeQuery:
		if n.isList {
			if n.ret == nil || n.ret.([]g.Map) == nil {
				return []g.Map{}, n.err
			}
		} else {
			if n.ret == nil || n.ret.(g.Map) == nil {
				return nil, n.err
			}
		}

	case NodeTypeRef:
		if strings.HasSuffix(n.simpleReqVal, "[]/total") {
			return n.total, nil
		}
	case NodeTypeStruct:
		if n.isList {
			var retList []g.Map

			var primaryList []g.Map

			if n.children[n.primaryTableKey].ret != nil {
				primaryList = n.children[n.primaryTableKey].ret.([]g.Map)

				for i := 0; i < len(primaryList); i++ {

					pItem := primaryList[i]

					item := g.Map{
						n.primaryTableKey: pItem,
					}

					// 遍历组装数据, 后续考虑使用别的方案优化 (暂未简单使用map的id->item ,主要考虑多字段问题)
					for childK, childNode := range n.children {
						if childNode.primaryTableKey == "" {
							if childNode.ret != nil {

								var resultList []g.Map

								for _, depRetItem := range childNode.ret.([]g.Map) {
									match := true
									for refK, refNode := range childNode.refKeyMap {
										if pItem[refNode.column] != depRetItem[refK] {
											match = false
											break
										}
									}
									if match {
										resultList = append(resultList, depRetItem)
									}

								}
								if len(resultList) > 0 {
									if strings.HasSuffix(childK, "[]") {
										item[childK] = resultList
									} else {
										item[childK] = resultList[0]
									}
								}

							}

						}
					}

					retList = append(retList, item)
				}
			}

			n.ret = retList
		} else {

			retMap := g.Map{}
			for k, node := range n.children {
				var err error
				if strings.HasSuffix(k, "@") {
					k = k[0 : len(k)-1]
				}
				if strings.HasSuffix(k, consts.FunctionsKeySuffix) {
					k = k[0 : len(k)-2]
				}

				retMap[k], err = node.Result()
				if node.Type == NodeTypeFunc && retMap[k] == nil {
					delete(retMap, k)
				}

				if err != nil {
					return nil, err
				}
			}
			n.ret = retMap
		}

	}

	return n.ret, n.err

}
