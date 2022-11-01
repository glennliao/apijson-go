package db

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	Column struct {
		// 字段名
		Name string
	}

	Table struct {
		// 表名
		Name    string
		Columns []Column
	}
)

var tableMap = map[string]Table{}

func loadTableMeta() {
	var ctx = context.TODO()

	_tableMap := make(map[string]Table)

	tables, err := g.DB().Tables(ctx)
	if err != nil {
		panic(err)
	}

	for _, table := range tables {
		fields, err := g.DB().TableFields(ctx, table)
		if err != nil {
			panic(err)
		}

		var columns []Column
		for field, _ := range fields {
			columns = append(columns, Column{Name: field})
		}

		_tableMap[table] = Table{
			Name:    table,
			Columns: columns,
		}
	}

	tableMap = _tableMap
}
