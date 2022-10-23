package query

import (
	"context"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"my-apijson/apijson/db"
	"my-apijson/apijson/util"
	"path/filepath"
	"strings"
	"time"
)

const (
	UNKNOWN = "UNKNOWN" // 未登录用户
	LOGIN   = "LOGIN"   // 登录用户 (用于需要登录才能查看的公开资源)
	OWNER   = "OWNER"   // 用户 自己创建的数据
	ADMIN   = "ADMIN"   // 管理员

	// 以下非apijson官方的

	DEPT = "DEPT" // 机构管理员
)

type RefNode struct {
	column string
	node   *Node
}

// 设置嵌套的最大深度 https://github.com/Tencent/APIJSON/issues/147

const MaxTreeWidth = 5
const MaxTreeDeep = 5

type Node struct {
	ctx          context.Context
	queryContext *Query

	Key  string
	Path string

	req         g.MapStrAny
	sqlExecutor *db.SqlExecutor

	IsList  bool
	startAt time.Time
	endAt   time.Time

	ret any
	err error

	depRetList []g.Map // 返回给主表组装

	children map[string]*Node

	refKeyMap  g.Map // 关联字段
	refNodeMap map[string]RefNode

	isPrimaryTable bool // 是否主查询表

	Finish bool // 执行完毕

	isAccess bool // 是否可访问
}

func newNode(ctx *Query, key string, path string, nodeReq g.Map) *Node {

	if len(strings.Split(path, "/")) > MaxTreeDeep {
		panic(gerror.Newf("width(%s) > 5", path))
	}

	g.Log().Debugf(ctx.ctx, "【node】(%s) <new> ", path)

	return &Node{
		ctx:          ctx.ctx,
		queryContext: ctx,
		Key:          key,
		Path:         path,
		req:          nodeReq,
		startAt:      time.Now(),
		Finish:       false,
		isAccess:     true,
	}
}

func (n *Node) buildChild() error {

	if len(n.req) > MaxTreeWidth {
		path := n.Path
		if path == "" {
			path = "root"
		}
		return gerror.Newf("width(%s) > 5", path)
	}

	n.children = make(map[string]*Node)

	for k, v := range n.req {
		nodeReq, ok := v.(g.Map)
		if ok {
			path := n.Path
			if path != "" { // 根节点时不带/
				path += "/"
			}
			path += k
			node := newNode(n.queryContext, k, path, nodeReq)

			err := node.buildChild()
			if err != nil {
				return err
			}

			n.children[k] = node
			n.queryContext.pathNodes[path] = node
		}
	}

	return nil
}

func (n *Node) parse() {

	g.Log().Debugf(n.ctx, "【node】(%s) <parse> ", n.Path)

	table, isList := parseTableKey(n.Key, n.Path)

	n.IsList = isList

	if table != "" {

		role := n.ctx.Value("ajg.role").([]string)

		access := db.AccessMap[table]

		canAccess := false
		for _, r := range role {
			if util.Contains(access.Get, r) {
				canAccess = true
				break
			}
		}
		g.Log().Debug(n.ctx, "userRole:", role, "accessRole", access.Get, "can?:", canAccess)
		if !canAccess {

			g.Log().Debug(n.ctx, n.Path, "userRole:", role, "accessRole", access.Get, "can?:", canAccess)
			return
		}

		refKeyMap, conditionMap := parseRefKey(n.req)

		n.sqlExecutor = db.NewSqlExecutor(n.ctx, table, isList)
		// 查询条件
		err := n.sqlExecutor.ParseCondition(conditionMap)

		if err != nil {
			n.err = err
			return
		}

		// access 限定条件
		if n.queryContext.AccessCondition != nil {
			where, err := n.queryContext.AccessCondition(n.ctx, table, n.req, access.Get)
			if err != nil {
				n.err = err
				return
			}
			if where != nil {
				err = n.sqlExecutor.ParseCondition(where)

				if err != nil {
					n.err = err
					return
				}
			}

		}

		if len(refKeyMap) > 0 {
			n.refKeyMap = refKeyMap
			n.refNodeMap = make(map[string]RefNode)

			hasRefBrother := false // 是否依赖同级节点

			for _refK, _refPath := range refKeyMap {

				refPath := gconv.String(_refPath)
				refPathCol := filepath.Base(refPath)                  // "id@":"[]/T-odo/userId"  ->  userId
				refPath = refPath[0 : len(refPath)-len(refPathCol)-1] // "id@":"[]/T-odo/userId"  ->  []/T-odo   不加横杠会自动变成goland的to_do 项

				if strings.HasPrefix(refPath, "/") { // 有点非正常绝对路径的写法, 这里/开头是相对同级
					refPath = filepath.Dir(n.Path) + refPath
				}

				if !hasRefBrother {
					if filepath.Dir(n.Path) == filepath.Dir(refPath) {
						hasRefBrother = true
					}
				}

				if refPath == n.Path { // 不能依赖自身
					panic(gerror.Newf("node cannot ref self: (%s) {%s:%s}", refPath, _refK, _refPath))
				}

				refNode := n.queryContext.pathNodes[refPath]

				if refNode == nil {
					panic(gerror.Newf("node %s is nil", refPath))
				}

				n.refNodeMap[_refK] = RefNode{
					column: refPathCol,
					node:   refNode,
				}
			}

			n.isPrimaryTable = !hasRefBrother

		} else {
			n.isPrimaryTable = true
		}

	} else { // key 为 []

		page := 1
		count := 10

		for k, v := range n.req {

			if _, ok := v.(g.Map); ok {
				continue
			}

			switch k {
			case "page":
				page = gconv.Int(v)

			case "count":
				count = gconv.Int(v)

			}
		}

		for _, childNode := range n.children {
			childNode.parse()
		}

		hasPrimary := false
		for _, n := range n.children {
			if n.isPrimaryTable { // 主查询表 才分页
				err := n.sqlExecutor.ParseCondition(g.Map{
					"page":  page,
					"count": count,
				})
				if err != nil {
					n.err = err
					return
				}
				hasPrimary = true
			} else {
				if n.sqlExecutor != nil {
					n.sqlExecutor.ParseCondition(g.Map{
						"page":  0,
						"count": 0,
					})
				}

			}
		}
		if n.Key == "[]" && !hasPrimary {
			panic(gerror.Newf("node must have  primary table: (%s)", n.Path))
		}
	}

	g.Log().Debugf(n.ctx, "【node】(%s) <parse-endAt> ", n.Path)

}

func (n *Node) fetch() {
	g.Log().Debugf(n.ctx, "【node】(%s) <fetch> hasFinish: 【%v】", n.Path, n.Finish)

	if n.Finish {
		return
	}

	if n.isPrimaryTable {
		if n.sqlExecutor != nil {
			n.ret, n.err = n.sqlExecutor.Fetch()
		}

		for _, node := range n.children {
			if len(node.refKeyMap) == 0 {
				node.fetch()
			}
		}
	} else {
		g.Log().Debug(n.ctx, "[dep]", n.Path, " -> ", n.refKeyMap)
		for k, refNode := range n.refNodeMap {

			ret, err := refNode.node.Result()
			if err != nil {
				g.Log().Error(n.ctx, "", err)
				n.err = err
			}

			if refNode.node.IsList {
				list := ret.(gdb.Result)
				valList := list.Array(refNode.column)

				set := gset.New()
				for _, value := range valList {
					set.Add(gconv.String(value))
				}

				err := n.sqlExecutor.ParseCondition(g.Map{
					k + "{}": set.Slice(), // todo @ 与 {}&等的结合 id{}@的处理
				})

				if err != nil {
					n.err = err
					return
				}

			} else {
				refConditionMap := g.Map{}
				item := ret.(gdb.Record)
				refConditionMap[k] = item.Map()[refNode.column]
				err := n.sqlExecutor.ParseCondition(refConditionMap)
				if err != nil {
					n.err = err
					return
				}
			}
		}

		if n.sqlExecutor != nil {
			ret, err := n.sqlExecutor.Fetch()
			if err != nil {
				n.err = err
				return
			}

			if n.IsList {
				retList := ret.(gdb.Result)
				var depRetList []g.Map
				for _, record := range retList {
					depRetList = append(depRetList, record.Map())
				}
				n.depRetList = depRetList
			} else {
				record := ret.(gdb.Record)
				n.ret = record.Map()
			}

			if n.err != nil {
				return
			}
		}

	}

	n.Finish = true
	n.endAt = time.Now()

	g.Log().Debugf(n.ctx, "【node】(%s) <fetch-endAt> ", n.Path)

}

func (n *Node) Result() (any, error) {

	if n.sqlExecutor != nil { // children == 0
		return n.ret, n.err
	}

	if n.IsList { // []组装数据

		retList := g.List{}

		var primaryList gdb.Result
		var primaryKey string

		for _, node := range n.children {
			ret, err := node.Result()
			if err != nil {
				panic(err)
			}
			if node.isPrimaryTable {
				primaryList = ret.(gdb.Result)
				primaryKey = node.Key
			}
		}

		for i := 0; i < len(primaryList); i++ {

			pItem := primaryList[i].Map()

			item := g.Map{
				primaryKey: pItem,
			}

			// 遍历组装数据, 后续考虑使用别的方案优化 (暂未简单使用map的id->item ,主要考虑多字段问题)
			for childK, childNode := range n.children {
				if !childNode.isPrimaryTable {
					for _, depRetItem := range childNode.depRetList {
						for refK, refNode := range childNode.refNodeMap {
							if pItem[refNode.column] == depRetItem[refK] {
								item[childK] = depRetItem
							}
						}
					}
				}
			}

			retList = append(retList, item)
		}

		n.ret = retList
	} else {
		retMap := g.Map{}
		for k, node := range n.children {
			var err error
			retMap[k], err = node.Result()
			if err != nil {
				panic(err)
			}
		}
		n.ret = retMap
	}

	return n.ret, n.err
}

func parseTableKey(k string, p string) (tableName string, isList bool) {

	if k == "@root" {
		return "", false
	}

	if strings.HasSuffix(k, "[]") {
		return k[0 : len(k)-2], true
	}

	if strings.Contains(p, "[]") {
		return k, true
	}

	return k, false
}
func parseRefKey(reqMap g.Map) (g.Map, g.Map) {
	depMap := g.Map{}
	otherKeyMap := g.Map{}
	for k, v := range reqMap {
		if strings.HasSuffix(k, "@") {
			depMap[k[0:len(k)-1]] = gconv.String(v)
		} else {
			otherKeyMap[k] = v
		}
	}

	return depMap, otherKeyMap
}
