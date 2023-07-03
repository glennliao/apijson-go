package config

import (
	"strings"

	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type RequestConfig struct {
	Debug     int8
	Version   string
	Method    string
	Tag       string
	Structure map[string]*Structure
	Detail    string
	CreatedAt *gtime.Time
	// 节点执行顺序
	ExecQueue []string
	Executor  map[string]string
	// 是否开启事务
	Transaction *bool
}

type Structure struct {
	Must   []string `json:"MUST,omitempty"`
	Refuse []string `json:"REFUSE,omitempty"`

	Unique []string `json:"UNIQUE,omitempty"`

	// 不存在时添加
	Insert g.Map `json:"INSERT,omitempty"`
	// 不存在时就添加，存在时就修改
	Update g.Map `json:"UPDATE,omitempty"`
	// 存在时替换
	Replace g.Map `json:"REPLACE,omitempty"`
	// 存在时移除
	Remove []string `json:"REMOVE,omitempty"`
}

type RequestConfigs struct {
	requestMap map[string]*RequestConfig
}

func NewRequestConfig(requestList []RequestConfig) *RequestConfigs {
	c := RequestConfigs{}
	requestMap := make(map[string]*RequestConfig)

	for _, _item := range requestList {
		item := _item

		if item.Structure == nil {
			item.Structure = make(map[string]*Structure)
		}

		if len(item.ExecQueue) == 0 {
			tag, _ := getTag(item.Tag)
			item.ExecQueue = []string{tag}
		}

		requestMap[getRequestFullKey(item.Tag, item.Method, gconv.String(item.Version))] = &item
		//  获取时version排序,所以此处最后一个为最新
		requestMap[getRequestFullKey(item.Tag, item.Method, "latest")] = &item
	}

	c.requestMap = requestMap
	return &c
}

func getRequestFullKey(tag string, method string, version string) string {
	return tag + "@" + method + "@" + version
}

func (c *RequestConfigs) GetRequest(tag string, method string, version string) (*RequestConfig, error) {

	if version == "" || version == "-1" || version == "0" {
		version = "latest"
	}

	key := getRequestFullKey(tag, method, version)
	request, ok := c.requestMap[key]

	if !ok {
		return nil, gerror.Newf("request[%s]: 404", key)
	}

	return request, nil
}

func getTag(tag string) (name string, isList bool) {
	if strings.HasSuffix(tag, consts.ListKeySuffix) {
		name = tag[0 : len(tag)-2]
		isList = true
	} else {
		name = tag
	}

	return
}
