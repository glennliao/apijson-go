package executor

import (
	"context"
	"github.com/glennliao/apijson-go/config/db"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

type QueryExecutor interface {
	ParseCondition(conditions g.MapStrAny, accessVerify bool) error
	ParseCtrl(ctrl g.Map) error
	List(page int, count int, needTotal bool) (list []g.Map, total int64, err error)
	One() (g.Map, error)
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
