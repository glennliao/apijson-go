package db

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
)

var RowKeyMap = map[string]string{
	"user": "userId",
	"todo": "id",
} // 从数据库中读入

var AccessMap = map[string]Access{}

type Access struct {
	Debug  int8
	Name   string
	Alias  string
	Get    []string
	Head   []string
	Gets   []string
	Heads  []string
	Post   []string
	Put    []string
	Delete []string
	Date   *gtime.Time
	Detail string
}

var RequestMap = map[string]Request{}

type Request struct {
	Debug   int8
	Version int16
	Method  string
	Tag     string
	// https://github.com/Tencent/APIJSON/blob/master/APIJSONORM/src/main/java/apijson/orm/Operation.java
	Structure g.Map
	Detail    string
	Date      *gtime.Time
}

func Init() {
	getAccessMap()
	getRequestMap()
}

func getAccessMap() {
	accessMap := make(map[string]Access)

	var accessList []Access
	g.DB().Model("Access").Scan(&accessList)

	for _, access := range accessList {
		accessMap[access.Name] = access
	}

	AccessMap = accessMap
}

func getRequestMap() {
	requestMap := make(map[string]Request)

	var requestList []Request
	g.DB().Model("Request").Scan(&requestList)

	for _, item := range requestList {

		tag := item.Tag
		if strings.ToLower(tag) != tag {
			// 本身大写, 如果没有外层, 则套一层
			// https://github.com/Tencent/APIJSON/issues/115#issuecomment-565733254
			if _, ok := item.Structure[tag]; !ok {
				item.Structure = g.Map{
					tag: item.Structure,
				}
			}
		}

		requestMap[item.Method+"@"+item.Tag] = item
	}

	RequestMap = requestMap
}
