package db

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	Column struct {
		Name string
	}

	Table struct {
		Name    string
		Columns []Column
	}
)

var tableMap = map[string]Table{}

func loadTableMeta() {
	var ctx = context.TODO()

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

		tableMap[table] = Table{
			Name:    table,
			Columns: columns,
		}
	}
}
