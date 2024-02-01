package query

const (
	NodeTypeStruct = iota // 结构节点
	NodeTypeQuery         // 查询节点
	NodeTypeRef           // 引用节点
	NodeTypeFunc          // functions 节点
)

type nodeHandler interface {
	parse()
	fetch()
	result()
	nodeType() int
}
