package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

type RowKeyGenFuncHandler func(ctx context.Context, genParam g.Map, table string, data g.Map) (g.Map, error)

var rowKeyGenFuncMap = map[string]RowKeyGenFuncHandler{}

func RowKeyGen(ctx context.Context, genFuncName string, table string, data g.Map) (g.Map, error) {
	if f, exists := rowKeyGenFuncMap[genFuncName]; exists {
		return f(ctx, g.Map{}, table, data)
	}

	return nil, nil
}

func RowKeyGenFunc(name string, f RowKeyGenFuncHandler) {
	rowKeyGenFuncMap[name] = f
}

func RowKeyGenList() []string {
	return lo.Keys(rowKeyGenFuncMap)
}
