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
	req        []model.Map
	ctx        context.Context
	action     *Action
	Key        string
	tableName  string
	AccessName string
	Role       string

	Data   []model.Map // 需写入数据库的数据
	Where  []model.Map // 条件
	RowKey string      // 主键

	structure *config.Structure
	executor  string

	keyNode map[string]*Node

	//access *config.Access
}

func newNode(key string, req []model.Map, structure *config.Structure, executor string) Node {

	accessName := key
	if strings.HasSuffix(accessName, "[]") {
		accessName = accessName[0 : len(accessName)-2]
	}

	return Node{
		Key: key, req: req, structure: structure, executor: executor, AccessName: accessName,
	}
}

func (n *Node) parseReq(method string) {
	n.Data = []model.Map{}
	n.Where = []model.Map{}

	for i, item := range n.req {

		n.Data = append(n.Data, model.Map{})
		n.Where = append(n.Where, model.Map{})

		for key, val := range item {
			if key == consts.Role {
				n.Role = util.String(val)
			} else {
				key = n.action.DbFieldStyle(n.ctx, n.tableName, key)

				if method == http.MethodDelete {
					n.Where[i][key] = val
				} else {
					if key == n.RowKey || key == n.RowKey+"{}" {
						if method == http.MethodPut {
							n.Where[i][key] = val
						} else {
							n.Data[i][key] = val
						}
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

	if n.action.NoAccessVerify == false {
		// 1. 检查权限, 无权限就不用做参数检查了
		var accessRoles []string

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

	return nil
}

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

	for i, item := range n.req {

		condition := config.NewConditionRet()

		conditionReq := config.ConditionReq{
			AccessName:          n.tableName,
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
				functionName, paramKeys := util.ParseFunctionsStr(updateVal.(string))
				var param = model.Map{}
				for _, paramKey := range paramKeys {
					if paramKey == consts.FunctionOriReqParam {
						param[paramKey] = n.Data[i]
					} else {
						param[paramKey] = n.Data[i][paramKey]
					}
				}
				k := key[0 : len(key)-2]
				val, err := n.action.Functions.Call(n.ctx, functionName, param)
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

func (n *Node) do(ctx context.Context, method string, dataIndex int) (ret model.Map, err error) {

	err = EmitHook(ctx, BeforeExecutorDo, n, method)
	if err != nil {
		return nil, err
	}

	var count int64

	switch method {
	case http.MethodPost:

		var rowKeyVal model.Map

		access, err := n.action.actionConfig.GetAccessConfig(n.Key, true)
		if err != nil {
			return nil, err
		}

		if access.RowKeyGen != "" {
			for i, _ := range n.Data {

				rowKeyVal, err = n.action.actionConfig.RowKeyGen(ctx, access.RowKeyGen, n.AccessName, n.Data[i])
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

		var id int64

		id, count, err = executor.GetActionExecutor(n.executor).Insert(ctx, n.tableName, n.Data)

		if err != nil {
			return nil, err
		}

		ret = model.Map{
			"code":  200,
			"count": count,
			"id":    id,
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
		count, err = executor.GetActionExecutor(n.executor).Update(ctx, n.tableName, n.Data[dataIndex], n.Where[dataIndex])
		if err != nil {
			return nil, err
		}

		ret = model.Map{
			"code":  200,
			"count": count,
		}
	case http.MethodDelete:
		count, err = executor.GetActionExecutor(n.executor).Delete(ctx, n.tableName, n.Where[dataIndex])
		if err != nil {
			return nil, err
		}

		ret = model.Map{
			"code":  200,
			"count": count,
		}
	}

	if ret == nil {
		return nil, gerror.New("undefined method:" + method)
	}

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

	if method == http.MethodPost { // 新增时可以合并新增
		ret, err := n.do(ctx, method, 0)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		for i, _ := range n.req {
			_, err := n.do(ctx, method, i)
			if err != nil {
				return nil, err
			}
		}
	}

	return model.Map{
		"code": 200,
	}, nil
}
