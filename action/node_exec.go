package action

import (
	"context"
	"net/http"
	"strings"

	"github.com/glennliao/apijson-go/config"

	"github.com/glennliao/apijson-go/util"

	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (n *Node) execute(ctx context.Context, method string) (model.Map, error) {
	err := n.updateReqBeforeDo()
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

// updateReqBeforeDo
// like ref -> 'xxx@'
func (n *Node) updateReqBeforeDo() error {
	for key, v := range n.Data {
		// ref
		if strings.HasSuffix(key, consts.RefKeySuffix) {
			refNodeKey, refCol := util.ParseRefCol(v.(string))
			refNodeKeyCol := n.Action.DbFieldStyle(n.ctx, n.TableName, refCol)

			n.Data[key] = n.keyNode[refNodeKey].Data[refNodeKeyCol]

			//if strings.HasSuffix(refNodeKey, consts.ListKeySuffix) { // 双列表
			//	n.Data[key] = n.keyNode[refNodeKey].Data[n.Action.DbFieldStyle(n.ctx, n.TableName, refCol)]
			//}
		}
	}
	return nil
}

func (n *Node) do(ctx context.Context, method string) (ret model.Map, err error) {
	access, err := n.Action.ActionConfig.GetAccessConfig(n.Key, n.Action.NoAccessVerify)
	if err != nil {
		return nil, err
	}

	var (
		rowKeyVal model.Map
		rowKey    = access.RowKey
	)

	// gen rowKey
	if method == http.MethodPost && access.RowKeyGen != "" {

		rowKeyVal, err = n.genRowKey(access)
		if err != nil {
			return nil, err
		}
	}

	executor, err := GetActionExecutor(n.executor)
	if err != nil {
		return nil, err
	}

	ret, err = executor.Do(ctx, ExecutorReq{
		Method:          method,
		Table:           n.TableName,
		Data:            n.Data,
		Where:           n.Where,
		Access:          access,
		AccessCondition: n.AccessCondition,
		Config:          n.Action.ActionConfig,
	})
	if err != nil {
		return nil, err
	}

	// return rowKey
	if rowKeyVal != nil {
		jsonStyle := n.Action.JsonFieldStyle
		for k, v := range rowKeyVal {
			if k == consts.RowKey {
				ret[jsonStyle(ctx, n.TableName, rowKey)] = v
			} else {
				ret[jsonStyle(ctx, n.TableName, k)] = v
			}
		}
	}

	n.Ret = ret

	return
}

func (n *Node) genRowKey(access *config.AccessConfig) (rowKeyVal model.Map, err error) {
	rowKeyVal, err = n.Action.ActionConfig.RowKeyGen(n.ctx, access.RowKeyGen, n.Key, n.TableName, n.Data)
	if err != nil {
		return nil, gerror.Wrap(err, "RowKeyGen")
	}

	for k, v := range rowKeyVal {
		if k == consts.RowKey {
			n.Data[access.RowKey] = v
		} else {
			n.Data[k] = v
		}
	}
	return
}

// reqUpdate 处理 structure 的 Update/Insert等
//func (n *Node) reqUpdate() error {
//	for i := range n.req {
//		for key, updateVal := range n.structure.Update {
//			if strings.HasSuffix(key, consts.FunctionsKeySuffix) {
//
//				k := key[0 : len(key)-2]
//
//				// call functions
//				{
//					actionConfig := n.Action.ActionConfig
//
//					functionName, paramKeys := util.ParseFunctionsStr(updateVal.(string))
//
//					_func := actionConfig.Func(functionName)
//
//					param := model.Map{}
//					for paramI, item := range _func.ParamList {
//						if item.Name == consts.FunctionOriReqParam {
//							param[item.Name] = n.Data[i]
//						} else {
//							param[item.Name] = n.Data[i][paramKeys[paramI]]
//						}
//					}
//
//					val, err := actionConfig.CallFunc(n.ctx, functionName, param)
//					if err != nil {
//						return err
//					}
//					if val != nil {
//						n.Data[i][k] = val
//					}
//				}
//
//			} else {
//				n.Data[i][key] = updateVal
//			}
//		}
//
//		for key, updateVal := range n.structure.Insert {
//			if _, exists := n.Data[i][key]; !exists {
//				n.Data[i][key] = updateVal
//			}
//		}
//
//	}
//
//	return nil
//}
