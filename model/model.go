package model

import (
	"github.com/gogf/gf/v2/util/gconv"
)

type Map map[string]any
type MapStrStr map[string]string
type MapStrAny map[string]any
type FuncParam map[string]Var

func (p *FuncParam) Scan(pointer interface{}, mapping ...map[string]string) error {
	return gconv.Scan(p, pointer, mapping...)
}
