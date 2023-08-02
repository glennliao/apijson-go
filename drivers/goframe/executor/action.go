package executor

import (
	"context"
	"net/http"
	"strings"

	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
)

type ActionExecutor struct {
	DbResolver DbResolver
}

func (a *ActionExecutor) Do(ctx context.Context, req action.ActionExecutorReq) (ret model.Map, err error) {
	switch req.Method {
	case http.MethodPost:
		return a.Insert(ctx, req.Table, req.Data)
	case http.MethodPut:

		for i, _ := range req.Data {
			ret, err = a.Update(ctx, req.Table, req.Data[i], req.Where[i])
			if err != nil {
				break
			}
		}
		if err != nil {
			ret = model.Map{
				"code": 200,
			}
		}
		return ret, err

	case http.MethodDelete:

		for i, _ := range req.Data {
			ret, err = a.Delete(ctx, req.Table, req.Where[i])
			if err != nil {
				break
			}
		}
		if err != nil {
			ret = model.Map{
				"code": 200,
			}
		}
		return ret, err
	}
	return nil, consts.NewMethodNotSupportErr(req.Method)
}

func (a *ActionExecutor) Insert(ctx context.Context, table string, data []model.Map) (ret model.Map, err error) {
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

func (a *ActionExecutor) Update(ctx context.Context, table string, data model.Map, where model.Map) (ret model.Map, err error) {
	m := a.DbResolver(ctx).Model(table).Ctx(ctx)

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
		}
		if k == consts.Raw {
			m = m.Where(v.(map[string][]any))
			delete(where, k)
			continue
		}

		if v == nil || gconv.String(v) == "" { // 暂只处理字符串为空的情况
			return nil, consts.NewValidReqErr("where的值不能为空:" + k)
		}
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

	_ret, err := m.Where(where).Update(data)
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

func (a *ActionExecutor) Delete(ctx context.Context, table string, where model.Map) (ret model.Map, err error) {
	if len(where) == 0 {
		return nil, consts.NewValidReqErr("where的值不能为空")
	}

	m := a.DbResolver(ctx).Model(table).Ctx(ctx)

	for k, v := range where {

		if k == consts.Raw {
			m = m.Where(v)
			continue
		}

		if strings.HasSuffix(k, consts.OpIn) {
			m = m.WhereIn(k[0:len(k)-2], v)
			delete(where, k)
			continue
		}

		if gconv.String(v) == "" || v == nil { // 暂只处理字符串为空的情况
			return nil, consts.NewValidReqErr("where的值不能为空:" + k)
		}

		m = m.Where(k, v)
	}

	_ret, err := m.Delete()
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
