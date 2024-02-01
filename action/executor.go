package action

import (
	"context"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type ExecutorReq struct {
	Method          string
	Table           string
	Data            model.Map
	Where           model.Map
	AccessCondition *config.ConditionRet
	Access          *config.AccessConfig
	Config          *config.ActionConfig
}

var actionExecutorMap = map[string]Executor{}

func RegExecutor(name string, e Executor) {
	actionExecutorMap[name] = e
}

func GetActionExecutor(name string) (Executor, error) {
	if name == "" {
		name = "default"
	}
	if v, exists := actionExecutorMap[name]; exists {
		return v, nil
	}
	return nil, consts.NewSysErr("action executor not found: " + name)
}

func ExecutorList() []string {
	return lo.Keys(actionExecutorMap)
}

type Executor interface {
	Do(ctx context.Context, req ExecutorReq) (ret model.Map, err error)
}

// TransactionHandler 事务处理函数

type TransactionHandler func(ctx context.Context, action func(ctx context.Context) error) error

type TransactionResolver func(ctx context.Context, req *Action) TransactionHandler

var noTransactionHandler = func(ctx context.Context, action func(ctx context.Context) error) error {
	return action(ctx)
}

var transactionResolver TransactionResolver

func RegTransactionResolver(r TransactionResolver) {
	transactionResolver = r
}

func GetTransactionHandler(ctx context.Context, req *Action) TransactionHandler {
	return transactionResolver(ctx, req)
}
