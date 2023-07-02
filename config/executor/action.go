package executor

import (
	"context"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type ActionExecutor interface {
	Do(ctx context.Context, req ActionExecutorReq) (ret model.Map, err error)
}

type ActionExecutorReq struct {
	Method string
	Table  string
	Data   []model.Map
	Where  []model.Map
	Access *config.AccessConfig
	Config *config.ActionConfig
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
