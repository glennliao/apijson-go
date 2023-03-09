package config

import (
	"context"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type RowKeyGenReq struct {
	AccessName string
	Data       model.Map
}

type RowKeyGenRet struct {
	data model.Map
}

func NewRowKeyGenRet() *RowKeyGenRet {
	return &RowKeyGenRet{data: map[string]any{}}
}

func (r *RowKeyGenRet) RowKey(id any) {
	r.data[consts.RowKey] = id
}
func (r *RowKeyGenRet) RowKeys(d model.Map) {
	for k, v := range d {
		r.data[k] = v
	}
}

type RowKeyGenFuncHandler func(ctx context.Context, req *RowKeyGenReq, ret *RowKeyGenRet) error

func (c *Config) RowKeyGen(ctx context.Context, genFuncName string, accessName string, data model.Map) (model.Map, error) {
	if f, exists := c.rowKeyGenFuncMap[genFuncName]; exists {
		req := &RowKeyGenReq{
			AccessName: accessName,
			Data:       data,
		}
		ret := NewRowKeyGenRet()
		err := f(ctx, req, ret)
		return ret.data, err
	}

	return nil, nil
}

func (c *Config) RowKeyGenFunc(name string, f RowKeyGenFuncHandler) {
	c.rowKeyGenFuncMap[name] = f
}

func (c *Config) RowKeyGenList() []string {
	return lo.Keys(c.rowKeyGenFuncMap)
}
