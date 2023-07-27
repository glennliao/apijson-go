package action

import (
	"context"
	"strings"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/query"
	"github.com/glennliao/apijson-go/util"
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

	NewQuery  func(ctx context.Context, req model.Map) *query.Query
	NewAction func(ctx context.Context, method string, req model.Map) *Action

	HooksMap map[string][]*Hook

	ret model.Map
}

func New(ctx context.Context, actionConfig *config.ActionConfig, method string, req model.Map) *Action {

	request, err := checkTag(req, method, actionConfig)
	if err != nil {
		return &Action{
			err: err,
		}
	}

	delete(req, consts.Tag)
	delete(req, consts.Version)

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

	if a.err != nil {
		return a.err
	}

	structures := a.tagRequest.Structure

	for key, v := range a.req {

		structuresKey := key
		if strings.HasSuffix(key, consts.ListKeySuffix) {
			structuresKey = util.RemoveSuffix(key, consts.ListKeySuffix)
		}

		structure, ok := structures[key]
		if !ok {
			if structure, ok = structures[structuresKey]; !ok { // User[]可读取User或者User[]
				return consts.NewStructureKeyNoFoundErr(key)
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

func (a *Action) hookExecute(i int, inTransaction bool) error {
	k := a.tagRequest.ExecQueue[i]
	node := a.children[k]
	nodeHookReq := &HookReq{
		Node:            node,
		Method:          a.method,
		ctx:             a.ctx,
		nextIdx:         -1,
		isInTransaction: inTransaction,
		hooks:           getHooksByAccessName(a.HooksMap, k),
	}

	nodeHookReq.handler = func(ctx context.Context, n *Node, method string) error {

		if i+1 < len(a.tagRequest.ExecQueue) {
			return a.hookExecute(i+1, inTransaction)
		}

		// 执行完了普通hook的before,开始执行事务内
		if !inTransaction {

			transactionHandler := noTransactionHandler

			if a.tagRequest.Transaction != nil && *a.tagRequest.Transaction == true {
				h := GetTransactionHandler(a.ctx, a)
				if h == nil {
					err := consts.NewSysErr("transaction handler is nil")
					return err
				}
				transactionHandler = h
			}

			err := transactionHandler(a.ctx, func(ctx context.Context) error {
				return a.hookExecute(0, !inTransaction)
			})

			return err
		}

		var err error
		a.ret[k], err = node.execute(ctx, a.method)
		return err
	}

	err := nodeHookReq.Next()
	return err
}

func (a *Action) Result() (model.Map, error) {

	err := a.parse()
	if err != nil {
		return nil, err
	}

	a.ret = model.Map{}

	err = a.hookExecute(0, false)
	if err != nil {
		a.err = err
	}
	return a.ret, err
}

func checkTag(req model.Map, method string, requestCfg *config.ActionConfig) (*config.RequestConfig, error) {
	_tag, ok := req[consts.Tag]
	if !ok {
		return nil, consts.ErrNoTag
	}

	tag := gconv.String(_tag)
	version := req[consts.Version]

	request, err := requestCfg.GetRequest(tag, method, gconv.String(version))
	if err != nil {
		return nil, err
	}

	return request, nil
}
