package query

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
)

const (
	totalInList = "[]/total"
)

type RefNode struct {
	node *Node
}

func newRefNode(n *Node) *RefNode {
	return &RefNode{node: n}
}

func (r *RefNode) parse() {
	n := r.node

	refStr := n.simpleReqVal
	if strings.HasPrefix(refStr, "/") { // 这里/开头是相对同级
		refStr = filepath.Dir(n.Path) + refStr
	}
	refPath, refCol := util.ParseRefCol(refStr)
	if refPath == n.Path { // 不能依赖自身
		n.err = consts.NewValidReqErr(fmt.Sprintf("node cannot ref self: (%s)", refPath))
		return
	}

	refNode := n.queryContext.pathNodes[refPath]
	if refNode == nil {
		n.err = consts.NewValidReqErr(fmt.Sprintf(" node %s is nil, but ref by %s", refPath, n.Path))
		return
	}

	n.refKeyMap = make(map[string]NodeRef)

	if strings.HasSuffix(n.simpleReqVal, totalInList) {
		setNeedTotal(refNode)
	}

	n.refKeyMap[n.Key] = NodeRef{
		column: refCol,
		node:   refNode,
	}
}

func (r *RefNode) fetch() {
	n := r.node
	for _, refNode := range n.refKeyMap {
		if strings.HasSuffix(refNode.column, consts.Total) {
			n.total = refNode.node.total
		} else {
			refRet := refNode.node.ret
			switch refRet.(type) {
			case model.Map:
				n.ret = refRet.(model.Map)[refNode.column]
			}
		}
	}
}

func (r *RefNode) result() {
	n := r.node
	if strings.HasSuffix(n.simpleReqVal, totalInList) {
		n.ret = n.total
	}
}

func (r *RefNode) nodeType() int {
	return NodeTypeRef
}
