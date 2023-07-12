package action

import (
	"context"
	"net/http"
)

// const (
// 	BeforeNodeExec = iota
// 	AfterNodeExec
// 	BeforeExecutorDo
// 	AfterExecutorDo
// )

type HookReq struct {
	Node            *Node
	Method          string
	ctx             context.Context
	hooks           []*Hook
	nextIdx         int
	isInTransaction bool
	handler         finishHandler
}

func (r *HookReq) IsPost() bool {
	return r.Method == http.MethodPost
}

func (r *HookReq) IsPut() bool {
	return r.Method == http.MethodPut
}

func (r *HookReq) IsDelete() bool {
	return r.Method == http.MethodDelete
}

func (r *HookReq) Next() error {

	for {

		var h *Hook

		for r.nextIdx < len(r.hooks) && h == nil {

			if r.nextIdx+1 >= len(r.hooks) {
				if r.isInTransaction {
					// finish all
					return r.handler(r.ctx, r.Node, r.Method)
				} else {
					r.nextIdx = -1
					r.isInTransaction = true
				}
			}

			r.nextIdx++

			_h := r.hooks[r.nextIdx]

			if r.isInTransaction {
				if _h.HandlerInTransaction == nil {
					continue
				}
				h = _h
			} else {
				if _h.Handler == nil {
					continue
				}
				h = _h
			}

		}

		if r.nextIdx < len(r.hooks) {
			if r.isInTransaction {
				return h.HandlerInTransaction(r.ctx, r)
			}

			return h.Handler(r.ctx, r)
		}
	}

}

type Hook struct {
	For []string
	// 事务外 ， 可执行参数校验，io等耗时操作
	Handler func(ctx context.Context, req *HookReq) error
	// 事务内，尽量少执行耗时操作 (无论request配置中是否开启事务， 都会先执行handler 然后 在范围内执行HandlerInTransaction)
	HandlerInTransaction func(ctx context.Context, req *HookReq) error
}

type finishHandler func(ctx context.Context, n *Node, method string) error

func getHooksByAccessName(hooksMap map[string][]*Hook, accessName string) []*Hook {
	hooks := append(hooksMap["*"], hooksMap[accessName]...)
	return hooks
}

//
// type Hook2 struct {
// 	For []string //
// 	// Exec 事务外
// 	BeforeNodeExec func(ctx context.Context, n *Node, method string) error
// 	AfterNodeExec  func(ctx context.Context, n *Node, method string) error
//
// 	// Do 事务内
// 	BeforeExecutorDo func(ctx context.Context, n *Node, method string) error
// 	AfterExecutorDo  func(ctx context.Context, n *Node, method string) error
// }
//
// func emitHook(ctx context.Context, hooksMap map[string][]Hook, hookAt int, node *Node, method string) error {
//
// 	hooks := append(hooksMap["*"], hooksMap[node.Key]...)
// 	for _, hook := range hooks {
//
// 		var handler func(ctx context.Context, n *Node, method string) error
// 		switch hookAt {
// 		case BeforeNodeExec:
// 			handler = hook.BeforeNodeExec
// 		case AfterNodeExec:
// 			handler = hook.AfterNodeExec
// 		case BeforeExecutorDo:
// 			handler = hook.BeforeExecutorDo
// 		case AfterExecutorDo:
// 			handler = hook.AfterExecutorDo
// 		}
//
// 		if handler != nil {
// 			err := handler(ctx, node, method)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
