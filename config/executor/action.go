package executor

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

type ActionExecutor interface {
	Insert(ctx context.Context, table string, data any) (id int64, count int64, err error)
	Update(ctx context.Context, table string, data g.Map, where g.Map) (count int64, err error)
	Delete(ctx context.Context, table string, where g.Map) (count int64, err error)
}

var actionExecutorMap = map[string]ActionExecutor{}

func RegActionExecutor(name string, e ActionExecutor) {
	actionExecutorMap[name] = e
}

func GetActionExecutor(name string) ActionExecutor {
	if v, exists := actionExecutorMap[name]; exists {
		return v
	}
	return actionExecutorMap["default"]
}

func ActionExecutorList() []string {
	return lo.Keys(actionExecutorMap)
}
