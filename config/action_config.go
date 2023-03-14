package config

import (
	"context"
	"github.com/glennliao/apijson-go/model"
)

type ActionConfig struct {
	requestConfig    *RequestConfig
	access           *Access
	functions        *functions
	rowKeyGenFuncMap map[string]RowKeyGenFuncHandler
	defaultRoleFunc  DefaultRole
}

func (c *ActionConfig) NoVerify() bool {
	return c.access.NoVerify
}

func (c *ActionConfig) DefaultRoleFunc() DefaultRole {
	return c.defaultRoleFunc
}

func (c *ActionConfig) GetAccessConfig(key string, noVerify bool) (*AccessConfig, error) {
	return c.access.GetAccess(key, noVerify)
}

func (c *ActionConfig) CallFunc(ctx context.Context, name string, param model.Map) (any, error) {
	return c.functions.Call(ctx, name, param)
}

func (c *ActionConfig) GetRequest(tag string, method string, version string) (*Request, error) {
	return c.requestConfig.GetRequest(tag, method, version)
}

func (c *ActionConfig) ConditionFunc(ctx context.Context, req ConditionReq, condition *ConditionRet) error {
	return c.access.ConditionFunc(ctx, req, condition)
}

func (c *ActionConfig) RowKeyGen(ctx context.Context, genFuncName string, accessName string, data model.Map) (model.Map, error) {
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
