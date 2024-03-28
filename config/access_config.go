package config

import (
	"net/http"

	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/samber/lo"
)

type FieldsGetValue struct {
	In       map[string][]string
	Out      map[string]string
	MaxCount *int // 可使用的最大分页大小,默认100
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
	FieldsGet map[string]*FieldsGetValue
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

func (a *AccessConfig) GetAccessRoles(method string) []string {
	switch method {
	case http.MethodGet:
		return a.Get
	case http.MethodHead:
		return a.Head
	case http.MethodPost:
		return a.Post
	case http.MethodPut:
		return a.Put
	case http.MethodDelete:
		return a.Delete
	}

	return make([]string, 0)
}

func (a *Access) GetAccess(accessName string, noVerify bool) (*AccessConfig, error) {
	access, ok := a.accessConfigMap[accessName]

	if !ok {
		if noVerify {
			return &AccessConfig{
				Debug: 0,
				Name:  accessName,
				Alias: accessName,
			}, nil
		}
		return nil, consts.NewAccessNoFoundErr(accessName)
	}

	return &access, nil
}

func (a *Access) GetAccessRole(accessName string, method string) ([]string, string, error) {
	access, ok := a.accessConfigMap[accessName]

	if !ok {
		return nil, "", consts.NewAccessNoFoundErr(accessName)
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
