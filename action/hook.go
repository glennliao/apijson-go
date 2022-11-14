package action

type Hook struct {
	Before func(n *Node, method string) error
	After  func(n *Node, method string) error
}

var hooks []Hook

func RegHook(h Hook) {
	hooks = append(hooks, h)
}
