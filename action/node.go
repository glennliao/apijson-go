package action

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

type Node struct {
	req       []model.Map
	ctx       context.Context
	action    *Action
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

	//access *config.Access
}

func newNode(key string, req []model.Map, structure *config.Structure, executor string) Node {

	n := Node{
		Key: key, req: req, structure: structure, executor: executor,
	}

	if strings.HasSuffix(key, consts.ListKeySuffix) {
		n.Key = key[0 : len(key)-len(consts.ListKeySuffix)]
		n.IsList = true
	}

	return n
}

// parse req data to data/where
func (n *Node) parseReq(method string) {
	n.Data = []model.Map{}
	n.Where = []model.Map{}

	for i, item := range n.req {

		n.Data = append(n.Data, model.Map{})
		n.Where = append(n.Where, model.Map{})

		for key, val := range item {

			if key == consts.Role {
				n.Role = util.String(val)
				continue
			}

			key = n.action.DbFieldStyle(n.ctx, n.tableName, key)

			switch method {
			case http.MethodPost:
				n.Data[i][key] = val
			case http.MethodDelete:
				n.Where[i][key] = val
			case http.MethodPut:
				if key == n.RowKey || key == n.RowKey+"{}" {
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
		key = key[0 : len(key)-2]
	}
	access, err := n.action.actionConfig.GetAccessConfig(key, true)

	if err != nil {
		return err
	}

	n.tableName = access.Name
	n.RowKey = access.RowKey

	n.parseReq(method)

	// 0. 角色替换

	err = n.roleUpdate()
	if err != nil {
		return err
	}
	var accessRoles []string
	if n.action.NoAccessVerify == false {
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

	role, err := n.action.actionConfig.DefaultRoleFunc()(ctx, config.RoleReq{
		AccessName: n.tableName,
		Method:     method,
		NodeRole:   n.Role,
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

	return nil
}

func (n *Node) whereUpdate(ctx context.Context, method string, accessRoles []string) error {

	for i, item := range n.req {

		condition := config.NewConditionRet()

		conditionReq := config.ConditionReq{
			AccessName:          n.Key,
			TableAccessRoleList: accessRoles,
			Method:              method,
			NodeRole:            n.Role,
			NodeReq:             item,
		}

		err := n.action.actionConfig.ConditionFunc(ctx, conditionReq, condition)

		if err != nil {
			return err
		}

		if method == http.MethodPost {
			for k, v := range condition.Where() {
				n.Data[i][k] = v
			}
		} else {
			for k, v := range condition.Where() {
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
		if len(n.structure.Refuse) > 0 && n.structure.Refuse[0] == "!" {
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

// reqUpdate 处理 Update/Insert等
func (n *Node) reqUpdate() error {

	for i, _ := range n.req {
		for key, updateVal := range n.structure.Update {

			if strings.HasSuffix(key, consts.FunctionsKeySuffix) {

				k := key[0 : len(key)-2]

				// call functions
				{
					queryConfig := n.action.actionConfig

					functionName, paramKeys := util.ParseFunctionsStr(updateVal.(string))

					_func := queryConfig.Func(functionName)

					param := model.Map{}
					for paramI, item := range _func.ParamList {
						if item.Name == consts.FunctionOriReqParam {
							param[item.Name] = n.Data[i]
						} else {
							param[item.Name] = n.Data[i][paramKeys[paramI]]
						}
					}

					val, err := _func.Handler(n.ctx, param)
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
					n.Data[i][k] = n.keyNode[refNodeKey].Data[i][n.action.DbFieldStyle(n.ctx, n.tableName, refCol)]
				} else {
					n.Data[i][k] = n.keyNode[refNodeKey].Data[0][n.action.DbFieldStyle(n.ctx, n.tableName, refCol)]
				}
			}
		}
	}

	return nil
}

func (n *Node) do(ctx context.Context, method string) (ret model.Map, err error) {

	err = EmitHook(ctx, BeforeExecutorDo, n, method)
	if err != nil {
		return nil, err
	}

	switch method {
	case http.MethodPost:

		var rowKeyVal model.Map

		access, err := n.action.actionConfig.GetAccessConfig(n.Key, true)
		if err != nil {
			return nil, err
		}

		if access.RowKeyGen != "" {
			for i, _ := range n.Data {

				rowKeyVal, err = n.action.actionConfig.RowKeyGen(ctx, access.RowKeyGen, n.Key, n.Data[i])
				if err != nil {
					return nil, err
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

		ret, err := executor.GetActionExecutor(n.executor).Do(ctx, executor.ActionExecutorReq{
			Method: method,
			Table:  n.tableName,
			Data:   n.Data,
			Where:  nil,
		})

		if err != nil {
			return nil, err
		}

		if len(n.Data) > 0 { //多条插入时返回值已经应该无意义了

			jsonStyle := n.action.JsonFieldStyle
			if rowKeyVal != nil {
				for k, v := range rowKeyVal {
					if k == consts.RowKey {
						ret[jsonStyle(ctx, n.tableName, access.RowKey)] = v
					} else {
						ret[jsonStyle(ctx, n.tableName, k)] = v
					}
				}
			}
		}

	case http.MethodPut:
	case http.MethodDelete:

	default:
		return nil, gerror.New("undefined method:" + method)
	}

	ret, err = executor.GetActionExecutor(n.executor).Do(ctx, executor.ActionExecutorReq{
		Method: method,
		Table:  n.tableName,
		Data:   n.Data,
		Where:  n.Where,
	})

	if err != nil {
		return nil, err
	}

	n.Ret = ret

	err = EmitHook(ctx, AfterExecutorDo, n, method)
	if err != nil {
		return nil, err
	}

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
