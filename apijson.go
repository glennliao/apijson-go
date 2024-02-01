package apijson

import (
	"context"

	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/query"
)

const CtxKey = "apiJsonApp"

type ApiJson struct {
	ctx    context.Context
	config *config.Config
	Debug  bool

	actionHooks   []action.Hook
	actionHookMap map[string][]*action.Hook
}

func New() *ApiJson {
	a := &ApiJson{
		ctx:    context.Background(),
		config: config.New(),
	}
	return a
}

func (a *ApiJson) SetCtx(ctx context.Context) *ApiJson {
	a.ctx = ctx
	return a
}

func (a *ApiJson) Use(p ...func(ctx context.Context, a *ApiJson)) *ApiJson {
	for _, plugin := range p {
		plugin(a.ctx, a)
	}
	return a
}

func (a *ApiJson) Load() error {
	return a.config.ReLoad(a.ctx)
}

func (a *ApiJson) Config() *config.Config {
	return a.config
}

func (a *ApiJson) NewQuery(ctx context.Context, req model.Map) *query.Query {
	ctx = context.WithValue(ctx, CtxKey, a)

	q := query.New(ctx, a.Config().QueryConfig(), req)

	q.DbMeta = a.config.DbMeta
	q.DbFieldStyle = a.config.DbFieldStyle
	q.JsonFieldStyle = a.config.JsonFieldStyle

	q.NoAccessVerify = a.config.Access.NoVerify
	q.AccessCondition = a.config.Access.ConditionFunc

	return q
}

func (a *ApiJson) NewAction(ctx context.Context, method string, req model.Map) *action.Action {
	ctx = context.WithValue(ctx, CtxKey, a)
	act := action.New(ctx, a.Config().ActionConfig(), method, req)

	act.NoAccessVerify = a.config.Access.NoVerify
	act.DbFieldStyle = a.config.DbFieldStyle
	act.JsonFieldStyle = a.config.JsonFieldStyle

	act.HooksMap = a.actionHookMap

	return act
}

func (a *ApiJson) RegActionHook(hook action.Hook) {
	if a.actionHookMap == nil {
		a.actionHookMap = make(map[string][]*action.Hook)
	}
	for _, item := range hook.For {
		a.actionHookMap[item] = append(a.actionHookMap[item], &hook)
	}
}

func GetApiJson(ctx context.Context) *ApiJson {
	return ctx.Value(CtxKey).(*ApiJson)
}
