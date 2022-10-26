package query

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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
	AccessCondition config.AccessCondition
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

// 输出节点信息
func printNode(n *Node, deep int) {

	for i := 0; i < deep; i++ {
		fmt.Print("|")
	}

	desc := gconv.String(n.Type)
	if n.IsList {
		desc += "[]"
	}

	format := fmt.Sprintf("- %%-%ds | %%s\n", 20-deep)

	fmt.Printf(format, n.Key, desc)

	for _, node := range n.children {
		printNode(node, deep+1)
	}

}

func analysisRef(p *Node, fetchNodeQueue *[]*Node, fetchNodeQueueWithRef *[]*Node, fetchNodeQueueRefNode *[]*Node) {

	// 分析依赖关系, 让无依赖的先执行， 然后在执行后续的
	// 需优化调整 更通用的 （目前问题点在于依赖的节点也有依赖时的优先顺序问题）
	for _, node := range p.children {
		if node.Type == NodeTypeQuery && len(node.refKeyMap) == 0 {
			*fetchNodeQueue = append(*fetchNodeQueue, node)
		} else {
			if node.Type == NodeTypeRef {
				*fetchNodeQueueRefNode = append(*fetchNodeQueueRefNode, node)
			} else {
				*fetchNodeQueueWithRef = append(*fetchNodeQueueWithRef, node)
			}
		}
		//if node.primaryTableKey != "" {
		//
		//} else {
		//
		//}
		analysisRef(node, fetchNodeQueue, fetchNodeQueueWithRef, fetchNodeQueueRefNode)
	}

}

func (q *Query) fetch() {
	// 分析依赖关系

	var fetchNodeQueue []*Node
	var fetchNodeQueueWithRef []*Node
	var fetchNodeQueueRefNode []*Node

	analysisRef(q.rootNode, &fetchNodeQueue, &fetchNodeQueueWithRef, &fetchNodeQueueRefNode)

	//fetchNodeQueue = lo.Reverse(fetchNodeQueue)
	//fetchNodeQueueWithRef = lo.Reverse(fetchNodeQueueWithRef)

	//for _, node := range fetchNodeQueueWithRef {
	//	fmt.Printf("%s\n", node.Path)
	//	for k, refPath := range node.refKeyMap {
	//		fmt.Printf("%s -> %s\n", k, refPath)
	//	}
	//	fmt.Println("---------------------")
	//}

	fmt.Println("fetch queue")
	for _, node := range append(fetchNodeQueue) {
		fmt.Printf(" 【%s】 > ", node.Path)
	}
	fmt.Println("")

	for _, node := range append(fetchNodeQueueWithRef) {
		fmt.Printf(" 【%s】 > ", node.Path)
	}

	fmt.Println("")

	for _, node := range append(fetchNodeQueueRefNode) {
		fmt.Printf(" 【%s】 > ", node.Path)
	}

	fmt.Println("")

	for _, node := range fetchNodeQueue {
		node.fetch()
	}

	for _, node := range fetchNodeQueueWithRef {
		node.fetch()
	}

	for _, node := range fetchNodeQueueRefNode {
		node.fetch()
	}

}

func (q *Query) Result() (g.Map, error) {

	g.Log().Debugf(q.ctx, "【query】 ============ [buildNodeTree]")

	// 构建节点树,并校验结构是否符合,  不符合则返回错误, 结束本次查询
	q.rootNode = newNode(q, "", "", q.req)

	err := q.rootNode.buildChild()

	if err != nil {
		return nil, err
	}

	printNode(q.rootNode, 0)

	g.Log().Debugf(q.ctx, "【query】 ============ [parse]")

	setNodeRole(q.rootNode, "", "")

	q.rootNode.parse()

	g.Log().Debugf(q.ctx, "【query】 ============ [fetch]")

	q.fetch()

	resultMap, err := q.rootNode.Result()

	if err != nil {
		if q.rootNode.err != nil {
			return nil, q.rootNode.err
		}

		resultMap := g.Map{}
		for k, node := range q.rootNode.children {
			resultMap[k] = node.err
		}
		return resultMap, err
	}

	g.Log().Debugf(q.ctx, "【query】 ^=======================^")
	return resultMap.(g.Map), err
}
