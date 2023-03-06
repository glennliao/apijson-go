package config

import (
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/samber/lo"
	"net/http"
)

type FieldsGetValue struct {
	In  map[string][]string
	Out map[string]string
}

type AccessConfig struct {
	Debug     int8
	Name      string
	Alias     string
	Get       []string
	Head      []string
	Gets      []string
	Heads     []string
	Post      []string
	Put       []string
	Delete    []string
	CreatedAt *gtime.Time
	Detail    string

	RowKeyGen string // 主键生成策略
	RowKey    string
	FieldsGet map[string]FieldsGetValue
	Executor  string
}

func (a *AccessConfig) GetFieldsGetOutByRole(role string) []string {
	var fieldsMap map[string]string

	if val, exists := a.FieldsGet[role]; exists {
		fieldsMap = val.Out
	} else {
		fieldsMap = a.FieldsGet["default"].Out
	}
	return lo.Keys(fieldsMap)
}

func (a *AccessConfig) GetFieldsGetInByRole(role string) map[string][]string {
	var inFieldsMap map[string][]string

	if val, exists := a.FieldsGet[role]; exists {
		inFieldsMap = val.In
	} else {
		inFieldsMap = a.FieldsGet["default"].In
	}

	return inFieldsMap
}

func (a *Access) GetAccess(tableAlias string, accessVerify bool) (*AccessConfig, error) {
	tableAlias, _ = util.ParseNodeKey(tableAlias)
	access, ok := a.accessConfigMap[tableAlias]

	if !ok {
		if accessVerify {
			return nil, gerror.Newf("access[%s]: 404", tableAlias)
		}
		return &AccessConfig{
			Debug: 0,
			Name:  tableAlias,
			Alias: tableAlias,
		}, nil
	}

	return &access, nil
}

func (a *Access) GetAccessRole(table string, method string) ([]string, string, error) {
	access, ok := a.accessConfigMap[table]

	if !ok {
		return nil, "", gerror.Newf("access[%s]: 404", table)
	}

	switch method {
	case http.MethodGet:
		return access.Get, access.Name, nil
	case http.MethodHead:
		return access.Head, access.Name, nil
	case http.MethodPost:
		return access.Post, access.Name, nil
	case http.MethodPut:
		return access.Put, access.Name, nil
	case http.MethodDelete:
		return access.Delete, access.Name, nil
	}

	return []string{}, access.Name, nil
}
