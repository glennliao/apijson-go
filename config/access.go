package config

import (
	"context"

	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"
)

type ConditionReq struct {
	AccessName          string // _access 中的alias
	TableAccessRoleList []string
	Method              string
	NodeReq             g.Map  // 节点的请求数据
	NodeRole            string // 节点的角色
}

type ConditionRet struct {
	condition    map[string]any
	rawCondition map[string][]any
}

func NewConditionRet() *ConditionRet {
	c := ConditionRet{
		condition:    map[string]any{},
		rawCondition: map[string][]any{},
	}
	return &c
}

func (c *ConditionRet) Add(k string, v any) {
	c.condition[k] = v
}

func (c *ConditionRet) AddRaw(k string, v ...any) {
	c.rawCondition[k] = v
}

func (c *ConditionRet) Where() map[string]any {
	if len(c.rawCondition) > 0 {
		c.condition[consts.Raw] = c.rawCondition
	}
	return c.condition
}

type AccessCondition func(ctx context.Context, req ConditionReq, condition *ConditionRet) error

type RoleReq struct {
	AccessName string
	Method     string
	NodeRole   string // 前端传入的节点的角色, 目前未传入则为空
}

type DefaultRole func(ctx context.Context, req RoleReq) (string, error)

func defaultRole(ctx context.Context, req RoleReq) (string, error) {
	return consts.UNKNOWN, nil
}

func defaultCondition(ctx context.Context, req ConditionReq, condition *ConditionRet) error {
	return nil
}

type Access struct {
	// 禁用_access权限校验, 默认为false, 需手动开启
	NoVerify bool

	// 用于 根据accessName+user来自定义添加sql条件, 完成数据的权限限制
	ConditionFunc AccessCondition

	// nodeRole 为前端显式指定的role, 需要此函数中判断该role是不是用户角色之一, 返回最终该节点的角色
	DefaultRoleFunc DefaultRole

	roleList []string

	accessConfigMap map[string]AccessConfig
}

func NewAccess() *Access {

	a := &Access{}
	a.ConditionFunc = defaultCondition
	a.DefaultRoleFunc = defaultRole
	a.roleList = []string{consts.UNKNOWN, consts.LOGIN, consts.OWNER, consts.ADMIN}

	return a
}

// AddRole 添加应用中额外的角色
func (a *Access) AddRole(roles []string) *Access {
	for _, role := range roles {
		if !lo.Contains(a.roleList, role) {
			a.roleList = append(a.roleList, role)
		}
	}
	return a
}

func (a *Access) RoleList() []string { return a.roleList }
