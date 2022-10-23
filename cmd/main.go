package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/db"
	"github.com/glennliao/apijson-go/handler"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/iancoleman/orderedmap"
	"github.com/samber/lo"
	"strings"
	"time"
)

func main() {

	db.Init()

	s := g.Server()

	s.BindMiddleware("/*", func(r *ghttp.Request) {
		corsOptions := r.Response.DefaultCORSOptions()
		corsOptions.AllowOrigin = r.Request.Header.Get("Origin")
		r.Response.CORS(corsOptions)
		r.Middleware.Next()
	})

	s.Group("/", func(group *ghttp.RouterGroup) {

		group.Middleware(func(r *ghttp.Request) {
			// 模拟认证, 获取用户角色、用户信息
			ctx := r.Context()
			authorization := r.Request.Header.Get("Authorization")
			if authorization != "" {
				ctx = context.WithValue(ctx, config.RoleKey, []string{consts.OWNER, "PARTNER", consts.LOGIN})
				ctx = context.WithValue(ctx, config.UserIdKey, authorization)
			} else {
				ctx = context.WithValue(ctx, config.RoleKey, []string{consts.UNKNOWN})
			}
			r.SetCtx(ctx)
			r.Middleware.Next()
		})

		group.POST("/get", gfHandler("get"))
		group.POST("/post", gfHandler("post"))
		group.POST("/head", gfHandler("head"))
		group.POST("/put", gfHandler("put"))
		group.POST("/delete", gfHandler("delete"))
	})

	config.AccessVerify = false
	config.AccessCondition = accessCondition
	config.Debug = true

	s.Run()
}

func accessCondition(ctx context.Context, table string, req g.Map, reqRole string, needRole []string) (g.Map, error) {

	userRole := ctx.Value(config.RoleKey).([]string)

	// 可改成switch方式

	//if lo.Contains(needRole, consts.UNKNOWN) {
	//	return nil, nil
	//}
	//
	//if lo.Contains(needRole, consts.LOGIN) && lo.Contains(userRole, consts.LOGIN) { // 登录后公开资源
	//	return nil, nil
	//}

	// 用户访问的角色为单次单个,  请求时候指定用户角色, 如果没有指定则默认OWNER (获取request中指定, 或者自定义 (不同app不同角色))
	// 此处的角色为系统用户角色, 即为未登录用户、普通用户、机构、 后台管理员、 （业务角色 （例如todo的伙伴））, 不是后台管理员总的角色,
	// 后台管理员中的角色 需要另外处理, 针对 ADMIN 角色, 通过读取系统配置表判断该用户是否对该数据表具有get,post,put,delete权限, 然后需要自定义实现他们如何做行控制条件, 以及字段控制
	// 后台导入、导出如何搞呢 -> 统一导入导出模块, 然后调用 apijson 模板完成数据查找、处理、然后再统一导入导出, 还可以注册自定义导出handler, 处理复杂导入导出

	if table == "t_todo" {
		if req["@role"] == "PARTNER" && lo.Contains(needRole, "PARTNER") {
			return g.Map{
				"partner": ctx.Value("ajg.userId"),
			}, nil
		}

		if req["@role"] == "ADMIN" && lo.Contains(needRole, "ADMIN") {

			if ctx.Value("ajg.userId").(string) == "10001" { // 模拟admin
				return g.Map{
					"partner": ctx.Value("ajg.userId"),
				}, nil
			} else {
				req["@role"] = "OWNER"
			}
		}

		if req["@role"] == "OWNER" && lo.Contains(needRole, "OWNER") {
			return g.Map{
				"user_id": ctx.Value("ajg.userId"),
			}, nil
		}
	}

	if lo.Contains(needRole, consts.OWNER) && lo.Contains(userRole, consts.OWNER) {
		if table == "User" {
			return g.Map{
				"id": ctx.Value("ajg.userId"),
			}, nil
		} else {
			return g.Map{
				"user_id": ctx.Value("ajg.userId"),
			}, nil
		}
	}

	return nil, nil
}

func gfHandler(p string) func(req *ghttp.Request) {

	var api func(ctx context.Context, req g.Map) (res g.Map, err error)

	switch p {
	case "get":
		api = handler.Get
	case "post":
		api = handler.Post
	case "head":
		api = handler.Head
	case "put":
		api = handler.Put
	case "delete":
		api = handler.Delete
	}
	return func(req *ghttp.Request) {
		commonResponse(req, api)
	}
}

func commonResponse(req *ghttp.Request, handler func(ctx context.Context, req g.Map) (res g.Map, err error)) {
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

		if config.Debug { //调试模式开启
			reqSortMap := orderedmap.New()

			err := json.Unmarshal(req.GetBody(), reqSortMap)
			if err != nil {
				g.Log().Error(req.Context(), err)
			}
			for _, k := range reqSortMap.Keys() {
				if strings.HasPrefix(k, "@") {
					continue
				}

				res.Set(k, ret[k])
			}

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
	res.Set("code", code)
	res.Set("msg", msg)
	res.Set("_span", fmt.Sprintf("%s", time.Since(time.UnixMilli(req.EnterTime))))
	req.Response.WriteJson(res.String())
}
