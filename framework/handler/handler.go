package handler

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/query"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/iancoleman/orderedmap"
	"net/http"
	"strings"
	"time"
)

func Get(ctx context.Context, req model.Map) (res model.Map, err error) {
	q := query.New(ctx, req)
	q.NoAccessVerify = config.NoAccessVerify
	q.AccessCondition = config.AccessConditionFunc
	return q.Result()
}

func Head(ctx context.Context, req model.Map) (res model.Map, err error) {
	return nil, err
}

func Post(ctx context.Context, req model.Map) (res model.Map, err error) {
	act := action.New(ctx, http.MethodPost, req)
	act.NoAccessVerify = config.NoAccessVerify
	return act.Result()
}

func Put(ctx context.Context, req model.Map) (res model.Map, err error) {
	act := action.New(ctx, http.MethodPut, req)
	act.NoAccessVerify = config.NoAccessVerify
	return act.Result()
}

func Delete(ctx context.Context, req model.Map) (res model.Map, err error) {
	act := action.New(ctx, http.MethodDelete, req)
	act.NoAccessVerify = config.NoAccessVerify
	return act.Result()
}

type Mode = func(data gmap.ListMap, meta gmap.ListMap) gmap.ListMap

func SpreadMode(data gmap.ListMap, meta gmap.ListMap) gmap.ListMap {

	res := gmap.ListMap{}
	for _, k := range data.Keys() {
		res.Set(k, data.Get(k))
	}
	for _, k := range meta.Keys() {
		res.Set(k, meta.Get(k))
	}

	return res
}

func InDataMode(data gmap.ListMap, meta gmap.ListMap) gmap.ListMap {
	res := gmap.ListMap{}
	res.Set("data", data)
	for _, k := range meta.Keys() {
		res.Set(k, meta.Get(k))
	}
	return res
}

func Bind(group *ghttp.RouterGroup, mode ...Mode) {
	if len(mode) == 0 {
		mode = []Mode{InDataMode}
	}
	group.POST("/get", CommonResponse(Get, mode[0]))
	group.POST("/post", CommonResponse(Post, mode[0]))
	group.POST("/head", CommonResponse(Head, mode[0]))
	group.POST("/put", CommonResponse(Put, mode[0]))
	group.POST("/delete", CommonResponse(Delete, mode[0]))
}

func CommonResponse(handler func(ctx context.Context, req model.Map) (res model.Map, err error), mode Mode) func(req *ghttp.Request) {
	return func(req *ghttp.Request) {
		metaRes := gmap.ListMap{}
		code := 200
		msg := "success"
		nodeRes := gmap.ListMap{}

		err := g.Try(req.Context(), func(ctx context.Context) {

			ret, err := handler(req.Context(), req.GetMap())

			if err == nil {
				code = 200
			} else {
				code = 500
				msg = err.Error()
			}

			if config.Debug {
				sortMap(ctx, req.GetBody(), &metaRes, ret)
			} else {
				for k, v := range ret {
					nodeRes.Set(k, v)
				}
			}

		})

		if err != nil {
			code = 500
			msg = err.Error()
			if e, ok := err.(*gerror.Error); ok {
				g.Log().Stack(false).Error(req.Context(), err, e.Stack())
			} else {
				g.Log().Stack(false).Error(req.Context(), err)
			}
		}

		metaRes.Set("ok", code == 200)
		metaRes.Set("code", code)
		metaRes.Set("msg", msg)
		metaRes.Set("span", fmt.Sprintf("%s", time.Since(time.UnixMilli(req.EnterTime))))

		res := mode(nodeRes, metaRes)
		req.Response.WriteJson(res.String())
	}
}

// 调试模式开启, 使用orderedmap输出结果
func sortMap(ctx context.Context, body []byte, res *gmap.ListMap, ret model.Map) *orderedmap.OrderedMap {
	reqSortMap := orderedmap.New()
	err := reqSortMap.UnmarshalJSON(body)
	if err != nil {
		g.Log().Warning(ctx, err)
	}

	for _, k := range reqSortMap.Keys() {
		if strings.HasPrefix(k, "@") {
			continue
		}
		if k == "tag" {
			continue
		}

		if strings.HasSuffix(k, "@") {
			k = k[:len(k)-1]
		}

		res.Set(k, ret[k])
	}

	return reqSortMap
}
