package apijson

import (
	"context"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/query"
)

type Plugin interface {
	Install(ctx context.Context, a *ApiJson)
}

type ApiJson struct {
	config *config.Config
	Debug  bool // 是否开启debug模式, 显示每步骤
	ctx    context.Context
}

var DefaultApiJson = New()

type App struct {
}

func New() *ApiJson {
	a := &ApiJson{}
	a.config = config.New()
	a.ctx = context.Background()
	return a
}

// Load load for defaultApiJson, 简化使用
func Load(apps ...func(ctx context.Context, a *ApiJson)) *ApiJson {

	for _, app := range apps {
		DefaultApiJson.Use(app)
	}

	DefaultApiJson.Load()
	return DefaultApiJson
}

func (a *ApiJson) Use(p ...func(ctx context.Context, a *ApiJson)) *ApiJson {
	for _, plugin := range p {
		plugin(a.ctx, a)
	}
	return a
}

func (a *ApiJson) Load() {
	a.config.ReLoad()
}

func (a *ApiJson) Config() *config.Config {
	return a.config
}

func (a *ApiJson) NewQuery(ctx context.Context, req model.Map) *query.Query {
	q := query.New(ctx, a.Config().QueryConfig(), req)

	q.DbMeta = a.config.DbMeta
	q.Functions = a.config.Functions
	q.DbFieldStyle = a.config.DbFieldStyle
	q.JsonFieldStyle = a.config.JsonFieldStyle

	q.NoAccessVerify = a.config.Access.NoVerify
	q.AccessCondition = a.config.Access.ConditionFunc

	return q
}

func (a *ApiJson) NewAction(ctx context.Context, method string, req model.Map) *action.Action {
	act := action.New(ctx, a.Config().ActionConfig(), method, req)

	act.NoAccessVerify = a.config.Access.NoVerify
	act.DbFieldStyle = a.config.DbFieldStyle
	act.JsonFieldStyle = a.config.JsonFieldStyle
	act.Functions = a.config.Functions

	return act
}
