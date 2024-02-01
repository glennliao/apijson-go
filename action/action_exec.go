package action

import (
	"context"

	"github.com/glennliao/apijson-go/model"

	"github.com/glennliao/apijson-go/consts"
)

func (a *Action) exec() {
	a.ret = model.Map{}
	a.err = a.hookExecute(0, false)
}

// hook like middleware 洋葱模型
func (a *Action) hookExecute(i int, inTransaction bool) error {
	var (
		execNodeKey = a.tagRequestConfig.ExecQueue[i]
		node        = a.children[execNodeKey]
	)

	nodeHookReq := &HookReq{
		ctx:             a.ctx,
		Node:            node,
		Method:          a.Method,
		nextId:          -1,
		isInTransaction: inTransaction,
		hooks:           getHooksByAccessName(a.HooksMap, execNodeKey),
	}

	nodeHookReq.handler = func(ctx context.Context, n *Node, method string) error {
		if i+1 < len(a.tagRequestConfig.ExecQueue) {
			return a.hookExecute(i+1, inTransaction)
		}

		// 执行完了普通hook的before,开始执行事务内
		if !inTransaction {

			transactionHandler := noTransactionHandler

			if a.tagRequestConfig.Transaction != nil && *a.tagRequestConfig.Transaction == true {
				h := GetTransactionHandler(a.ctx, a)
				if h == nil {
					err := consts.NewSysErr("transaction handler is nil")
					return err
				}
				transactionHandler = h
			}

			err := transactionHandler(a.ctx, func(ctx context.Context) error {
				return a.hookExecute(0, !inTransaction)
			})

			return err
		}

		// 执行完 全部前置hook
		var err error

		for _, name := range a.tagRequestConfig.ExecQueue {
			node := a.children[name]
			a.ret[name], err = node.execute(ctx, a.Method)
			if err != nil {
				return err
			}
		}

		return err
	}

	err := nodeHookReq.Next()
	return err
}

func getHooksByAccessName(hooksMap map[string][]*Hook, accessName string) []*Hook {
	hooks := append(hooksMap["*"], hooksMap[accessName]...)
	return hooks
}
