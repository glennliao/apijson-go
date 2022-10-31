package config

import (
	"context"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/samber/lo"
	"strings"
)

var (
	Debug = false
)

type AccessConditionReq struct {
	Table               string
	TableAccessRoleList []string
	Method              string
	NodeReq             g.Map  //节点的请求数据
	NodeRole            string // 节点的角色
}

type RoleReq struct {
	Table    string
	Method   string
	NodeRole string // 前端传入的节点的角色, 目前未传入则为空
}

// AccessCondition 根据传入的ctx获取用户信息, 结合req 中的信息 返回需要添加到sql的where条件
type AccessCondition func(ctx context.Context, req AccessConditionReq) (g.Map, error)

// DefaultRole nodeRole 为前端显式指定的role, 需要此函数中判断该role是不是用户角色之一, 返回最终该节点的角色
type DefaultRole func(ctx context.Context, req RoleReq) (string, error)

var (
	// AccessVerify 是否权限验证
	AccessVerify = false
	// AccessConditionFunc 自定义权限限制条件
	AccessConditionFunc AccessCondition
	// DefaultRoleFunc 自定义获取节点的默认角色
	DefaultRoleFunc DefaultRole = func(ctx context.Context, req RoleReq) (string, error) {
		return consts.UNKNOWN, nil
	}
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
	roleList = []string{consts.UNKNOWN, consts.LOGIN, consts.OWNER, consts.ADMIN}
)

// AddRole 增加自定义角色
func AddRole(name string) {
	if !lo.Contains(roleList, name) {
		roleList = append(roleList, name)
	}
}

func RoleList() []string { return roleList }

// =========================  请求方法 =======================

var methodList = []string{consts.MethodGet, consts.MethodHead, consts.MethodPost, consts.MethodPut, consts.MethodDelete}

func MethodList() []string { return methodList }

// =========================  字段配置 =======================

// jsonFieldStyle dbFieldStyle 配置Json字段, 数据库命名风格
var jsonFieldStyle = consts.CaseCamel
var dbFieldStyle = consts.CaseSnake

// TODO 从配置文件读取命名风格

// SetJsonFieldStyle 设置返回的 json字段风格, 默认为 lowerCamelCase风格, 参考 gstr.CaseCamelLower
func SetJsonFieldStyle(style string) {
	jsonFieldStyle = fieldStyle(style)
}

// SetDbFieldStyle 设置数据库的字段风格, 默认为 snake_case风格, 参考 gstr.CaseSnake
func SetDbFieldStyle(style string) {
	dbFieldStyle = fieldStyle(style)
}

func fieldStyle(style string) int {

	switch strings.ToUpper(style) {
	case consts.CASE_CAMEL:
		return consts.CaseCamel
	case consts.CASE_CAMEL_UPPER:
		return consts.CaseCamelUpper
	case consts.CASE_SNAKE:
		return consts.CaseSnake
	}

	return consts.Origin
}

func convFieldStyle(style int, field string) string {
	switch style {

	case consts.CaseCamel:
		return gstr.CaseCamelLower(field)
	case consts.CaseCamelUpper:
		return gstr.CaseCamel(field)
	case consts.CaseSnake:
		return gstr.CaseSnake(field)
	default:
		return field
	}
}

func ToDbField(field string) string {
	return convFieldStyle(dbFieldStyle, field)
}

func ToJsonField(field string) string {
	return convFieldStyle(jsonFieldStyle, field)
}
