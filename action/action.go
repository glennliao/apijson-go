package action

import (
	"context"
	"strings"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// Action 非get查询的request表中的请求
type Action struct {
	ctx        context.Context
	tagRequest *config.RequestConfig
	method     string

	req model.Map

	err error

	children map[string]*Node
	keyNode  map[string]*Node

	// 关闭 access 权限验证, 默认否
	NoAccessVerify bool
	// 关闭 request 验证开关, 默认否
	NoRequestVerify bool

	// Access *config.Access

	// dbFieldStyle 数据库字段命名风格 请求传递到数据库中
	DbFieldStyle config.FieldStyle

	// jsonFieldStyle 数据库返回的字段
	JsonFieldStyle config.FieldStyle

	ActionConfig *config.ActionConfig
}

func New(ctx context.Context, actionConfig *config.ActionConfig, method string, req model.Map) *Action {

	request, err := checkTag(req, method, actionConfig)
	if err != nil {
		panic(err)
	}

	delete(req, consts.Tag)
	delete(req, "version")

	a := &Action{
		ctx:          ctx,
		tagRequest:   request,
		method:       method,
		req:          req,
		children:     map[string]*Node{},
		keyNode:      map[string]*Node{},
		ActionConfig: actionConfig,
	}
	return a
}

func (a *Action) parse() error {

	structures := a.tagRequest.Structure

	for key, v := range a.req {

		structuresKey := key
		if strings.HasSuffix(key, consts.ListKeySuffix) {
			structuresKey = structuresKey[0 : len(structuresKey)-2]
		}
		structure, ok := structures[key]
		if !ok {
			if structure, ok = structures[structuresKey]; !ok { // User[]可读取User或者User[]
				return gerror.New("structure错误: 400, 缺少" + key)
			}
		}

		var list []model.Map
		_v, ok := v.(model.Map)
		if ok { // 将所有node都假设成列表, 如果单个则看成一个元素的批量
			list = []model.Map{_v}
		} else {
			for _, m := range gconv.Maps(v) {
				list = append(list, m)
			}
		}

		node := newNode(key, list, structure, a.tagRequest.Executor[key])
		node.ctx = a.ctx
		node.Action = a
		a.keyNode[key] = &node
		node.keyNode = a.keyNode
		err := node.parse(a.ctx, a.method)
		if err != nil {
			return err
		}

		a.children[key] = &node
	}

	return nil
}

func (a *Action) Result() (model.Map, error) {

	err := a.parse()
	if err != nil {
		return nil, err
	}

	ret := model.Map{}

	for _, k := range a.tagRequest.ExecQueue {
		node := a.children[k]
		err = EmitHook(a.ctx, BeforeNodeExec, node, a.method)
		if err != nil {
			return nil, err
		}
	}

	for _, k := range a.tagRequest.ExecQueue {

		node := a.children[k]
		err = node.reqUpdate()
		if err != nil {
			return nil, err
		}
	}

	err = g.DB().Transaction(a.ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, k := range a.tagRequest.ExecQueue {
			node := a.children[k]
			ret[k], err = node.execute(ctx, a.method)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, k := range a.tagRequest.ExecQueue {
		node := a.children[k]
		err = EmitHook(a.ctx, AfterNodeExec, node, a.method)
		if err != nil {
			return nil, err
		}
	}

	return ret, err
}

func checkTag(req model.Map, method string, requestCfg *config.ActionConfig) (*config.RequestConfig, error) {
	_tag, ok := req[consts.Tag]
	if !ok {
		return nil, gerror.New("tag 缺失")
	}

	tag := gconv.String(_tag)
	version := req["version"]

	request, err := requestCfg.GetRequest(tag, method, gconv.String(version))
	if err != nil {
		return nil, err
	}

	return request, nil
}
