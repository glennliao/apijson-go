package gf_orm

import (
	"context"
	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

type ActionExecutor struct {
	DbName string
}

func (a *ActionExecutor) Insert(ctx context.Context, table string, data any) (id int64, count int64, err error) {
	ret, err := g.DB(a.DbName).Insert(ctx, table, data)
	if err != nil {
		return 0, 0, err
	}
	id, err = ret.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	count, err = ret.RowsAffected()
	return id, count, nil
}

func (a *ActionExecutor) Update(ctx context.Context, table string, data model.Map, where model.Map) (count int64, err error) {
	m := g.DB(a.DbName).Model(table).Ctx(ctx)

	for k, v := range where {
		if strings.HasSuffix(k, "{}") {
			if vStr, ok := v.(string); ok {
				if vStr == "" {
					return 0, gerror.New("where的值不能为空")
				}
			}
			m = m.WhereIn(k[0:len(k)-2], v)
			delete(where, k)
			continue
		}
		if v.(string) == "" || v == nil { //暂只处理字符串为空的情况
			return 0, gerror.New("where的值不能为空")
		}
	}

	_ret, err := m.Where(where).Update(data)
	if err != nil {
		return 0, err
	}

	count, err = _ret.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err
}

func (a *ActionExecutor) Delete(ctx context.Context, table string, where model.Map) (count int64, err error) {
	if len(where) == 0 {
		return 0, gerror.New("where不能为空")
	}

	m := g.DB(a.DbName).Model(table).Ctx(ctx)

	for k, v := range where {
		if strings.HasSuffix(k, "{}") {
			m = m.WhereIn(k[0:len(k)-2], v)
			delete(where, k)
			continue
		}
		if v.(string) == "" || v == nil { //暂只处理字符串为空的情况
			return 0, gerror.New("where的值不能为空")
		}
	}

	_ret, err := m.Where(where).Delete()
	if err != nil {
		return 0, err
	}

	count, err = _ret.RowsAffected()

	if err != nil {
		return 0, err
	}

	return count, err
}

// init 暂先自动注册,后续改成可手动配置
func init() {
	executor.RegActionExecutor("default", &ActionExecutor{})
}
