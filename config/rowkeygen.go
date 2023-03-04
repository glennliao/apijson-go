package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

type RowKeyGenFuncHandler func(ctx context.Context, genParam g.Map, table string, data g.Map) (g.Map, error)

func (a *Config) RowKeyGen(ctx context.Context, genFuncName string, table string, data g.Map) (g.Map, error) {
	if f, exists := a.rowKeyGenFuncMap[genFuncName]; exists {
		return f(ctx, g.Map{}, table, data)
	}

	return nil, nil
}

func (a *Config) RowKeyGenFunc(name string, f RowKeyGenFuncHandler) {
	a.rowKeyGenFuncMap[name] = f
}

func (a *Config) RowKeyGenList() []string {
	return lo.Keys(a.rowKeyGenFuncMap)
}
