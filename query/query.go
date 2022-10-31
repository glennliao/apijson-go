package query

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"strings"
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
	if n.isList {
		desc += "[]"
	}

	format := fmt.Sprintf("- %%-%ds | %%s\n", 20-deep)

	fmt.Printf(format, n.Key, desc)

	for _, node := range n.children {
		printNode(node, deep+1)
	}

}

func (q *Query) fetch() {
	// 分析依赖关系

	var prerequisites [][]string
	analysisRef(q.rootNode, &prerequisites)
	fetchQueue, err := analysisOrder(prerequisites)

	if err != nil {
		q.err = err
		return
	}

	for k, _ := range q.pathNodes {
		if !lo.Contains(fetchQueue, k) {
			fetchQueue = append(fetchQueue, k)
		}
	}

	g.Log().Debugf(q.ctx, "fetch queue： %s", strings.Join(fetchQueue, " > "))

	for _, path := range fetchQueue {
		q.pathNodes[path].fetch()
	}

	q.rootNode.fetch()
}

func (q *Query) Result() (g.Map, error) {

	g.Log().Debugf(q.ctx, "【query】 ============ [buildNodeTree]")

	// 构建节点树,并校验结构是否符合,  不符合则返回错误, 结束本次查询
	q.rootNode = newNode(q, "", "", q.req)

	err := q.rootNode.buildChild()

	if err != nil {
		return nil, err
	}

	if config.Debug {
		printNode(q.rootNode, 0)
	}

	g.Log().Debugf(q.ctx, "【query】 ============ [parse]")

	setNodeRole(q.rootNode, "", "")

	q.rootNode.parse()

	g.Log().Debugf(q.ctx, "【query】 ============ [fetch]")

	q.fetch()

	if q.err != nil {
		return nil, q.err
	}

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
