package query

import (
	"context"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type Executor interface {
	ParseCondition(conditions *config.ConditionRet, accessVerify bool) error
	ParseCtrl(ctrl model.Map) error
	List(page int, count int) (list []model.Map, err error)
	Count() (total int64, err error)
	One() (model.Map, error)
	SetEmptyResult()
}

type queryExecutorBuilder func(ctx context.Context, config *config.ExecutorConfig) (Executor, error)

var queryExecutorBuilderMap = map[string]queryExecutorBuilder{}

func RegExecutor(name string, e queryExecutorBuilder) {
	queryExecutorBuilderMap[name] = e
}

func NewExecutor(name string, ctx context.Context, config *config.ExecutorConfig) (Executor, error) {
	if name == "" {
		name = "default"
	}

	if v, exists := queryExecutorBuilderMap[name]; exists {
		return v(ctx, config)
	}

	return nil, consts.NewSysErr("query executor not found: " + name)
}

func ExecutorList() []string {
	return lo.Keys(queryExecutorBuilderMap)
}
