package query

import (
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
)

type funcNode struct {
	node *Node
}

func newFuncNode(n *Node) *funcNode {
	return &funcNode{node: n}
}

func (h *funcNode) parse() {

}

func (h *funcNode) fetch() {
	n := h.node
	queryConfig := n.queryContext.queryConfig

	functionName, paramKeys := util.ParseFunctionsStr(n.simpleReqVal)

	_func := queryConfig.Func(functionName)

	if n.isList && _func.Batch {
		n.later = true
		return
	}

	param := model.Map{}

	for i, item := range _func.ParamList {
		valNode := n.queryContext.pathNodes[paramKeys[i]]
		// if valNode == nil {
		// 	continue
		// }
		if valNode.ret != nil {
			param[item.Name] = valNode.ret
		} else {
			param[item.Name] = valNode.simpleReqVal
		}
	}

	n.ret, n.err = _func.Handler(n.ctx, param)
}

func (h *funcNode) result() {

}

func (h *funcNode) nodeType() int {
	return NodeTypeStruct
}
