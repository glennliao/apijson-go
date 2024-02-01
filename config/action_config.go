package config

import (
	"context"
	"strings"

	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/container/gvar"
)

type ActionConfig struct {
	requestConfig    *RequestConfigs
	access           *Access
	functions        *functions
	rowKeyGenFuncMap map[string]RowKeyGenerator
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

func (c *ActionConfig) Func(name string) *Func {
	return c.functions.funcMap[name]
}

func (c *ActionConfig) CallFunc(ctx context.Context, name string, param model.Map) (res any, err error) {
	return c.functions.Call(ctx, name, param)
}

func (c *ActionConfig) GetRequest(tag string, method string, version string) (*RequestConfig, error) {
	return c.requestConfig.GetRequest(tag, method, version)
}

func (c *ActionConfig) ConditionFunc(ctx context.Context, req ConditionReq, condition *ConditionRet) error {
	return c.access.ConditionFunc(ctx, req, condition)
}

func (c *ActionConfig) RowKeyGen(ctx context.Context, genFuncName string, accessName string, tableName string, data model.Map) (model.Map, error) {
	var paramKeys []string

	if strings.Contains(genFuncName, "(") {
		genFuncName, paramKeys = util.ParseFunctionsStr(genFuncName)
	}

	if f, exists := c.rowKeyGenFuncMap[genFuncName]; exists {

		req := &RowKeyGenReq{
			AccessName: accessName,
			TableName:  tableName,
			Data:       data,
		}

		if len(paramKeys) > 0 {
			param := model.FuncParam{}
			for i, item := range f.ParamList {
				if len(paramKeys) >= i {
					param[item.Name] = gvar.New(paramKeys[i])
				} else {
					param[item.Name] = gvar.New(item.Default)
				}
			}
			req.GenParam = param
		}

		ret := NewRowKeyGenRet()
		err := f.Handler(ctx, req, ret)
		return ret.data, err
	}

	return nil, consts.NewSysErr("rowKey RowKeyGenerator not found：" + genFuncName)
}
