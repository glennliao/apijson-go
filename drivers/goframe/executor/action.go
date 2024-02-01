package executor

import (
	"context"
	"net/http"
	"strings"

	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
)

type ActionExecutor struct {
	DbResolver DbResolver
}

func (a *ActionExecutor) Do(ctx context.Context, req action.ExecutorReq) (ret model.Map, err error) {
	switch req.Method {
	case http.MethodPost:
		return a.Insert(ctx, req.Table, req.Data)
	case http.MethodPut:

		ret, err = a.Update(ctx, req.Table, req.Data, req.Where, req.AccessCondition)
		if err != nil {
			ret = model.Map{
				"code": 200,
			}
		}
		return ret, err

	case http.MethodDelete:

		ret, err = a.Delete(ctx, req.Table, req.Where, req.AccessCondition)
		if err != nil {
			ret = model.Map{
				"code": 200,
			}
		}
		return ret, err
	}
	return nil, consts.NewMethodNotSupportErr(req.Method)
}

func (a *ActionExecutor) Insert(ctx context.Context, table string, data model.Map) (ret model.Map, err error) {
	result, err := a.DbResolver(ctx).Insert(ctx, table, data)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	count, err := result.RowsAffected()

	ret = model.Map{
		"code":  200,
		"count": count,
		"id":    id,
	}

	return ret, nil
}

func (a *ActionExecutor) Update(ctx context.Context, table string, data model.Map, where model.Map, accessCondition *config.ConditionRet) (ret model.Map, err error) {
	conditions, builders := accessCondition.AllCondition()
	//if len(conditions)+len(builders) == 0 {
	//	return nil, consts.NewValidReqErr("where的值不能为空")
	//}

	m := a.DbResolver(ctx).Model(table).Ctx(ctx)

	whereBuilder := m.Builder()

	for k, v := range where {
		// TODO 统一到condition中
		if strings.HasSuffix(k, consts.OpIn) {
			if vStr, ok := v.(string); ok {
				if vStr == "" {
					return nil, consts.NewValidReqErr("where的值不能为空")
				}
			}
			m = m.WhereIn(k[0:len(k)-2], v)
			delete(where, k)
			continue
		} else if k == consts.Raw {
			m = m.Where(v.(map[string][]any))
			delete(where, k)
			continue
		} else {
			whereBuilder = whereBuilder.Where(k, v)
		}
	}

	for _, v := range conditions {

		whereBuilder = whereBuilder.Where(v.Column, v.Args)

		if strings.Contains(v.Column, "?") {
			// TODO 校验必须有条件地删除/更新
			//if len(v.Args) == 0 || gconv.String(v.Args[0]) == "" { // 暂只处理字符串为空的情况, // TODO 此处 args[n]? 存在的可能?
			//	return nil, consts.NewValidReqErr("where的值不能为空:" + v.Column)
			//}
		}
	}

	// 子查询, 看下如何与上面的统一
	for _, v := range builders {
		whereBuilder = whereBuilder.Where(v)
	}

	for k, v := range data {
		if strings.HasSuffix(k, consts.OpPLus) {
			field := util.RemoveSuffix(k, consts.OpPLus)
			data[field] = &gdb.Counter{
				Field: field,
				Value: gconv.Float64(v),
			}
			delete(data, k)
			continue
		}
		if strings.HasSuffix(k, consts.OpSub) {
			field := util.RemoveSuffix(k, consts.OpSub)
			data[field] = &gdb.Counter{
				Field: field,
				Value: -gconv.Float64(v),
			}
			delete(data, k)
			continue
		}

		if data[k] == nil {
			// 此处目前不允许外部设置null
			delete(data, k)
		}

	}

	_ret, err := m.Where(whereBuilder).Update(data)
	if err != nil {
		return nil, err
	}

	count, err := _ret.RowsAffected()
	if err != nil {
		return nil, err
	}

	ret = model.Map{
		"code":  200,
		"count": count,
	}

	return ret, err
}

func (a *ActionExecutor) Delete(ctx context.Context, table string, where model.Map, accessCondition *config.ConditionRet) (ret model.Map, err error) {
	// TODO access 校验判断
	//conditions, builders := accessCondition.AllCondition()
	//if len(conditions)+len(builders) == 0 {
	//	return nil, consts.NewValidReqErr("where的值不能为空")
	//}

	m := a.DbResolver(ctx).Model(table).Ctx(ctx)

	whereBuilder := m.Builder()

	for k, v := range where {
		if strings.HasSuffix(k, consts.OpIn) {
			if vStr, ok := v.(string); ok {
				if vStr == "" {
					return nil, consts.NewValidReqErr("where的值不能为空")
				}
			}
			m = m.WhereIn(k[0:len(k)-2], v)
			delete(where, k)
			continue
		} else if k == consts.Raw {
			m = m.Where(v.(map[string][]any))
			delete(where, k)
			continue
		} else {
			whereBuilder = whereBuilder.Where(k, v)
		}
	}

	//for _, v := range conditions {
	//
	//	whereBuilder = whereBuilder.Where(v.Column, v.Args)
	//
	//	if strings.Contains(v.Column, "?") {
	//		//if len(v.Args) == 0 || gconv.String(v.Args[0]) == "" { // 暂只处理字符串为空的情况, // TODO 此处 args[n]? 存在的可能?
	//		//	return nil, consts.NewValidReqErr("where的值不能为空:" + v.Column)
	//		//}
	//	}
	//
	//}

	//for _, v := range builders {
	//	whereBuilder = whereBuilder.Where(v)
	//}

	_ret, err := m.Where(whereBuilder).Delete()
	if err != nil {
		return nil, err
	}

	count, err := _ret.RowsAffected()
	if err != nil {
		return nil, err
	}

	ret = model.Map{
		"code":  200,
		"count": count,
	}

	return ret, err
}
