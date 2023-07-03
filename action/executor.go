package action

import (
	"context"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/query"
	"github.com/samber/lo"
)

type ActionExecutorReq struct {
	Method   string
	Table    string
	Data     []model.Map
	Where    []model.Map
	Access   *config.AccessConfig
	Config   *config.ActionConfig
	NewQuery func(ctx context.Context, req model.Map) *query.Query
}

var actionExecutorMap = map[string]ActionExecutor{}

func RegExecutor(name string, e ActionExecutor) {
	actionExecutorMap[name] = e
}

func GetActionExecutor(name string) (ActionExecutor, error) {
	if name == "" {
		name = "default"
	}
	if v, exists := actionExecutorMap[name]; exists {
		return v, nil
	}
	return nil, consts.NewSysErr("action executor not found: " + name)
}

func ActionExecutorList() []string {
	return lo.Keys(actionExecutorMap)
}

type ActionExecutor interface {
	Do(ctx context.Context, req ActionExecutorReq) (ret model.Map, err error)
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
