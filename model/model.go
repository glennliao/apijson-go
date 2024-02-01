package model

import (
	"github.com/gogf/gf/v2/util/gconv"
)

type (
	Map       map[string]interface{}
	MapStrStr map[string]string
	MapStrAny map[string]any
	FuncParam map[string]Var
)

func (p *FuncParam) Scan(pointer interface{}, mapping ...map[string]string) error {
	return gconv.Scan(p, pointer, mapping...)
}
