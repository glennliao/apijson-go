package action

import (
	"context"
	"github.com/glennliao/apijson-go/db"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

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

// Action 非get查询的request表中的请求
type Action struct {
	ctx        context.Context
	tagRequest *db.Request
	method     string

	req g.Map

	err error

	children map[string]Node
	keyNode  map[string]*Node
}

func New(ctx context.Context, method string, req g.Map) *Action {

	request, err := checkTag(req, method)
	if err != nil {
		panic(err)
	}

	delete(req, "tag")

	a := &Action{
		ctx:        ctx,
		tagRequest: request,
		method:     method,
		req:        req,
		children:   map[string]Node{},
		keyNode:    map[string]*Node{},
	}
	return a
}

func (a *Action) parse() error {

	structures := a.tagRequest.Structure

	for key, v := range a.req {

		structuresKey := key
		if strings.HasSuffix(key, "[]") {
			structuresKey = structuresKey[0 : len(structuresKey)-2]
		}
		structureMap, ok := structures[key]
		if !ok {
			if structureMap, ok = structures[structuresKey]; !ok { //User[]可读取User或者User[]
				return gerror.New("structure错误: 400, 缺少" + key)
			}
		}

		structure := Structure{}
		err := gconv.Scan(structureMap, &structure)
		if err != nil {
			return err
		}

		// todo 初始化时完成map2struct,不用每次都scan生成
		structure.Must = strings.Split(structure.Must[0], ",")
		structure.Refuse = strings.Split(structure.Refuse[0], ",")

		var list []g.Map
		_v, ok := v.(g.Map)
		if ok { // 将所有node都假设成列表, 如果单个则看成一个元素的批量
			list = []g.Map{_v}
		} else {
			list = v.([]g.Map)
		}

		node := newNode(key, list, structure)
		node.ctx = a.ctx
		a.keyNode[key] = &node
		node.keyNode = a.keyNode
		err = node.parse(a.ctx, a.method)
		if err != nil {
			return err
		}

		a.children[key] = node
	}

	return nil
}

func (a *Action) Result() (g.Map, error) {

	err := a.parse()
	if err != nil {
		return nil, err
	}

	ret := g.Map{}

	err = g.DB().Transaction(a.ctx, func(ctx context.Context, tx *gdb.TX) error {
		for _, k := range a.tagRequest.ExecQueue {

			node := a.children[k]
			ret[k], err = node.execute(ctx, a.method)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return ret, err
}
