package action

import (
	"context"
	"net/http"

	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
)

// parse node
func (n *Node) parse() error {
	ctx, method := n.ctx, n.Action.Method

	access, err := n.Action.ActionConfig.GetAccessConfig(n.Key, n.Action.NoAccessVerify)
	if err != nil {
		return err
	}

	n.TableName = access.Name
	n.RowKey = access.RowKey

	accessRoles := access.GetAccessRoles(method)

	// 0. 角色替换

	//err = n.roleUpdate()
	//if err != nil {
	//	return err
	//}

	if !n.Action.NoAccessVerify {
		err = CheckRoleAccess(ctx, n, accessRoles)
		if err != nil {
			return err
		}
	}

	err = CheckReqStructure(n.req, n.structure)
	if err != nil {
		return err
	}

	req := model.Map{}
	for key, v := range n.req {
		k := n.Action.DbFieldStyle(ctx, n.TableName, key)
		req[k] = v
	}
	n.req = req

	err = n.parseAccessCondition(ctx, accessRoles)
	if err != nil {
		return err
	}

	n.parseReq()

	return err
}

// parseAccessCondition
func (n *Node) parseAccessCondition(ctx context.Context, accessRoles []string) error {
	n.AccessCondition = &config.ConditionRet{}

	conditionReq := config.ConditionReq{
		AccessName:          n.Key,
		TableAccessRoleList: accessRoles,
		Method:              n.Action.Method,
		NodeRole:            n.Role,
		NodeReq:             n.req,
	}

	err := n.Action.ActionConfig.ConditionFunc(ctx, conditionReq, n.AccessCondition)

	return err
}

// parse req to data/where
func (n *Node) parseReq() {
	for key, val := range n.req {

		if key == consts.Role {
			n.Role = util.String(val)
			continue
		}

		// key = n.Action.DbFieldStyle(n.ctx, n.TableName, key)

		switch n.Action.Method {
		case http.MethodPost:
			n.Data[key] = val
		case http.MethodDelete:
			n.Where[key] = val
		case http.MethodPut:
			if key == n.RowKey || key == n.RowKey+consts.OpIn {
				// only rowKey is where, others is update data
				n.Where[key] = val
			} else {
				n.Data[key] = val
			}
		}

	}
}

// update node role
//func (n *Node) roleUpdate() error {
//	if val, exists := n.structure.Insert[consts.Role]; exists {
//		if n.Role == "" {
//			n.Role = util.String(val)
//		}
//	}
//
//	if val, exists := n.structure.Update[consts.Role]; exists {
//		n.Role = util.String(val)
//	}
//	return nil
//}
