package config

import (
	"context"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
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

func AddRole(name string) {
	if !lo.Contains(roleList, name) {
		roleList = append(roleList, name)
	}
}

func RoleList() []string { return roleList }
