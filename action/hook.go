package action

import "context"

const (
	BeforeExec = iota
	AfterExec
	BeforeDo
	AfterDo
)

type Hook struct {
	For string //
	// Exec 事务外
	BeforeExec func(ctx context.Context, n *Node, method string) error
	AfterExec  func(ctx context.Context, n *Node, method string) error

	// Do 事务内
	BeforeDo func(ctx context.Context, n *Node, method string) error
	AfterDo  func(ctx context.Context, n *Node, method string) error
}

var hooksMap = map[string][]Hook{}

func RegHook(h Hook) {
	hooksMap[h.For] = append(hooksMap[h.For], h)
}

func EmitHook(ctx context.Context, hookAt int, node *Node, method string) error {

	hooks := append(hooksMap["*"], hooksMap[node.TableName]...)
	for _, hook := range hooks {

		var handler func(ctx context.Context, n *Node, method string) error
		switch hookAt {
		case BeforeExec:
			handler = hook.BeforeExec
		case AfterExec:
			handler = hook.AfterExec
		case BeforeDo:
			handler = hook.BeforeDo
		case AfterDo:
			handler = hook.AfterDo
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
