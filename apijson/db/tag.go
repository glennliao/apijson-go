package db

import (
	"context"
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

func Update(ctx context.Context, table string, data g.Map) (any, int64, error) {
	// table = gstr.CaseSnake(table)

	rowKey := RowKeyMap[table]
	id := data[rowKey]

	_ret, err := g.DB().Update(ctx, table, data, g.Map{rowKey: id})
	if err != nil {
		return id, 0, err
	}

	count, err := _ret.RowsAffected()
	if err != nil {
		return id, 0, err
	}

	return id, count, err

}

func Delete(ctx context.Context, table string, data g.Map) (any, int64, error) {
	// table = gstr.CaseSnake(table)

	rowKey := RowKeyMap[table]
	id := data[rowKey]

	if _, ok := data[rowKey+"{}"]; ok {
		id = data[rowKey+"{}"]
	}

	_ret, err := g.DB().Model(table).Ctx(ctx).Delete(g.Map{rowKey: id})
	if err != nil {
		return id, 0, err
	}

	count, err := _ret.RowsAffected()
	if err != nil {
		return id, 0, err
	}

	return id, count, err

}
