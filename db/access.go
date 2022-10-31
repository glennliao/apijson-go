package db

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var accessMap = map[string]Access{}

type Access struct {
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

	// ext

	RowKey string
}

func loadAccessMap() {
	_accessMap := make(map[string]Access)

	var accessList []Access
	g.DB().Model(config.TableAccess).Scan(&accessList)

	for _, access := range accessList {
		name := access.Alias
		if name == "" {
			name = access.Name
		}
		_accessMap[name] = access
	}

	accessMap = _accessMap
}

func GetAccess(table string, accessVerify bool) (*Access, error) {
	// 暂未使用version
	// 读取配置时将最新的版本额外增加一个@latest的版本, 传入为-1时候, 读取最新版本
	access, ok := accessMap[table]

	if !ok {
		if accessVerify {
			return nil, gerror.Newf("access[%s]: 404", table)
		}
		return &Access{
			Debug: 0,
			Name:  table,
			Alias: table,
		}, nil
	}

	return &access, nil
}

func GetAccessRole(table string, method string) ([]string, string, error) {
	// 暂未使用version
	// 读取配置时将最新的版本额外增加一个@latest的版本, 传入为-1时候, 读取最新版本
	access, ok := accessMap[table]

	if !ok {
		return nil, "", gerror.Newf("access[%s]: 404", table)
	}

	switch method {
	case consts.MethodGet:
		return access.Get, access.Name, nil
	case consts.MethodHead:
		return access.Head, access.Name, nil
	case consts.MethodPost:
		return access.Post, access.Name, nil
	case consts.MethodPut:
		return access.Put, access.Name, nil
	case consts.MethodDelete:
		return access.Delete, access.Name, nil
	}

	return []string{}, access.Name, nil
}

func Init() {
	Reload()
}

// Reload 重载刷新配置
func Reload() {
	loadAccessMap()
	loadRequestMap()
	loadTableMeta()
}
