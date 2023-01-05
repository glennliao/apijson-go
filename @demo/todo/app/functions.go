package app

import (
	"context"
	"errors"
	"github.com/glennliao/apijson-go/config/functions"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

const (
	UserIdWM = "10001"
	UserIdSQ = "10002"
)

func init() {

	functions.Reg("sayHello", functions.Func{
		Handler: func(ctx context.Context, param g.Map) (res any, err error) {
			return "world", err
			//return nil, nil
		},
	})

	functions.Reg("sayHi", functions.Func{
		Handler: func(ctx context.Context, param g.Map) (res any, err error) {
			return "你好:" + gconv.String(param["realname"]), err
		},
	})

	functions.Reg("checkTodoTitle", functions.Func{
		Handler: func(ctx context.Context, param g.Map) (res any, err error) {
			user, _ := ctx.Value(UserIdKey).(*CurrentUser)

			if user.UserId == UserIdSQ && strings.HasSuffix(gconv.String(param["title"]), "喝茶") {
				return nil, errors.New("操作不允许")
			}

			return nil, nil
		},
	})

	functions.Reg("updateTodoTitle", functions.Func{
		Handler: func(ctx context.Context, param g.Map) (res any, err error) {
			user, _ := ctx.Value(UserIdKey).(*CurrentUser)

			if user.UserId == UserIdSQ && strings.HasSuffix(gconv.String(param["title"]), "找林云逛街") {
				return strings.Replace(gconv.String(param["title"]), "找", "保护", -1), nil
			}

			return gconv.String(param["title"]), nil
		},
	})

}
