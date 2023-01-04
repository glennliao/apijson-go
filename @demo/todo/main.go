package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/db"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/iancoleman/orderedmap"
	"strings"
	"time"
	"todo/app"
)

func main() {

	db.Init()

	config.AccessVerify = false // 全局配置验证权限开关
	config.AccessConditionFunc = app.AccessCondition
	config.DefaultRoleFunc = app.Role
	config.Debug = true

	s := g.Server()

	s.BindMiddleware("/*", func(r *ghttp.Request) {
		corsOptions := r.Response.DefaultCORSOptions()
		corsOptions.AllowOrigin = r.Request.Header.Get("Origin")
		r.Response.CORS(corsOptions)
		r.Middleware.Next()
	})

	s.Group("/", func(group *ghttp.RouterGroup) {

		group.Middleware(auth)

		group.POST("/get", commonResponse(app.Get))
		group.POST("/post", commonResponse(app.Post))
		group.POST("/head", commonResponse(app.Head))
		group.POST("/put", commonResponse(app.Put))
		group.POST("/delete", commonResponse(app.Delete))
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

// commonResponse 创建统一响应
func commonResponse(handler func(ctx context.Context, req g.Map) (res g.Map, err error)) func(req *ghttp.Request) {

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

			if config.Debug { //调试模式开启, 使用orderedmap输出结果
				reqSortMap := orderedmap.New()

				err = json.Unmarshal(req.GetBody(), reqSortMap)
				if err != nil {
					g.Log().Error(req.Context(), err)
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
