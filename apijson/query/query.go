package query

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"my-apijson/apijson/util"
	"time"
)

type Query struct {
	ctx context.Context

	req      g.Map // json请求内容
	rootNode *Node // 节点树根节点

	pathNodes map[string]*Node // 用于根据path获取节点

	startAt time.Time
	endAt   time.Time

	err error

	// 是否权限验证
	AccessVerify bool
	// 自定义可访问权限的限定, 例如添加用户id的where条件
	AccessCondition func(ctx context.Context, table string, req g.Map, needRole []string) (g.Map, error)
}

func New(ctx context.Context, req g.Map) *Query {

	g.Log().Debugf(ctx, "【query】 ============ [begin]")

	return &Query{
		ctx:       ctx,
		req:       req,
		startAt:   time.Now(),
		pathNodes: map[string]*Node{},
	}
}

func p(n *Node, deep int) {

	for i := 0; i < deep; i++ {
		fmt.Print("|")
	}
	fmt.Println("-", n.Key)

	for _, node := range n.children {
		p(node, deep+1)
	}

}

func analysisRef(p *Node, fetchNodeQueue *[]*Node, fetchNodeQueueWithRef *[]*Node) {

	for _, node := range p.children {
		if node.isPrimaryTable {
			*fetchNodeQueue = append(*fetchNodeQueue, node)
		} else {
			*fetchNodeQueueWithRef = append(*fetchNodeQueueWithRef, node)
		}
		analysisRef(node, fetchNodeQueue, fetchNodeQueueWithRef)
	}

}

func (q *Query) fetch() {
	// 分析依赖关系

	var fetchNodeQueue []*Node
	var fetchNodeQueueWithRef []*Node

	analysisRef(q.rootNode, &fetchNodeQueue, &fetchNodeQueueWithRef)

	util.Reverse(&fetchNodeQueueWithRef)

	for _, node := range append(fetchNodeQueue) {
		fmt.Printf(" [%s] > ", node.Path)
	}
	fmt.Println("")

	for _, node := range append(fetchNodeQueueWithRef) {
		fmt.Printf(" [%s] > ", node.Path)
	}
	fmt.Println("")

	for _, node := range fetchNodeQueue {
		node.fetch()
	}

	for _, node := range fetchNodeQueueWithRef {
		node.fetch()
	}

}

func (q *Query) Result() (g.Map, error) {
	g.Log().Debugf(q.ctx, "【query】 ============ [buildNodeTree]")
	// buildNodeTree 构建节点树,并校验结构是否符合,  不符合则返回错误, 结束本次查询
	// 最大深度、宽度 为5
	q.rootNode = newNode(q, "@root", "", q.req)
	err := q.rootNode.buildChild()

	if err != nil {
		return nil, err
	}

	g.Log().Debugf(q.ctx, "【query】 ============ [parse]")
	q.rootNode.parse()

	p(q.rootNode, 0)

	g.Log().Debugf(q.ctx, "【query】 ============ [fetch]")
	q.fetch()

	resultMap, err := q.rootNode.Result()

	g.Log().Debugf(q.ctx, "【query】 ^=======================^")
	return resultMap.(g.Map), err
}
