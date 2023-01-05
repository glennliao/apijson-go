package handler

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
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

func Bind(group *ghttp.RouterGroup) {
	group.POST("/get", CommonResponse(Get))
	group.POST("/post", CommonResponse(Post))
	group.POST("/head", CommonResponse(Head))
	group.POST("/put", CommonResponse(Put))
	group.POST("/delete", CommonResponse(Delete))
}

func CommonResponse(handler func(ctx context.Context, req g.Map) (res g.Map, err error)) func(req *ghttp.Request) {
	return func(req *ghttp.Request) {
		res := gmap.ListMap{}
		code := 200
		msg := "success"

		err := g.Try(req.Context(), func(ctx context.Context) {

			ret, err := handler(req.Context(), req.GetMap())

			if err == nil {
				code = 200
			} else {
				code = 500
				msg = err.Error()
			}

			if config.Debug {
				sortMap(ctx, req.GetBody(), &res, ret)
			} else {
				for k, v := range ret {
					res.Set(k, v)
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

		res.Set("ok", code == 200)
		res.Set("code", code)
		res.Set("msg", msg)
		res.Set("span", fmt.Sprintf("%s", time.Since(time.UnixMilli(req.EnterTime))))
		req.Response.WriteJson(res.String())
	}
}

// 调试模式开启, 使用orderedmap输出结果
func sortMap(ctx context.Context, body []byte, res *gmap.ListMap, ret g.Map) *orderedmap.OrderedMap {
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
