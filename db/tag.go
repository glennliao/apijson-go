package db

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func Insert(ctx context.Context, table string, data any) (int64, int64, error) {

	access, ok := accessMap[table]
	if !ok {
		panic(gerror.New("table 不存在:" + table))
	}

	table = access.Name

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

	access, ok := accessMap[table]
	if !ok {
		panic(gerror.New("table 不存在:" + table))
	}

	table = access.Name

	rowKey := "id"
	id := data[rowKey]

	where := g.Map{rowKey: id}
	for k, v := range data {
		if k == "_where" {

			_v := v.(g.Map)
			for __k, __v := range _v {
				where[__k] = __v
			}

		}
	}

	delete(data, "_where")
	delete(data, rowKey)

	_ret, err := g.DB().Update(ctx, table, data, where)
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
	access, ok := accessMap[table]
	if !ok {
		panic(gerror.New("table 不存在:" + table))
	}

	table = access.Name

	rowKey := "id"
	id := data[rowKey]

	if _, ok := data[rowKey+"{}"]; ok {
		id = data[rowKey+"{}"]
	}

	where := g.Map{rowKey: id}
	for k, v := range data {
		where[k] = v
	}

	_ret, err := g.DB().Model(table).Ctx(ctx).Delete(where)
	if err != nil {
		return id, 0, err
	}

	count, err := _ret.RowsAffected()
	if err != nil {
		return id, 0, err
	}

	return id, count, err

}
