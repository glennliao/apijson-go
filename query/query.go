package query

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
)

type Query struct {
	ctx context.Context

	// json请求内容
	req model.Map
	// 节点树根节点
	rootNode *Node

	// 用于根据path获取节点
	pathNodes map[string]*Node

	startAt time.Time
	endAt   time.Time

	err error

	// 输出过程
	PrintProcessLog bool

	// 关闭权限验证 , 默认否
	NoAccessVerify bool

	queryConfig *config.QueryConfig

	// 自定义可访问权限的限定, 例如添加用户id的where条件
	AccessCondition config.AccessCondition

	DbMeta *config.DBMeta

	// dbFieldStyle 数据库字段命名风格 请求传递到数据库中
	DbFieldStyle config.FieldStyle

	// jsonFieldStyle 数据库返回的字段
	JsonFieldStyle config.FieldStyle
}

func New(ctx context.Context, qc *config.QueryConfig, req model.Map) *Query {
	q := &Query{
		queryConfig: qc,
	}

	q.init(ctx, req)
	q.NoAccessVerify = qc.NoVerify()

	return q
}

func (q *Query) init(ctx context.Context, req model.Map) {
	q.ctx = ctx
	q.req = req

	q.startAt = time.Now()
	q.pathNodes = make(map[string]*Node)
}

func (q *Query) Result() (model.Map, error) {
	if q.PrintProcessLog {
		g.Log().Debugf(q.ctx, "【query】 ============ [begin]")
		g.Log().Debugf(q.ctx, "【query】 ============ [buildNodeTree]")
	}

	// 构建节点树,并校验结构是否符合,  不符合则返回错误, 结束本次查询
	q.rootNode = newNode(q, "", "", q.req)

	err := q.rootNode.buildChild()
	if err != nil {
		return nil, err
	}

	if q.PrintProcessLog {
		q.printNode(q.rootNode, 0)
	}

	if q.PrintProcessLog {
		g.Log().Debugf(q.ctx, "【query】 ============ [parse]")
	}

	setNodeRole(q.rootNode, "", "")

	q.rootNode.parse()

	if q.PrintProcessLog {
		g.Log().Debugf(q.ctx, "【query】 ============ [fetch]")
	}

	q.fetch()

	if q.err != nil {
		return nil, q.err
	}

	resultMap, err := q.rootNode.Result()
	if err != nil {
		if q.rootNode.err != nil {
			return nil, q.rootNode.err
		}

		ret := model.Map{}
		for k, node := range q.rootNode.children {
			ret[k] = node.err
		}
		return ret, err
	}

	if q.PrintProcessLog {
		g.Log().Debugf(q.ctx, "【query】 ^=======================^")
	}

	return resultMap.(model.Map), err
}

func (q *Query) fetch() {
	// 分析依赖关系

	var prerequisites [][]string
	analysisRef(q.rootNode, &prerequisites)
	fetchQueue, err := util.AnalysisOrder(prerequisites)
	if err != nil {
		q.err = err
		return
	}

	for k := range q.pathNodes {
		if !lo.Contains(fetchQueue, k) {
			fetchQueue = append(fetchQueue, k)
		}
	}

	if q.PrintProcessLog {
		g.Log().Debugf(q.ctx, "fetch queue： %s", strings.Join(fetchQueue, " > "))
	}

	for _, path := range fetchQueue {
		q.pathNodes[path].fetch()
	}

	q.rootNode.fetch()
}

// 输出节点信息
func (q *Query) printNode(n *Node, deep int) {
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
		q.printNode(node, deep+1)
	}
}
