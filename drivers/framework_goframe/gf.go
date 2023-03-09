package framework_goframe

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/iancoleman/orderedmap"
	"net/http"
	"strings"
	"time"
)

type GF struct {
	apijson *apijson.ApiJson
}

func New(a *apijson.ApiJson) *GF {
	return &GF{
		apijson: a,
	}
}

func (gf *GF) Run(s ...*ghttp.Server) {

	var server *ghttp.Server

	if len(s) == 0 {
		server = g.Server("apijson")
	} else {
		server = s[0]
	}

	server.Group("/", func(group *ghttp.RouterGroup) {
		gf.Bind(group)
	})
	server.Run()
}

func (gf *GF) Bind(group *ghttp.RouterGroup, mode ...Mode) {
	if len(mode) == 0 {
		mode = []Mode{InDataMode}
	}
	group.POST("/get", gf.commonResponse(gf.Get, mode[0]))
	group.POST("/post", gf.commonResponse(gf.Post, mode[0]))
	group.POST("/head", gf.commonResponse(gf.Head, mode[0]))
	group.POST("/put", gf.commonResponse(gf.Put, mode[0]))
	group.POST("/delete", gf.commonResponse(gf.Delete, mode[0]))
}

func (gf *GF) Get(ctx context.Context, req model.Map) (res model.Map, err error) {
	q := gf.apijson.NewQuery(ctx, req)
	return q.Result()
}

func (gf *GF) Head(ctx context.Context, req model.Map) (res model.Map, err error) {
	return nil, err
}

func (gf *GF) Post(ctx context.Context, req model.Map) (res model.Map, err error) {
	act := gf.apijson.NewAction(ctx, http.MethodPost, req)
	return act.Result()
}

func (gf *GF) Put(ctx context.Context, req model.Map) (res model.Map, err error) {
	act := gf.apijson.NewAction(ctx, http.MethodPut, req)
	return act.Result()
}

func (gf *GF) Delete(ctx context.Context, req model.Map) (res model.Map, err error) {
	act := gf.apijson.NewAction(ctx, http.MethodDelete, req)
	return act.Result()
}

func (gf *GF) commonResponse(handler func(ctx context.Context, req model.Map) (res model.Map, err error), mode Mode) func(req *ghttp.Request) {
	return func(req *ghttp.Request) {
		metaRes := &gmap.ListMap{}
		code := 200
		msg := "success"
		nodeRes := &gmap.ListMap{}

		err := g.Try(req.Context(), func(ctx context.Context) {

			ret, err := handler(req.Context(), req.GetMap())

			if err == nil {
				code = 200
			} else {
				code = 500
				msg = err.Error()
			}

			if gf.apijson.Debug {
				sortMap(ctx, req.GetBody(), metaRes, ret)
			} else {
				for k, v := range ret {
					nodeRes.Set(k, v)
				}
			}

			if err != nil {
				panic(err)
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
