package action

type Hook struct {
	// todo 区分 事务内和事务外的hook, 减少事务时长
	Before func(n *Node, method string) error
	After  func(n *Node, method string) error
}

var hooks []Hook

func RegHook(h Hook) {
	hooks = append(hooks, h)
}
