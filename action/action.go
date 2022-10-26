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

// Structure https://github.com/Tencent/APIJSON/blob/master/APIJSONORM/src/main/java/apijson/orm/Operation.java
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
	}
	return a
}

func (a *Action) parse() error {

	structures := a.tagRequest.Structure

	for key, v := range a.req {

		structureMap, ok := structures[key]
		if !ok {
			return gerror.New("structure错误: 400, 缺少" + key)
		}

		structure := Structure{}
		err := gconv.Scan(structureMap, &structure)
		if err != nil {
			return err
		}

		structure.Must = strings.Split(structure.Must[0], ",")
		structure.Refuse = strings.Split(structure.Refuse[0], ",")

		node := newNode(key, v.(g.Map), structure)
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
		for k, node := range a.children {
			ret[k], err = node.execute(ctx, a.method)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return ret, err
}
