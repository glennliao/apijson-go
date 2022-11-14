package db

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

func Insert(ctx context.Context, table string, data any) (int64, int64, error) {

	ret, err := g.DB().Insert(ctx, table, data)
	if err != nil {
		return 0, 0, err
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	count, err := ret.RowsAffected()
	return id, count, nil
}

func Update(ctx context.Context, table string, data g.Map, where g.Map) (int64, error) {

	m := g.DB().Model(table).Ctx(ctx)

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

	count, err := _ret.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err

}

func Delete(ctx context.Context, table string, where g.Map) (int64, error) {

	if len(where) == 0 {
		return 0, gerror.New("where不能为空")
	}

	m := g.DB().Model(table).Ctx(ctx)

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

	count, err := _ret.RowsAffected()

	if err != nil {
		return 0, err
	}

	return count, err

}
