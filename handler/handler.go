package handler

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/db"
	"github.com/glennliao/apijson-go/query"
	"github.com/gogf/gf/v2/frame/g"
)

func Get(ctx context.Context, req g.Map) (res g.Map, err error) {
	q := query.New(ctx, req)
	q.AccessVerify = config.AccessVerify
	q.AccessCondition = config.AccessCondition
	return q.Result()
}

func Post(ctx context.Context, req g.Map) (res g.Map, err error) {

	req, err = checkByRequest(ctx, req, "POST")
	if err != nil {
		return nil, err
	}

	var ret = g.Map{}

	for k, v := range req {
		if val, ok := v.(map[string]any); ok {

			id, count, err := db.Insert(ctx, k, val)
			if err != nil {
				ret[k] = g.Map{
					"code": 500,
					"msg":  err.Error(),
				}
			} else {
				ret[k] = g.Map{
					"code":  200,
					"id":    id,
					"count": count,
				}
			}
		}
	}
	return ret, err
}

func Head(ctx context.Context, req g.Map) (res g.Map, err error) {
	return nil, err
}

func Put(ctx context.Context, req g.Map) (res g.Map, err error) {

	req, err = checkByRequest(ctx, req, "PUT")
	if err != nil {
		return nil, err
	}

	ret := g.Map{}

	for k, v := range req {
		id, count, err := db.Update(ctx, k, v.(g.Map))

		if err != nil {
			ret[k] = g.Map{
				"code": 500,
				"msg":  err.Error(),
			}
		} else {
			ret[k] = g.Map{
				"code":  200,
				"id":    id,
				"count": count,
			}
		}

	}
	return ret, err
}

func Delete(ctx context.Context, req g.Map) (res g.Map, err error) {

	req, err = checkByRequest(ctx, req, "DELETE")
	if err != nil {
		return nil, err
	}

	ret := g.Map{}

	for k, v := range req {
		id, count, err := db.Delete(ctx, k, v.(g.Map))
		if err != nil {
			ret[k] = g.Map{
				"code": 500,
				"msg":  err.Error(),
			}
		} else {
			ret[k] = g.Map{
				"code":  200,
				"id":    id,
				"count": count,
			}
		}

	}
	return ret, err
}
