package main

import (
	"context"
	"fmt"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"my-apijson/apijson"
	"my-apijson/apijson/db"
	"net/http"
	"time"
)

func main() {
	port := 8088
	s := g.Server()
	s.SetPort(port)
	s.SetDumpRouterMap(true)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			if r.Method == http.MethodOptions {
				r.Response.CORSDefault()
				return
			}
			r.Middleware.Next()
		})

		group.POST("/get", gfHandler("get"))
		group.POST("/post", gfHandler("post"))
		group.POST("/head", gfHandler("head"))
		group.POST("/put", gfHandler("put"))
		group.POST("/delete", gfHandler("delete"))

		// group.POST("/login", func() {
		//
		// })

	})

	db.Init()

	s.Run()
}

func gfHandler(p string) func(req *ghttp.Request) {

	var api func(ctx context.Context, req g.Map) (res g.Map, err error)

	switch p {
	case "get":
		api = apijson.Get
	case "post":
		api = apijson.Post
	case "head":
		api = apijson.Head
	case "put":
		api = apijson.Put
	case "delete":
		api = apijson.Delete
	}
	return func(req *ghttp.Request) {
		commonResponse(req, api)
	}
}

func commonResponse(req *ghttp.Request, handler func(ctx context.Context, req g.Map) (res g.Map, err error)) {
	res := g.Map{}

	req.GetMap()

	err := g.Try(req.Context(), func(ctx context.Context) {

		gmap.NewListMap()
		ret, err := handler(req.Context(), req.GetMap())

		if err == nil {
			res["code"] = 200
		} else {
			res["code"] = 500
			res["msg"] = err.Error()
			// g.Log().Stack(false)
			// g.Log().Error(req.Context(), err)
			// g.Log().Stack(true)
		}
		for k, v := range ret {
			res[k] = v
		}
	})
	if err != nil {
		res["code"] = 500
		res["msg"] = err.Error()
		if e, ok := err.(*gerror.Error); ok {
			g.Log().Stack(false).Error(req.Context(), err, e.Stack())
		} else {
			g.Log().Stack(false).Error(req.Context(), err)
		}
	}
	res["_span"] = fmt.Sprintf("%s", time.Since(time.UnixMilli(req.EnterTime)))
	req.Response.WriteJson(res)
}
