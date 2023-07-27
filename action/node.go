package action

import (
	"context"
	"net/http"
	"strings"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/samber/lo"
)

type Node struct {
	req       []model.Map
	ctx       context.Context
	Action    *Action
	Key       string
	IsList    bool
	tableName string
	Role      string

	Data   []model.Map // 需写入数据库的数据
	Where  []model.Map // 条件
	Ret    model.Map   // 节点返回值
	RowKey string      // 主键

	structure *config.Structure
	executor  string

	keyNode map[string]*Node

	// access *config.Access
}

func newNode(key string, req []model.Map, structure *config.Structure, executor string) Node {

	n := Node{
		Key: key, req: req, structure: structure, executor: executor,
	}

	n.Data = []model.Map{}
	n.Where = []model.Map{}

	for _ = range n.req {

		n.Data = append(n.Data, model.Map{})
		n.Where = append(n.Where, model.Map{})

	}

	if strings.HasSuffix(key, consts.ListKeySuffix) {
		n.Key = util.RemoveSuffix(key, consts.ListKeySuffix)
		n.IsList = true
	}

	return n
}

// parse req data to data/where
func (n *Node) parseReq(method string) {

	for i, item := range n.req {
		for key, val := range item {

			if key == consts.Role {
				n.Role = util.String(val)
				continue
			}

			key = n.Action.DbFieldStyle(n.ctx, n.tableName, key)

			switch method {
			case http.MethodPost:
				n.Data[i][key] = val
			case http.MethodDelete:
				n.Where[i][key] = val
			case http.MethodPut:
				if key == n.RowKey || key == n.RowKey+consts.OpIn {
					n.Where[i][key] = val
				} else {
					n.Data[i][key] = val
				}
			}

		}
	}
}

// parse node
func (n *Node) parse(ctx context.Context, method string) error {

	key := n.Key
	if strings.HasSuffix(key, consts.ListKeySuffix) {
		key = util.RemoveSuffix(key, consts.ListKeySuffix)
	}
	access, err := n.Action.ActionConfig.GetAccessConfig(key, true)

	if err != nil {
		return err
	}

	n.tableName = access.Name
	n.RowKey = access.RowKey

	// 0. 角色替换

	err = n.roleUpdate()
	if err != nil {
		return err
	}
	var accessRoles []string
	if n.Action.NoAccessVerify == false {
		// 1. 检查权限, 无权限就不用做参数检查了

		switch method {
		case http.MethodPost:
			accessRoles = access.Post
		case http.MethodPut:
			accessRoles = access.Put
		case http.MethodDelete:
			accessRoles = access.Delete
		}

		err = n.checkAccess(ctx, method, accessRoles)
		if err != nil {
			return err
		}
	}

	// 2. 检查参数
	err = n.checkReq()
	if err != nil {
		return err
	}

	// 3. get where by accessCondition
	err = n.whereUpdate(ctx, method, accessRoles)

	n.parseReq(method)

	return err
}

// update node role
func (n *Node) roleUpdate() error {

	if val, exists := n.structure.Insert[consts.Role]; exists {
		if n.Role == "" {
			n.Role = util.String(val)
		}
	}

	if val, exists := n.structure.Update[consts.Role]; exists {
		n.Role = util.String(val)
	}

	return nil
}

func (n *Node) checkAccess(ctx context.Context, method string, accessRoles []string) error {

	role, err := n.Action.ActionConfig.DefaultRoleFunc()(ctx, config.RoleReq{
		AccessName: n.tableName,
		Method:     method,
		NodeRole:   n.Role,
	})

	if err != nil {
		return err
	}

	if role == consts.DENY {
		return consts.NewDenyErr(n.Key, n.Role)
	}

	n.Role = role

	if !lo.Contains(accessRoles, role) {
		return consts.NewNoAccessErr(n.Key, n.Role)
	}

	return nil
}

func (n *Node) whereUpdate(ctx context.Context, method string, accessRoles []string) error {

	for i, item := range n.req {

		condition := config.NewConditionRet()

		req := model.Map{}

		for k, v := range item {
			k := n.Action.DbFieldStyle(ctx, n.RowKey, k)
			req[k] = v
		}

		conditionReq := config.ConditionReq{
			AccessName:          n.Key,
			TableAccessRoleList: accessRoles,
			Method:              method,
			NodeRole:            n.Role,
			NodeReq:             req,
		}

		err := n.Action.ActionConfig.ConditionFunc(ctx, conditionReq, condition)

		if err != nil {
			return err
		}

		if method == http.MethodPost {
			for k, v := range condition.AllWhere() {
				n.Data[i][k] = v
			}
		} else {
			for k, v := range condition.AllWhere() {
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
				return consts.NewStructureKeyNoFoundErr(n.Key + "." + key)
			}
		}

		// refuse
		if len(n.structure.Refuse) > 0 && n.structure.Refuse[0] == "!" {
			if len(n.structure.Must) == 0 {
				return consts.NewValidStructureErr("REFUSE为!时必须指定MUST:" + n.Key)
			}

			for key, _ := range item {
				if !lo.Contains(n.structure.Must, key) {
					return consts.NewValidStructureErr("不能包含:" + n.Key + "." + key)
				}
			}

		} else {
			for _, key := range n.structure.Refuse {
				if _, exists := item[key]; exists {
					return consts.NewValidStructureErr("不能包含:" + n.Key + "." + key)
				}
			}
		}
	}

	return nil
}

// reqUpdate 处理 Update/Insert等
func (n *Node) reqUpdate() error {

	for i, _ := range n.req {
		for key, updateVal := range n.structure.Update {

			if strings.HasSuffix(key, consts.FunctionsKeySuffix) {

				k := key[0 : len(key)-2]

				// call functions
				{
					actionConfig := n.Action.ActionConfig

					functionName, paramKeys := util.ParseFunctionsStr(updateVal.(string))

					_func := actionConfig.Func(functionName)

					param := model.Map{}
					for paramI, item := range _func.ParamList {
						if item.Name == consts.FunctionOriReqParam {
							param[item.Name] = n.Data[i]
						} else {
							param[item.Name] = n.Data[i][paramKeys[paramI]]
						}
					}

					val, err := actionConfig.CallFunc(n.ctx, functionName, param)
					if err != nil {
						return err
					}
					if val != nil {
						n.Data[i][k] = val
					}
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

	}

	return nil
}

// reqUpdate 处理 Update/Insert等  (事务内)
func (n *Node) reqUpdateBeforeDo() error {

	for i, _ := range n.req {

		for k, v := range n.Data[i] {
			// 处理 ref
			if strings.HasSuffix(k, consts.RefKeySuffix) {
				refNodeKey, refCol := util.ParseRefCol(v.(string))
				if strings.HasSuffix(refNodeKey, consts.ListKeySuffix) { // 双列表
					n.Data[i][k] = n.keyNode[refNodeKey].Data[i][n.Action.DbFieldStyle(n.ctx, n.tableName, refCol)]
				} else {
					n.Data[i][k] = n.keyNode[refNodeKey].Data[0][n.Action.DbFieldStyle(n.ctx, n.tableName, refCol)]
				}
			}
		}
	}

	return nil
}

func (n *Node) do(ctx context.Context, method string) (ret model.Map, err error) {

	var rowKeyVal model.Map
	var rowKey string

	access, err := n.Action.ActionConfig.GetAccessConfig(n.Key, true)
	if err != nil {
		return nil, err
	}

	switch method {
	case http.MethodPost:

		if access.RowKeyGen != "" {
			for i, _ := range n.Data {

				rowKeyVal, err = n.Action.ActionConfig.RowKeyGen(ctx, access.RowKeyGen, n.Key, n.tableName, n.Data[i])
				if err != nil {
					return nil, gerror.Wrap(err, "RowKeyGen")
				}

				for k, v := range rowKeyVal {
					if k == consts.RowKey {
						n.Data[i][access.RowKey] = v
					} else {
						n.Data[i][k] = v
					}
				}

			}
		}

		rowKey = access.RowKey

	case http.MethodPut:
	case http.MethodDelete:

	default:
		return nil, consts.NewMethodNotSupportErr(method)
	}

	executor, err := GetActionExecutor(n.executor)
	if err != nil {
		return nil, err
	}

	ret, err = executor.Do(ctx, ActionExecutorReq{
		Method:   method,
		Table:    n.tableName,
		Data:     n.Data,
		Where:    n.Where,
		Access:   access,
		Config:   n.Action.ActionConfig,
		NewQuery: n.Action.NewQuery,
	})

	if err != nil {
		return nil, err
	}

	if len(n.Data) == 1 {

		jsonStyle := n.Action.JsonFieldStyle
		if rowKeyVal != nil {
			for k, v := range rowKeyVal {
				if k == consts.RowKey {
					ret[jsonStyle(ctx, n.tableName, rowKey)] = v
				} else {
					ret[jsonStyle(ctx, n.tableName, k)] = v
				}
			}
		}
	}

	n.Ret = ret

	return
}

func (n *Node) execute(ctx context.Context, method string) (model.Map, error) {

	err := n.reqUpdateBeforeDo()
	if err != nil {
		return nil, err
	}

	ret, err := n.do(ctx, method)
	if err != nil {
		return nil, err
	}

	n.Ret = ret
	return n.Ret, nil
}
