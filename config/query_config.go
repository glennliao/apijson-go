package config

import (
	"context"
	"net/http"

	"github.com/glennliao/apijson-go/model"
	"github.com/samber/lo"
)

type QueryConfig struct {
	access          *Access
	functions       *functions
	maxTreeDeep     int
	maxTreeWidth    int
	defaultRoleFunc DefaultRole
}

func (c *QueryConfig) NoVerify() bool {
	return c.access.NoVerify
}

func (c *QueryConfig) DefaultRoleFunc() DefaultRole {
	return c.defaultRoleFunc
}

func (c *QueryConfig) GetAccessConfig(key string, noVerify bool) (*AccessConfig, error) {
	return c.access.GetAccess(key, noVerify)
}

func (c *QueryConfig) Func(name string) *Func {
	return c.functions.funcMap[name]
}

func (c *QueryConfig) CallFunc(ctx context.Context, name string, param model.Map) (res any, err error) {
	return c.functions.Call(ctx, name, param)
}

func (c *QueryConfig) MaxTreeDeep() int {
	return c.maxTreeDeep
}

func (c *QueryConfig) MaxTreeWidth() int {
	return c.maxTreeWidth
}

type ExecutorConfig struct {
	NoVerify       bool
	accessConfig   *AccessConfig
	method         string
	role           string
	DBMeta         *DBMeta
	DbFieldStyle   FieldStyle
	JsonFieldStyle FieldStyle
}

func NewExecutorConfig(accessConfig *AccessConfig, method string, noVerify bool) *ExecutorConfig {
	return &ExecutorConfig{
		accessConfig: accessConfig,
		method:       method,
		NoVerify:     noVerify,
	}
}

func (c *ExecutorConfig) SetRole(role string) {
	c.role = role
}

func (c *ExecutorConfig) TableName() string {
	return c.accessConfig.Name
}

func (c *ExecutorConfig) TableColumns() []string {
	return c.DBMeta.GetTableColumns(c.accessConfig.Name)
}

func (c *ExecutorConfig) GetFieldsGetByRole() *FieldsGetValue {
	if val, exists := c.accessConfig.FieldsGet[c.role]; exists {
		return val
	}

	return c.accessConfig.FieldsGet["default"]
}

func (c *ExecutorConfig) GetFieldsGetOutByRole() []string {
	if c.accessConfig.FieldsGet == nil {
		return make([]string, 0)
	}

	var fieldsMap map[string]string

	if val, exists := c.accessConfig.FieldsGet[c.role]; exists {
		fieldsMap = val.Out
	} else {
		fieldsMap = c.accessConfig.FieldsGet["default"].Out
	}
	return lo.Keys(fieldsMap)
}

func (c *ExecutorConfig) GetFieldsGetInByRole() map[string][]string {
	var inFieldsMap map[string][]string

	if val, exists := c.accessConfig.FieldsGet[c.role]; exists {
		inFieldsMap = val.In
	} else {
		inFieldsMap = c.accessConfig.FieldsGet["default"].In
	}

	return inFieldsMap
}

func (c *ExecutorConfig) AccessRoles() []string {
	switch c.method {
	case http.MethodGet:
		return c.accessConfig.Get
	case http.MethodHead:
		return c.accessConfig.Head
	case http.MethodPost:
		return c.accessConfig.Post
	case http.MethodPut:
		return c.accessConfig.Put
	case http.MethodDelete:
		return c.accessConfig.Delete
	}
	return []string{}
}

func (c *ExecutorConfig) Executor() string {
	return c.accessConfig.Executor
}
