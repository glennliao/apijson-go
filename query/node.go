package query

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
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

type nodeHandler interface {
	parse()
	fetch()
	result()
	nodeType() int
}

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
	page   model.Map // 分页参数

	// 访问当前节点的角色
	role string

	// 节点的请求数据
	req          model.Map
	simpleReqVal string //非对象结构

	// 节点数据执行器
	executor executor.QueryExecutor

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

	nodeHandler nodeHandler

	executorConfig *config.ExecutorConfig
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

	if query.PrintProcessLog {
		g.Log().Debugf(query.ctx, "【node】(%s) <new> ", path)
	}

	node := &Node{
		ctx:          query.ctx,
		queryContext: query,
		Key:          key,
		Path:         path,
		startAt:      time.Now(),
		finish:       false,
	}

	// 节点类型判断
	k, isList := util.ParseNodeKey(key)

	if util.IsFirstUp(k) { // 大写开头, 为查询节点(对应数据库)
		node.Type = NodeTypeQuery
	} else if strings.HasSuffix(k, consts.RefKeySuffix) {
		node.Type = NodeTypeRef
	} else if strings.HasSuffix(k, consts.FunctionsKeySuffix) {
		node.Type = NodeTypeFunc
	} else {
		node.Type = NodeTypeStruct // 结构节点下应该必须存在查询节点

		if query.NoAccessVerify == false {
			if lo.Contains(query.DbMeta.GetTableNameList(), k) {
				node.Type = NodeTypeQuery
			}
		}

	}

	if isList || strings.HasSuffix(filepath.Dir(path), consts.ListKeySuffix) {
		node.isList = true
	}

	switch node.Type {
	case NodeTypeQuery:
		node.nodeHandler = newQueryNode(node)
	case NodeTypeRef:
		node.nodeHandler = newRefNode(node)
	case NodeTypeStruct:
		node.nodeHandler = newStructNode(node)
	case NodeTypeFunc:
		node.nodeHandler = newFuncNode(node)
	}

	switch nodeReq.(type) {
	case map[string]any:
		node.req = nodeReq.(map[string]any)
	case model.Map:
		node.req = nodeReq.(model.Map)
	default:
		node.simpleReqVal = gconv.String(nodeReq)
	}

	return node
}

func (n *Node) buildChild() error {

	if n.Type == NodeTypeQuery && !util.HasFirstUpKey(n.req) { // 查询节点嵌套查询节点, 目前不支持
		return nil
	}

	//最大深度检查
	maxDeep := n.queryContext.queryConfig.MaxTreeDeep()
	if len(strings.Split(n.Path, "/")) > maxDeep {
		return gerror.Newf("deep(%s) > %d", n.Path, maxDeep)
	}

	children := make(map[string]*Node)

	for key, v := range n.req {

		if strings.HasPrefix(key, consts.RefKeySuffix) {
			continue
		}

		if n.Type == NodeTypeQuery && !util.IsFirstUp(key) { // 查询节点嵌套查询节点, 目前不支持
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

		// 最大宽度检查, 为当前节点的子节点数
		maxWidth := n.queryContext.queryConfig.MaxTreeWidth()
		if len(children) > maxWidth {
			path := n.Path
			if path == "" {
				path = "root"
			}
			return gerror.Newf("width(%s) > %d", path, maxWidth)
		}

		n.children = children

		for _, node := range children {
			n.queryContext.pathNodes[node.Path] = node
		}
	}

	return nil
}

func (n *Node) parse() {

	if n.queryContext.PrintProcessLog {
		g.Log().Debugf(n.ctx, "【node】(%s) <parse> ", n.Path)
	}

	n.nodeHandler.parse()

	if n.queryContext.PrintProcessLog {
		g.Log().Debugf(n.ctx, "【node】(%s) <parse-endAt> ", n.Path)
	}

}

func (n *Node) fetch() {

	defer func() {
		n.finish = true
		n.endAt = time.Now()
		if n.queryContext.PrintProcessLog {
			g.Log().Debugf(n.ctx, "【node】(%s) <fetch-endAt> ", n.Path)
		}
	}()

	if n.queryContext.PrintProcessLog {
		g.Log().Debugf(n.ctx, "【node】(%s) <fetch> hasFinish: 【%v】", n.Path, n.finish)
	}

	if n.finish {
		g.Log().Error(n.ctx, "再次执行", n.Path)
		return
	}

	if n.err != nil {
		return
	}

	n.nodeHandler.fetch()

}

func (n *Node) Result() (any, error) {
	if n.err != nil {
		return nil, n.err
	}

	n.nodeHandler.result()

	return n.ret, n.err

}
