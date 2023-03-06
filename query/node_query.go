package query

type queryNode struct {
	node *Node
}

func newQueryNode(n *Node) queryNode {
	return queryNode{node: n}
}

func (q *queryNode) parse() {

}

func (q *queryNode) fetch() {

}

func (q *queryNode) result() {

}

func (q *queryNode) nodeType() int {
	return NodeTypeQuery
}
