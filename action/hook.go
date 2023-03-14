package action

import "context"

const (
	BeforeNodeExec = iota
	AfterNodeExec
	BeforeExecutorDo
	AfterExecutorDo
)

type Hook struct {
	For string //
	// Exec 事务外
	BeforeNodeExec func(ctx context.Context, n *Node, method string) error
	AfterNodeExec  func(ctx context.Context, n *Node, method string) error

	// Do 事务内
	BeforeExecutorDo func(ctx context.Context, n *Node, method string) error
	AfterExecutorDo  func(ctx context.Context, n *Node, method string) error
}

var hooksMap = map[string][]Hook{}

func RegHook(h Hook) {
	hooksMap[h.For] = append(hooksMap[h.For], h)
}

func EmitHook(ctx context.Context, hookAt int, node *Node, method string) error {

	hooks := append(hooksMap["*"], hooksMap[node.Key]...)
	for _, hook := range hooks {

		var handler func(ctx context.Context, n *Node, method string) error
		switch hookAt {
		case BeforeNodeExec:
			handler = hook.BeforeNodeExec
		case AfterNodeExec:
			handler = hook.AfterNodeExec
		case BeforeExecutorDo:
			handler = hook.BeforeExecutorDo
		case AfterExecutorDo:
			handler = hook.AfterExecutorDo
		}

		if handler != nil {
			err := handler(ctx, node, method)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
