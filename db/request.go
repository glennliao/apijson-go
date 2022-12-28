package db

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

var requestMap = map[string]*Request{}

type Request struct {
	Debug       int8
	Version     int16
	Method      string
	Tag         string
	StructureDb map[string]any        `orm:"structure"`
	Structure   map[string]*Structure `orm:"-"`
	Detail      string
	CreatedAt   *gtime.Time
	// 节点执行顺序
	ExecQueue []string
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

func loadRequestMap() {
	_requestMap := make(map[string]*Request)

	var requestList []Request
	err := g.DB().Model(config.TableRequest).OrderAsc("version").Scan(&requestList)
	if err != nil {
		panic(err)
	}

	for _, _item := range requestList {
		item := _item
		tag, _ := getTag(item.Tag)

		if strings.ToLower(tag) != tag {
			// 本身大写, 如果没有外层, 则套一层
			if _, ok := item.StructureDb[tag]; !ok {
				item.StructureDb = map[string]any{
					tag: item.StructureDb,
				}
			}
		}

		item.Structure = make(map[string]*Structure)
		for k, v := range item.StructureDb {
			structure := Structure{}
			err = gconv.Scan(v, &structure)
			if err != nil {
				panic(err)
			}

			if structure.Must != nil {
				structure.Must = strings.Split(structure.Must[0], ",")
			}
			if structure.Refuse != nil {
				structure.Refuse = strings.Split(structure.Refuse[0], ",")
			}

			item.Structure[k] = &structure
		}

		if item.ExecQueue != nil {
			item.ExecQueue = strings.Split(item.ExecQueue[0], ",")
		} else {
			item.ExecQueue = []string{tag}
		}

		_requestMap[getRequestFullKey(item.Tag, item.Method, gconv.String(item.Version))] = &item
		//  获取时version排序,所以此处最后一个为最新
		_requestMap[getRequestFullKey(item.Tag, item.Method, "latest")] = &item
	}

	requestMap = _requestMap
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

func getRequestFullKey(tag string, method string, version string) string {
	return tag + "@" + method + "@" + version
}

func GetRequest(tag string, method string, version string) (*Request, error) {

	if version == "" || version == "-1" || version == "0" {
		version = "latest"
	}

	key := getRequestFullKey(tag, method, version)
	request, ok := requestMap[key]

	if !ok {
		return nil, gerror.Newf("request[%s]: 404", key)
	}

	return request, nil
}
