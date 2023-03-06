package apijson

import (
	"context"
	"github.com/glennliao/apijson-go/config"
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

}

func (a *ApiJson) Config() *config.Config {
	return a.config
}
