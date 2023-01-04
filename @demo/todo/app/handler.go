package app

import (
	"context"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/query"
	"github.com/gogf/gf/v2/frame/g"
	"net/http"
)

func Get(ctx context.Context, req g.Map) (res g.Map, err error) {
	q := query.New(ctx, req)
	q.AccessVerify = config.AccessVerify
	q.AccessCondition = config.AccessConditionFunc
	return q.Result()
}

func Head(ctx context.Context, req g.Map) (res g.Map, err error) {
	return nil, err
}

func Post(ctx context.Context, req g.Map) (res g.Map, err error) {
	act := action.New(ctx, http.MethodPost, req)
	return act.Result()
}

func Put(ctx context.Context, req g.Map) (res g.Map, err error) {
	act := action.New(ctx, http.MethodPut, req)
	return act.Result()
}

func Delete(ctx context.Context, req g.Map) (res g.Map, err error) {
	act := action.New(ctx, http.MethodDelete, req)
	return act.Result()
}
