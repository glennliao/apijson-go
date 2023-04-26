package executor

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type QueryExecutor interface {
	ParseCondition(conditions model.MapStrAny, accessVerify bool) error
	ParseCtrl(ctrl model.Map) error
	List(page int, count int) (list []model.Map, err error)
	Count() (total int64, err error)
	One() (model.Map, error)
	SetEmptyResult()
}

type queryExecutorBuilder func(ctx context.Context, config *config.ExecutorConfig) (QueryExecutor, error)

var queryExecutorBuilderMap = map[string]queryExecutorBuilder{}

func RegQueryExecutor(name string, e queryExecutorBuilder) {
	queryExecutorBuilderMap[name] = e
}

func NewQueryExecutor(name string, ctx context.Context, config *config.ExecutorConfig) (QueryExecutor, error) {
	if v, exists := queryExecutorBuilderMap[name]; exists {
		return v(ctx, config)
	}
	return queryExecutorBuilderMap["default"](ctx, config)
}

func QueryExecutorList() []string {
	return lo.Keys(queryExecutorBuilderMap)
}
