package config

import (
	"context"
	"github.com/glennliao/apijson-go/model"
)

type ActionConfig struct {
	requestConfig    RequestConfig
	access           *Access
	functions        *Functions
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

//
//type ExecutorConfig struct {
//	NoVerify       bool
//	accessConfig   *AccessConfig
//	method         string
//	role           string
//	DBMeta         *DBMeta
//	DbFieldStyle   FieldStyle
//	JsonFieldStyle FieldStyle
//}
//
//func NewExecutorConfig(accessConfig *AccessConfig, method string, noVerify bool) *ExecutorConfig {
//	return &ExecutorConfig{
//		accessConfig: accessConfig,
//		method:       method,
//		NoVerify:     noVerify,
//	}
//}
//
//func (c *ExecutorConfig) SetRole(role string) {
//	c.role = role
//}
//
//func (c *ExecutorConfig) TableName() string {
//	return c.accessConfig.Name
//}
//
//func (c *ExecutorConfig) TableColumns() []string {
//	return c.DBMeta.GetTableColumns(c.accessConfig.Name)
//}
//
//func (c *ExecutorConfig) GetFieldsGetOutByRole() []string {
//	var fieldsMap map[string]string
//
//	if val, exists := c.accessConfig.FieldsGet[c.role]; exists {
//		fieldsMap = val.Out
//	} else {
//		fieldsMap = c.accessConfig.FieldsGet["default"].Out
//	}
//	return lo.Keys(fieldsMap)
//}
//
//func (c *ExecutorConfig) GetFieldsGetInByRole() map[string][]string {
//	var inFieldsMap map[string][]string
//
//	if val, exists := c.accessConfig.FieldsGet[c.role]; exists {
//		inFieldsMap = val.In
//	} else {
//		inFieldsMap = c.accessConfig.FieldsGet["default"].In
//	}
//
//	return inFieldsMap
//}
//
//func (c *ExecutorConfig) AccessRoles() []string {
//	switch c.method {
//	case http.MethodGet:
//		return c.accessConfig.Get
//	case http.MethodHead:
//		return c.accessConfig.Head
//	case http.MethodPost:
//		return c.accessConfig.Post
//	case http.MethodPut:
//		return c.accessConfig.Put
//	case http.MethodDelete:
//		return c.accessConfig.Delete
//	}
//	return []string{}
//
//}
//
//func (c *ExecutorConfig) Executor() string {
//	return c.accessConfig.Executor
//
//}
