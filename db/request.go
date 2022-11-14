package db

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
)

var requestMap = map[string]Request{}

type Request struct {
	Debug     int8
	Version   int16
	Method    string
	Tag       string
	Structure g.Map
	Detail    string
	CreatedAt *gtime.Time

	ExecQueue []string
}

func loadRequestMap() {
	_requestMap := make(map[string]Request)

	var requestList []Request
	g.DB().Model(config.TableRequest).Scan(&requestList)

	for _, item := range requestList {

		tag := item.Tag
		if strings.HasSuffix(tag, "[]") {
			tag = tag[0 : len(tag)-2]
		}
		if strings.ToLower(tag) != tag {
			// 本身大写, 如果没有外层, 则套一层
			if _, ok := item.Structure[tag]; !ok {
				item.Structure = g.Map{
					tag: item.Structure,
				}
			}
		}

		// todo 改成列表读取数据库, 避免多次查询
		type ext struct {
			ExecQueue string
		}
		var _ext *ext
		g.DB().Model(config.TableRequestExt).Where(g.Map{
			"version": item.Version,
			"method":  item.Method,
			"tag":     item.Tag,
		}).Scan(&_ext)

		if _ext != nil {
			item.ExecQueue = strings.Split(_ext.ExecQueue, ",")
		} else {
			tag := item.Tag
			if strings.HasSuffix(tag, "[]") {
				tag = tag[0 : len(tag)-2]
			}
			item.ExecQueue = strings.Split(tag, ",")
		}

		_requestMap[item.Method+"@"+item.Tag] = item
	}

	requestMap = _requestMap
}

func GetRequest(tag string, method string, version int16) (*Request, error) {
	// 暂未使用version
	// 读取配置时将最新的版本额外增加一个@latest的版本, 传入为-1时候, 读取最新版本
	key := method + "@" + tag
	request, ok := requestMap[key]

	if !ok {
		return nil, gerror.Newf("request[%s]: 404", key)
	}

	return &request, nil
}
