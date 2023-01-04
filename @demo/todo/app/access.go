package app

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

// 自定义设置从ctx获取用户id和角色的key

const UserIdKey = "ajg.userId"

type CurrentUser struct {
	UserId string
}

const PARTNER = "PARTNER" // 伙伴角色, 可查看指定与自己有关的todo

// 以下角色、权限判断有些复杂，if判断太多， 可能需要考虑弄个规则引擎将配置化数据自动转成条件判断

func Role(ctx context.Context, req config.RoleReq) (string, error) {
	_, ok := ctx.Value(UserIdKey).(*CurrentUser)

	if !ok {
		return consts.UNKNOWN, nil //未登录
	}

	if req.NodeRole == "" {

		switch req.Table {
		case "t_todo", "t_user":
			return consts.OWNER, nil
		}

	} else {

		switch req.Table {
		case "t_todo", "t_user":

			if req.NodeRole == consts.LOGIN {
				req.NodeRole = consts.OWNER
			}

			if lo.Contains([]string{consts.OWNER, PARTNER}, req.NodeRole) {
				return req.NodeRole, nil
			}

			return consts.DENY, nil // 非拥有的角色

		default:
			return req.NodeRole, nil
		}
	}

	return consts.LOGIN, nil

}

// 用户访问的角色为单次单个,  请求时候指定用户角色, 如果没有指定则默认OWNER (获取request中指定, 或者自定义 (不同app不同角色))
// 此处的角色为系统用户角色, 即为未登录用户、普通用户、机构、 后台管理员、 （业务角色 （例如todo的伙伴））, 不是后台管理员总的角色,
// 后台管理员中的角色 需要另外处理, 针对 ADMIN 角色, 通过读取系统配置表判断该用户是否对该数据表具有get,post,put,delete权限, 然后需要自定义实现他们如何做行控制条件, 以及字段控制
// 后台导入、导出如何搞呢 -> 统一导入导出模块, 然后调用 apijson 模板完成数据查找、处理、然后再统一导入导出, 还可以注册自定义导出handler, 处理复杂导入导出

func AccessCondition(ctx context.Context, req config.AccessConditionReq) (g.Map, error) {

	user, ok := ctx.Value(UserIdKey).(*CurrentUser)

	if !ok {
		return nil, nil
	}

	switch req.Table {
	case "t_user":
		if req.NodeRole == consts.OWNER {
			return g.Map{
				"user_id": user.UserId,

				"@raw": g.Map{
					"id > 0": "",
				},
				// g.Map{
				//    "uid <=" : 1000,
				//    "age >=" : 18,
				//    "x in (?)":[]string{"1","2","3"}
				//}
			}, nil
		}
	case "t_todo":
		if req.NodeRole == consts.OWNER {
			return g.Map{"user_id": user.UserId}, nil
		}
		if req.NodeRole == PARTNER {
			return g.Map{"partner": user.UserId}, nil
		}
	}

	return nil, nil
}
