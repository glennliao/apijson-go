package action

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/db"
	"github.com/glennliao/apijson-go/functions"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"strings"
)

type Node struct {
	req       []g.Map
	ctx       context.Context
	Key       string
	TableName string
	Role      string

	Data   []g.Map // 需写入数据库的数据
	Where  []g.Map // 条件
	RowKey string  // 主键

	structure Structure

	keyNode map[string]*Node
}

func newNode(key string, req []g.Map, structure Structure) Node {
	return Node{
		Key: key, req: req, structure: structure,
	}
}

func (n *Node) parseReq(method string) {
	n.Data = []g.Map{}
	n.Where = []g.Map{}

	for i, item := range n.req {

		n.Data = append(n.Data, g.Map{})
		n.Where = append(n.Where, g.Map{})

		for key, val := range item {
			if key == consts.Role {
				n.Role = gconv.String(val)
			} else {
				key = config.GetDbFieldStyle()(n.ctx, n.TableName, key)

				if method == consts.MethodDelete {
					n.Where[i][key] = val
				} else {
					if key == n.RowKey || key == n.RowKey+"{}" {
						if method == consts.MethodPut {
							n.Where[i][key] = val
						}
						// Post 暂原则上不让传递这个rowKey值
					} else {
						n.Data[i][key] = val
					}
				}
			}
		}
	}

}

func (n *Node) parse(ctx context.Context, method string) error {

	key := n.Key
	if strings.HasSuffix(key, consts.ListKeySuffix) {
		key = key[0 : len(key)-2] // todo 提取util, 获取非数组的key
	}
	access, err := db.GetAccess(key, true)

	if err != nil {
		return err
	}

	n.TableName = access.Name
	n.RowKey = access.RowKey

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
		if n.Role == "" {
			n.Role = gconv.String(val)
		}
	}

	if val, exists := n.structure.Update[consts.Role]; exists {
		n.Role = gconv.String(val)
	}

	return nil
}

func (n *Node) checkAccess(ctx context.Context, method string, accessRoles []string) error {

	role, err := config.DefaultRoleFunc(ctx, config.RoleReq{
		Table:    n.TableName,
		Method:   method,
		NodeRole: n.Role,
	})

	if err != nil {
		return err
	}

	if role == consts.DENY {
		return gerror.Newf("deny node: %s with %s", n.Key, n.Role)
	}

	n.Role = role

	if !lo.Contains(accessRoles, role) {
		return gerror.Newf("node not access: %s with %s", n.Key, n.Role)
	}

	for i, item := range n.req {
		where, err := config.AccessConditionFunc(ctx, config.AccessConditionReq{
			Table:               n.TableName,
			TableAccessRoleList: accessRoles,
			Method:              method,
			NodeRole:            n.Role,
			NodeReq:             item,
		})

		if err != nil {
			return err
		}

		if method == consts.MethodPost {
			for k, v := range where {
				n.Data[i][k] = v
			}
		} else {
			for k, v := range where {
				n.Where[i][k] = v
			}
		}
	}

	return nil
}

func (n *Node) checkReq() error {

	for _, item := range n.req {
		// must
		for _, key := range n.structure.Must {
			if _, exists := item[key]; !exists {
				return gerror.New("structure错误: 400, 缺少" + n.Key + "." + key)
			}
		}

		// refuse
		if n.structure.Refuse[0] == "!" {
			if len(n.structure.Must) == 0 {
				return gerror.New("structure错误: 400, REFUSE为!时必须指定MUST" + n.Key)
			}

			for key, _ := range item {
				if !lo.Contains(n.structure.Must, key) {
					return gerror.New("structure错误: 400, 不能包含" + n.Key + "." + key)
				}
			}

		} else {
			for _, key := range n.structure.Refuse {
				if _, exists := item[key]; exists {
					return gerror.New("structure错误: 400, 不能包含" + n.Key + "." + key)
				}
			}
		}
	}

	return nil
}

func (n *Node) reqUpdate() error {

	for i, _ := range n.req {
		for key, updateVal := range n.structure.Update {

			if strings.HasSuffix(key, consts.FunctionsKeySuffix) {
				functionName, paramKeys := functions.ParseFunctionsStr(updateVal.(string))
				var param = g.Map{}
				for _, key := range paramKeys {
					if key == "$req" {
						param[key] = n.Data[i]
					} else {
						param[key] = n.Data[i][key]
					}
				}
				k := key[0 : len(key)-2]
				val, err := functions.Call(n.ctx, functionName, param)
				if err != nil {
					return err
				}
				if val != nil {
					n.Data[i][k] = val
				}
			} else {
				n.Data[i][key] = updateVal
			}
		}

		for key, updateVal := range n.structure.Insert {
			if _, exists := n.Data[i][key]; !exists {
				n.Data[i][key] = updateVal
			}
		}

		for k, v := range n.Data[i] {
			if strings.HasSuffix(k, "@") {
				refNodeKey, refCol := parseRefCol(v.(string))
				if strings.HasSuffix(refNodeKey, "[]") { // 双列表
					n.Data[i][k] = n.keyNode[refNodeKey].Data[i][config.GetDbFieldStyle()(n.ctx, n.TableName, refCol)]
				} else {
					n.Data[i][k] = n.keyNode[refNodeKey].Data[0][config.GetDbFieldStyle()(n.ctx, n.TableName, refCol)]
				}
			}
		}
	}

	return nil
}

func (n *Node) do(ctx context.Context, method string, i int) (ret g.Map, err error) {

	// todo 此处运行会导致事务时长与hook时长相关,特别是hook中运行了io类型的操作, 故需要调整到事务外去执行, 且如果事务失败, 则不执行after, 可以改成增加error
	for _, hook := range hooks {
		if hook.Before != nil {
			err := hook.Before(n, method)
			if err != nil {
				return nil, err
			}
		}
	}

	switch method {
	case consts.MethodPost:

		var rowKeyVal g.Map
		for i, _ := range n.Data {
			rowKeyVal, err = rowKeyGen(ctx, n.TableName, n.Data[i])
			if err != nil {
				return nil, err
			}

			if rowKeyVal != nil {
				for k, v := range rowKeyVal {
					n.Data[i][k] = v
				}
			}

		}

		id, count, err := db.Insert(ctx, n.TableName, n.Data)
		if err != nil {
			return nil, err
		}

		ret = g.Map{
			"code":  200,
			"count": count,
			"id":    id,
		}

		if len(n.Data) > 0 { //多条插入时返回值已经应该无意义了

			jsonStyle := config.GetJsonFieldStyle()
			if rowKeyVal != nil {
				for k, v := range rowKeyVal {
					ret[jsonStyle(ctx, n.TableName, k)] = v
				}
			}
		}

	case consts.MethodPut:
		count, err := db.Update(ctx, n.TableName, n.Data[i], n.Where[i])
		if err != nil {
			return nil, err
		}

		ret = g.Map{
			"code":  200,
			"count": count,
		}
	case consts.MethodDelete:
		count, err := db.Delete(ctx, n.TableName, n.Where[i])
		if err != nil {
			return nil, err
		}

		ret = g.Map{
			"code":  200,
			"count": count,
		}
	}

	if ret == nil {
		return nil, gerror.New("undefined method:" + method)
	}

	for _, hook := range hooks {
		if hook.After != nil {
			err := hook.After(n, method)
			if err != nil {
				return nil, err
			}
		}
	}

	return
}

func (n *Node) execute(ctx context.Context, method string) (g.Map, error) {

	// 参数替换
	err := n.reqUpdate() // todo 处理放到事务外, 减短事务时长
	if err != nil {
		return nil, err
	}

	// 执行操作

	if method == consts.MethodPost { // 新增时可以合并新增
		ret, err := n.do(ctx, method, 0)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		for i, _ := range n.req {
			_, err = n.do(ctx, method, i)
			if err != nil {
				return nil, err
			}
		}

	}

	return g.Map{
		"code": 200,
	}, nil
}
