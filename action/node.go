package action

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/db"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
)

type Node struct {
	req       g.Map
	key       string
	tableName string
	role      string

	data   g.Map  // 需写入数据库的数据
	where  g.Map  // 条件
	rowKey string // 主键

	structure Structure
}

func newNode(key string, req g.Map, structure Structure) Node {
	return Node{
		key: key, req: req, structure: structure,
	}
}

func (n *Node) parseReq(method string) {
	n.data = g.Map{}
	n.where = g.Map{}
	for key, val := range n.req {
		if key == consts.Role {
			n.role = gconv.String(val)
		} else {
			if method == consts.MethodDelete {
				n.where[key] = val
			} else {
				if key == n.rowKey {
					if method == consts.MethodPut {
						n.where[key] = val
					}
					// Post 暂原则上不让传递这个值
				} else {
					n.data[key] = val
				}
			}
		}
	}
}

func (n *Node) parse(ctx context.Context, method string) error {

	access, err := db.GetAccess(n.key, true)

	if err != nil {
		return err
	}

	n.tableName = access.Name
	n.rowKey = access.RowKey

	n.parseReq(method)

	// 0. 角色替换

	err = n.roleUpdate()
	if err != nil {
		return err
	}

	// 1. 检查权限, 无权限就不用做参数检查了
	var accessRoles []string

	switch method {
	case consts.MethodPost:
		accessRoles = access.Post
	case consts.MethodPut:
		accessRoles = access.Put
	case consts.MethodDelete:
		accessRoles = access.Delete
	}

	err = n.checkAccess(ctx, method, accessRoles)
	if err != nil {
		return err
	}

	// 2. 检查参数
	err = n.checkReq()
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) roleUpdate() error {

	if val, exists := n.structure.Insert[consts.Role]; exists {
		if n.role == "" {
			n.role = gconv.String(val)
		}
	}

	if val, exists := n.structure.Update[consts.Role]; exists {
		n.role = gconv.String(val)
	}

	return nil
}

func (n *Node) checkAccess(ctx context.Context, method string, accessRoles []string) error {

	role, err := config.DefaultRoleFunc(ctx, config.RoleReq{
		Table:    n.tableName,
		Method:   method,
		NodeRole: n.role,
	})

	if err != nil {
		return err
	}

	if role == consts.DENY {
		return gerror.Newf("deny node: %s with %s", n.key, n.role)
	}

	n.role = role

	if !lo.Contains(accessRoles, role) {
		return gerror.Newf("node not access: %s with %s", n.key, n.role)
	}

	where, err := config.AccessConditionFunc(ctx, config.AccessConditionReq{
		Table:               n.tableName,
		TableAccessRoleList: accessRoles,
		Method:              method,
		NodeRole:            n.role,
		NodeReq:             n.req,
	})

	if err != nil {
		return err
	}

	if method == consts.MethodPost {
		for k, v := range where {
			n.data[k] = v
		}
	} else {
		for k, v := range where {
			n.where[k] = v
		}
	}

	return nil
}

func (n *Node) checkReq() error {

	// must
	for _, key := range n.structure.Must {
		if _, exists := n.req[key]; !exists {
			return gerror.New("structure错误: 400, 缺少" + n.key + "." + key)
		}
	}

	// refuse
	if n.structure.Refuse[0] == "!" {
		if len(n.structure.Must) == 0 {
			return gerror.New("structure错误: 400, REFUSE为!时必须指定MUST" + n.key)
		}

		for key, _ := range n.req {
			if !lo.Contains(n.structure.Must, key) {
				return gerror.New("structure错误: 400, 不能包含" + n.key + "." + key)
			}
		}

	} else {
		for _, key := range n.structure.Refuse {
			if _, exists := n.req[key]; exists {
				return gerror.New("structure错误: 400, 不能包含" + n.key + "." + key)
			}
		}
	}

	return nil
}

func (n *Node) reqUpdate() error {

	for key, updateVal := range n.structure.Update {
		n.data[key] = updateVal
	}

	for key, updateVal := range n.structure.Insert {
		if _, exists := n.data[key]; !exists {
			n.data[key] = updateVal
		}
	}

	return nil
}

func (n *Node) do(ctx context.Context, method string) (g.Map, error) {

	switch method {
	case consts.MethodPost:
		id, count, err := db.Insert(ctx, n.tableName, n.data)
		if err != nil {
			return nil, err
		}

		return g.Map{
			"code":  200,
			"id":    id,
			"count": count,
		}, nil
	case consts.MethodPut:
		count, err := db.Update(ctx, n.tableName, n.data, n.where)
		if err != nil {
			return nil, err
		}

		return g.Map{
			"code":  200,
			"count": count,
		}, nil
	case consts.MethodDelete:
		count, err := db.Delete(ctx, n.tableName, n.where)
		if err != nil {
			return nil, err
		}

		return g.Map{
			"code":  200,
			"count": count,
		}, nil
	}

	return nil, gerror.New("undefined method:" + method)
}

func (n *Node) execute(ctx context.Context, method string) (g.Map, error) {

	// 参数替换
	err := n.reqUpdate()
	if err != nil {
		return nil, err
	}

	// 执行操作
	return n.do(ctx, method)
}
