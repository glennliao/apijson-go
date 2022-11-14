package action

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type RowKeyGenFuncHandler func(ctx context.Context, table string, data g.Map) (g.Map, error)

var rowKeyGen RowKeyGenFuncHandler = func(ctx context.Context, table string, data g.Map) (g.Map, error) {
	return nil, nil
}

func RowKeyGenFunc(f RowKeyGenFuncHandler) {
	rowKeyGen = f
}
