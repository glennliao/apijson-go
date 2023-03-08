package query

import (
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"path/filepath"
	"strings"
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
		panic(gerror.Newf("node cannot ref self: (%s)", refPath))
	}

	refNode := n.queryContext.pathNodes[refPath]
	if refNode == nil {
		panic(gerror.Newf(" node %s is nil, but ref by %s", refPath, n.Path))
	}

	n.refKeyMap = make(map[string]NodeRef)

	if strings.HasSuffix(n.simpleReqVal, "[]/total") {
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
		if strings.HasSuffix(refNode.column, "total") && strings.HasSuffix(refNode.node.Path, consts.ListKeySuffix) {
			n.total = refNode.node.total
		} else {
			n.ret = refNode.node.ret.(model.Map)[refNode.column] //todo fei model.Map 时候
		}
	}
}

func (r *RefNode) result() {
	n := r.node
	if strings.HasSuffix(n.simpleReqVal, "[]/total") {
		n.ret = n.total
	}
}

func (r *RefNode) nodeType() int {
	return NodeTypeRef
}
