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
	//n := h.node
	//functionName, _ := util.ParseFunctionsStr(n.simpleReqVal)
	//n.simpleReqVal = functionName
}

func (h *funcNode) fetch() {
	n := h.node
	functionName, paramKeys := util.ParseFunctionsStr(n.simpleReqVal)
	//n.simpleReqVal = functionName
	// todo 如何传递参数

	param := model.Map{}
	for _, key := range paramKeys {
		param[key] = n.queryContext.pathNodes[key].simpleReqVal
	}
	n.ret, n.err = n.queryContext.queryConfig.CallFunc(n.ctx, functionName, param)
}

func (h *funcNode) result() {
	//n := h.node

}

func (h *funcNode) nodeType() int {
	return NodeTypeStruct
}
