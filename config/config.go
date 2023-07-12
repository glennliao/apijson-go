package config

import (
	"context"
)

type AccessListProvider func(ctx context.Context) []AccessConfig

var accessListProviderMap = make(map[string]AccessListProvider)

func RegAccessListProvider(name string, provider AccessListProvider) {
	accessListProviderMap[name] = provider
}

type RequestListProvider func(ctx context.Context) []RequestConfig

var requestListProviderMap = make(map[string]RequestListProvider)

func RegRequestListProvider(name string, provider RequestListProvider) {
	requestListProviderMap[name] = provider
}

type DbMetaProvider func(ctx context.Context) []Table

var dbMetaProviderMap = make(map[string]DbMetaProvider)

func RegDbMetaProvider(name string, provider DbMetaProvider) {
	dbMetaProviderMap[name] = provider
}

type Config struct {
	Access *Access

	Functions *functions

	MaxTreeWidth int
	MaxTreeDeep  int

	rowKeyGenFuncMap map[string]RowKeyGenerator

	// dbFieldStyle 数据库字段命名风格 请求传递到数据库中
	DbFieldStyle FieldStyle

	// jsonFieldStyle 数据库返回的字段
	JsonFieldStyle FieldStyle

	DbMeta *DBMeta

	AccessListProvider  string
	RequestListProvider string
	DbMetaProvider      string

	accessList []AccessConfig

	requestConfigs *RequestConfigs
	queryConfig    *QueryConfig
	actionConfig   *ActionConfig
}

func New() *Config {
	c := &Config{}
	c.Access = NewAccess()
	c.AccessListProvider = "db"
	c.RequestListProvider = "db"
	c.DbMetaProvider = "db"

	c.MaxTreeWidth = 5
	c.MaxTreeDeep = 5

	c.rowKeyGenFuncMap = make(map[string]RowKeyGenerator)

	c.DbFieldStyle = CaseSnake
	c.JsonFieldStyle = CaseCamel

	c.Functions = &functions{}
	c.Functions.funcMap = make(map[string]*Func)

	return c
}

func (c *Config) ReLoad() {

	accessConfigMap := make(map[string]AccessConfig)

	ctx := context.Background()

	accessListProvider := accessListProviderMap[c.AccessListProvider]

	if accessListProvider != nil {
		c.accessList = accessListProvider(ctx)

		defaultMaxCount := 100

		for _, access := range c.accessList {
			name := access.Alias
			if name == "" {
				name = access.Name
			}

			if access.FieldsGet == nil {
				access.FieldsGet = map[string]*FieldsGetValue{}
			}

			if _, exists := access.FieldsGet["default"]; !exists {
				access.FieldsGet["default"] = &FieldsGetValue{}
			}

			for role, _ := range access.FieldsGet {
				if access.FieldsGet[role].MaxCount == nil {
					access.FieldsGet[role].MaxCount = &defaultMaxCount
				}
			}
			accessConfigMap[access.Alias] = access
		}
	}

	c.Access.accessConfigMap = accessConfigMap

	requestListProvider := requestListProviderMap[c.RequestListProvider]
	if requestListProvider != nil {
		requestList := requestListProvider(ctx)
		c.requestConfigs = NewRequestConfig(requestList)
	}

	dbMetaProvider := dbMetaProviderMap[c.DbMetaProvider]
	if dbMetaProvider != nil {
		c.DbMeta = NewDbMeta(dbMetaProvider(ctx))
	}

	c.queryConfig = &QueryConfig{
		access:          c.Access,
		functions:       c.Functions,
		maxTreeDeep:     c.MaxTreeDeep,
		maxTreeWidth:    c.MaxTreeWidth,
		defaultRoleFunc: c.Access.DefaultRoleFunc,
	}

	c.actionConfig = &ActionConfig{
		requestConfig:    c.requestConfigs,
		access:           c.Access,
		functions:        c.Functions,
		rowKeyGenFuncMap: c.rowKeyGenFuncMap,
		defaultRoleFunc:  c.Access.DefaultRoleFunc,
	}

}

func (c *Config) QueryConfig() *QueryConfig {
	return c.queryConfig
}

func (c *Config) ActionConfig() *ActionConfig {
	return c.actionConfig
}
