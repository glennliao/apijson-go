package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

type RowKeyGenFuncHandler func(ctx context.Context, genParam g.Map, table string, data g.Map) (g.Map, error)

func (c *Config) RowKeyGen(ctx context.Context, genFuncName string, table string, data g.Map) (g.Map, error) {
	if f, exists := c.rowKeyGenFuncMap[genFuncName]; exists {
		return f(ctx, g.Map{}, table, data)
	}

	return nil, nil
}

func (c *Config) RowKeyGenFunc(name string, f RowKeyGenFuncHandler) {
	c.rowKeyGenFuncMap[name] = f
}

func (c *Config) RowKeyGenList() []string {
	return lo.Keys(c.rowKeyGenFuncMap)
}
