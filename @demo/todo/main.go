package main

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/framework"
	"github.com/glennliao/apijson-go/framework/handler"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"todo/app"
)

func main() {

	framework.Init()

	config.AccessVerify = false // 全局配置验证权限开关
	config.AccessConditionFunc = app.AccessCondition
	config.DefaultRoleFunc = app.Role
	config.Debug = true

	s := g.Server()

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(cors, auth)
		handler.Bind(group)
	})

	s.Run()
}

func auth(r *ghttp.Request) {
	// 模拟认证, 获取用户信息, 此处为简化直接Authorization传递userId, 实际项目中请勿这么操作
	authorization := r.Request.Header.Get("Authorization")

	if authorization != "" {
		ctx := r.Context()
		ctx = context.WithValue(ctx, app.UserIdKey, &app.CurrentUser{UserId: authorization})
		r.SetCtx(ctx)
	} else {
		if r.URL.Path != "/get" { // 此处限制非查询的都需要登录, 可结合实际调整
			r.Response.WriteJson(g.Map{
				"code": 401,
				"msg":  "未登录",
			})
			return
		}
	}

	r.Middleware.Next()
}

func cors(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowOrigin = r.Request.Header.Get("Origin")
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}
