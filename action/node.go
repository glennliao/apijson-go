package action

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/db"
	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/config/functions"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"net/http"
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

	structure *db.Structure
	executor  string

	keyNode map[string]*Node
}

func newNode(key string, req []g.Map, structure *db.Structure, executor string) Node {
	return Node{
		Key: key, req: req, structure: structure, executor: executor,
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

				if method == http.MethodDelete {
					n.Where[i][key] = val
				} else {
					if key == n.RowKey || key == n.RowKey+"{}" {
						if method == http.MethodPut {
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
		key = key[0 : len(key)-2]
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

		if method == http.MethodPost {
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
				var param = g.Map{}
				for _, paramKey := range paramKeys {
					if paramKey == consts.FunctionOriReqParam {
						param[paramKey] = n.Data[i]
					} else {
						param[paramKey] = n.Data[i][paramKey]
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

	}

	return nil
}

// reqUpdate 处理 Update/Insert等  (事务内)
func (n *Node) reqUpdateBeforeDo() error {

	for i, _ := range n.req {

		for k, v := range n.Data[i] {
			if strings.HasSuffix(k, consts.RefKeySuffix) {
				refNodeKey, refCol := util.ParseRefCol(v.(string))
				if strings.HasSuffix(refNodeKey, consts.ListKeySuffix) { // 双列表
					n.Data[i][k] = n.keyNode[refNodeKey].Data[i][config.GetDbFieldStyle()(n.ctx, n.TableName, refCol)]
				} else {
					n.Data[i][k] = n.keyNode[refNodeKey].Data[0][config.GetDbFieldStyle()(n.ctx, n.TableName, refCol)]
				}
			}
		}
	}

	return nil
}

func (n *Node) do(ctx context.Context, method string, dataIndex int) (ret g.Map, err error) {

	err = EmitHook(ctx, BeforeDo, n, method)
	if err != nil {
		return nil, err
	}

	var count int64

	switch method {
	case http.MethodPost:

		var rowKeyVal g.Map

		access, err := db.GetAccess(n.Key, true)
		if err != nil {
			return nil, err
		}

		if access.RowKeyGen != "" {
			for i, _ := range n.Data {

				rowKeyVal, err = config.RowKeyGen(ctx, access.RowKeyGen, n.TableName, n.Data[i])
				if err != nil {
					return nil, err
				}

				for k, v := range rowKeyVal {
					if k == "rowKey" {
						n.Data[i][access.RowKey] = v
					} else {
						n.Data[i][k] = v
					}
				}

			}
		}

		var id int64

		id, count, err = executor.GetActionExecutor(n.executor).Insert(ctx, n.TableName, n.Data)

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
					if k == "rowKey" {
						ret[jsonStyle(ctx, n.TableName, access.RowKey)] = v
					} else {
						ret[jsonStyle(ctx, n.TableName, k)] = v
					}
				}
			}
		}

	case http.MethodPut:
		count, err = executor.GetActionExecutor(n.executor).Update(ctx, n.TableName, n.Data[dataIndex], n.Where[dataIndex])
		if err != nil {
			return nil, err
		}

		ret = g.Map{
			"code":  200,
			"count": count,
		}
	case http.MethodDelete:
		count, err = executor.GetActionExecutor(n.executor).Delete(ctx, n.TableName, n.Where[dataIndex])
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

	err = EmitHook(ctx, AfterDo, n, method)
	if err != nil {
		return nil, err
	}

	return
}

func (n *Node) execute(ctx context.Context, method string) (g.Map, error) {

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

	return g.Map{
		"code": 200,
	}, nil
}
