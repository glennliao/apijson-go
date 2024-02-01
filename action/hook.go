package action

import (
	"context"
	"net/http"
)

type HookReq struct {
	Node            *Node
	Method          string
	ctx             context.Context
	hooks           []*Hook
	nextId          int
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

		for r.nextId+1 < len(r.hooks) && h == nil {

			r.nextId++

			_h := r.hooks[r.nextId]

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

		if h != nil {
			if r.nextId < len(r.hooks) {
				if r.isInTransaction {
					return h.HandlerInTransaction(r.ctx, r)
				}

				return h.Handler(r.ctx, r)
			}
		}

		return r.handler(r.ctx, r.Node, r.Method)
	}
}

func (r *HookReq) RowKey() any {
	id := r.Node.Data[r.Node.RowKey]
	if r.IsPut() {
		id = r.Node.Where[r.Node.RowKey]
	}
	return id
}

type Hook struct {
	For []string
	// 事务外 ， 可执行参数校验，io等耗时操作
	Handler func(ctx context.Context, req *HookReq) error
	// 事务内，尽量少执行耗时操作 (无论request配置中是否开启事务， 都会先执行handler 然后 在范围内执行HandlerInTransaction)
	HandlerInTransaction func(ctx context.Context, req *HookReq) error
}

type finishHandler func(ctx context.Context, n *Node, method string) error
