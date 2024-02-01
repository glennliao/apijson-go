package action

import (
	"context"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
)

type Action struct {
	ctx context.Context
	err error
	req model.Map
	ret model.Map

	Method           string
	tagRequestConfig *config.RequestConfig

	children map[string]*Node
	keyNode  map[string]*Node

	// 关闭 access 权限验证, 默认否
	NoAccessVerify bool
	// 关闭 request 验证开关, 默认否
	NoRequestVerify bool

	// dbFieldStyle 数据库字段命名风格 请求传递到数据库中
	DbFieldStyle config.FieldStyle
	// jsonFieldStyle 数据库返回的字段
	JsonFieldStyle config.FieldStyle

	ActionConfig *config.ActionConfig

	HooksMap map[string][]*Hook
}

func New(ctx context.Context, actionConfig *config.ActionConfig, method string, req model.Map) *Action {
	a := &Action{
		ctx:          ctx,
		req:          req,
		Method:       method,
		ActionConfig: actionConfig,
	}

	request, err := CheckTag(req, method, actionConfig)
	if err != nil {
		a.err = err
		return a
	}

	a.tagRequestConfig = request

	return a
}

func (a *Action) Result() (model.Map, error) {
	if a.err != nil {
		return nil, a.err
	}

	a.parse()
	if a.err != nil {
		return nil, a.err
	}

	a.exec()
	return a.ret, a.err
}
