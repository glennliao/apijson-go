package config

import (
	"context"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

var (
	Debug = false
)

type AccessConditionFunc func(ctx context.Context, table string, req g.Map, reqRole string, needRole []string) (g.Map, error)

var (
	// AccessVerify 是否权限验证
	AccessVerify = false
	// AccessCondition 自定义权限限制条件
	AccessCondition AccessConditionFunc
)

// 自定义设置从ctx获取用户id和角色的key
var (
	RoleKey   = "ajg.role" // ctx 中role 的key
	UserIdKey = "ajg.userId"
)

// 设置 _access/_request 自定义表名
var (
	TableAccess  = "_access"
	TableRequest = "_request"
)

// =========================  角色 =======================
// 角色列表
// access 中填写的角色应在角色列表中

var (
	RoleList = []string{consts.UNKNOWN, consts.LOGIN, consts.OWNER, consts.ADMIN}
)

// AddRole 增加自定义角色
func AddRole(name string) {
	if !lo.Contains(RoleList, name) {
		RoleList = append(RoleList, name)
	}
}

// =========================  字段配置 =======================

// JsonFieldStyle DbFieldStyle 字段命名风格
var JsonFieldStyle = consts.CaseCamel
var DbFieldStyle = consts.CaseSnake

// todo 如果配置 DbFieldStyle 风格下划线， 则前端传进来的字段查询时候都转成下划线, 返回时则根据JsonFieldStyle转换, 如果JsonFieldStyle， DbFieldStyle 一致， 则可以看成不用转
// sqlexecute中处理数据库端的转换
// 返回的可在查询的field字段中使用 user_id as userId 完成转换
