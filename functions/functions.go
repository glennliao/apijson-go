package functions

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type Func struct {
	Name    string
	Handler func(ctx context.Context, param g.Var) (res g.Var, err error)
}

func Reg(name string) {
}
