package config

import (
	"context"

	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type RowKeyGenReq struct {
	AccessName string
	TableName  string

	GenParam model.FuncParam

	Data model.Map
}

type RowKeyGenFuncHandler func(ctx context.Context, req *RowKeyGenReq, ret *RowKeyGenRet) error

type RowKeyGenerator struct {
	Name      string
	ParamList []ParamItem
	Handler   RowKeyGenFuncHandler
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

// func (c *Config) RowKeyGen(ctx context.Context, genFuncName string, accessName string, data model.Map) (model.Map, error) {
// 	if f, exists := c.rowKeyGenFuncMap[genFuncName]; exists {
// 		req := &RowKeyGenReq{
// 			AccessName: accessName,
// 			Data:       data,
// 		}
// 		ret := NewRowKeyGenRet()
// 		err := f.Handler(ctx, req, ret)
// 		return ret.data, err
// 	}
//
// 	return nil, nil
// }

func (c *Config) RowKeyGenFunc(f RowKeyGenerator) {
	if f.Handler == nil {
		panic("RowKeyGenFunc Handler is nil")
	}
	if _, ok := c.rowKeyGenFuncMap[f.Name]; ok {
		panic("RowKeyGenFunc " + f.Name + " already exists")
	}
	c.rowKeyGenFuncMap[f.Name] = f
}

func (c *Config) RowKeyGenList() []string {
	return lo.Keys(c.rowKeyGenFuncMap)
}
