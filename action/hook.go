package action

type Hook struct {
	// Exec 事务外
	BeforeExec func(n *Node, method string) error
	AfterExec  func(n *Node, method string) error

	// Do 事务内
	BeforeDo func(n *Node, method string) error
	AfterDo  func(n *Node, method string) error
}

var hooks []Hook

func RegHook(h Hook) {
	hooks = append(hooks, h)
}
