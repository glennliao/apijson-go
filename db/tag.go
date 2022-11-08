package db

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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

	_ret, err := g.DB().Update(ctx, table, data, where)
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

	for _, v := range where {
		if v.(string) == "" || v == nil { //暂只处理字符串为空的情况
			return 0, gerror.New("where的值不能为空")
		}
	}

	_ret, err := g.DB().Model(table).Ctx(ctx).Delete(where)
	if err != nil {
		return 0, err
	}

	count, err := _ret.RowsAffected()

	if err != nil {
		return 0, err
	}

	return count, err

}
