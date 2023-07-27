package query

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/glennliao/apijson-go/consts"
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

}

func (h *funcNode) result() {
	n := h.node
	queryConfig := n.queryContext.queryConfig

	functionName, paramKeys := util.ParseFunctionsStr(n.simpleReqVal)

	_func := queryConfig.Func(functionName)

	if _func == nil {
		n.err = consts.NewValidReqErr("functions not exists: " + functionName)
		return
	}

	// todo batch support

	param := model.Map{}

	for i, item := range _func.ParamList {
		paramName := paramKeys[i]
		if strings.HasPrefix(paramName, "/") { // 这里/开头是相对同级
			dir := filepath.Dir(n.Path)
			if dir == "." {
				dir = ""
				paramName = paramName[1:]
			}
			paramName = dir + paramName
		}
		refPath, paramName := util.ParseRefCol(paramName)
		if refPath == n.Path { // 不能依赖自身
			n.err = consts.NewValidReqErr(fmt.Sprintf("node cannot ref self: (%s)", refPath))
			return
		}

		valNode := n.queryContext.pathNodes[refPath]
		if valNode == nil {
			h.node.err = consts.NewValidReqErr(fmt.Sprintf("param %s no found on %s", paramKeys[i], functionName))
			return
		}
		if valNode.ret != nil {

			switch valNode.ret.(type) {
			case model.Map:
				param[item.Name] = util.String(valNode.ret.(model.Map)[paramName])
			case string:
				param[item.Name] = valNode.ret.(string)
			}

		} else {
			param[item.Name] = util.String(valNode.simpleReqVal)
		}
	}

	n.ret, n.err = queryConfig.CallFunc(n.ctx, functionName, param)
}

func (h *funcNode) nodeType() int {
	return NodeTypeStruct
}
