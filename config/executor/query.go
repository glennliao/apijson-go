package executor

import (
	"context"
	"github.com/glennliao/apijson-go/config/db"
	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type QueryExecutor interface {
	ParseCondition(conditions model.MapStrAny, accessVerify bool) error
	ParseCtrl(ctrl model.Map) error
	List(page int, count int, needTotal bool) (list []model.Map, total int64, err error)
	One() (model.Map, error)
	EmptyResult()
}

type queryExecutorBuilder func(ctx context.Context, accessVerify bool, role string, access *db.Access) (QueryExecutor, error)

var queryExecutorBuilderMap = map[string]queryExecutorBuilder{}

func RegQueryExecutor(name string, e queryExecutorBuilder) {
	queryExecutorBuilderMap[name] = e
}

func NewQueryExecutor(name string, ctx context.Context, accessVerify bool, role string, access *db.Access) (QueryExecutor, error) {
	if v, exists := queryExecutorBuilderMap[name]; exists {
		return v(ctx, accessVerify, role, access)
	}
	return queryExecutorBuilderMap["default"](ctx, accessVerify, role, access)
}

func QueryExecutorList() []string {
	return lo.Keys(queryExecutorBuilderMap)
}
