package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/iancoleman/orderedmap"
)

type GF struct {
	apijson          *apijson.ApiJson
	ResponseResolver func(handler func(ctx context.Context, req model.Map) (res model.Map, err error), mode Mode, debug bool) func(req *ghttp.Request)
}

func New(a *apijson.ApiJson) *GF {
	return &GF{
		apijson:          a,
		ResponseResolver: CommonResponse,
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
	group.POST("/get", gf.ResponseResolver(gf.Get, mode[0], gf.apijson.Debug))
	group.POST("/post", gf.ResponseResolver(gf.Post, mode[0], gf.apijson.Debug))
	group.POST("/head", gf.ResponseResolver(gf.Head, mode[0], gf.apijson.Debug))
	group.POST("/put", gf.ResponseResolver(gf.Put, mode[0], gf.apijson.Debug))
	group.POST("/delete", gf.ResponseResolver(gf.Delete, mode[0], gf.apijson.Debug))
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

// 调试模式开启, 使用orderedmap输出结果
func sortMap(ctx context.Context, body []byte, res *gmap.ListMap, ret model.Map) *orderedmap.OrderedMap {
	reqSortMap := orderedmap.New()
	err := reqSortMap.UnmarshalJSON(body)
	if err != nil {
		g.Log().Warning(ctx, err)
	}

	for _, k := range reqSortMap.Keys() {
		if strings.HasPrefix(k, consts.RefKeySuffix) {
			continue
		}
		if k == consts.Tag {
			continue
		}

		if strings.HasSuffix(k, consts.RefKeySuffix) {
			k = k[:len(k)-1]
		}

		res.Set(k, ret[k])
	}

	return reqSortMap
}

func try(ctx context.Context, try func(ctx context.Context) error) (err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.Newf(`%+v`, exception)
			}
		}
	}()
	err = try(ctx)
	return
}

type CodeErr interface {
	Code() int
	Error() string
}

func CommonResponse(handler func(ctx context.Context, req model.Map) (res model.Map, err error), mode Mode, debug bool) func(req *ghttp.Request) {
	return func(req *ghttp.Request) {
		metaRes := &gmap.ListMap{}
		code := 200
		msg := "success"
		nodeRes := &gmap.ListMap{}

		err := try(req.Context(), func(ctx context.Context) (err error) {

			ret, err := handler(ctx, req.GetMap())

			if debug {
				sortMap(ctx, req.GetBody(), metaRes, ret)
			} else {
				for k, v := range ret {
					nodeRes.Set(k, v)
				}
			}
			return
		})

		if err != nil {

			if e, ok := err.(CodeErr); ok {
				code = e.Code()
				if strconv.Itoa(e.Code())[0] == '4' {
					code = e.Code()
					msg = e.Error()
				} else {
					code = 500
					msg = "系统异常"
				}
			} else {

				if _, ok := err.(*gerror.Error); ok {
					// if e.Code() == gcode.CodeNil {
					// 	code = 400
					// 	msg = e.Error()
					// } else {
					// 	code = 500
					// 	msg = "系统异常"
					// }
					code = 500
					msg = "系统异常"
				} else {
					code = 500
					msg = "系统异常"
				}
			}

			if code >= 500 {
				g.Log().Stack(false).Errorf(req.Context(), "%+v", err)

				// if e, ok := err.(*gerror.Error); ok {
				// 	g.Log().Stack(false).Error(req.Context(), err, e.Stack())
				// } else {
				// 	g.Log().Stack(false).Error(req.Context(), err)
				// }
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
