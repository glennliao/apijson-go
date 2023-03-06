package config

import "github.com/samber/lo"

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

type DBMeta struct {
	tableMap map[string]Table
}

func NewDbMeta(tables []Table) *DBMeta {
	d := &DBMeta{}
	d.tableMap = make(map[string]Table)

	for _, table := range tables {
		d.tableMap[table.Name] = table
	}

	return d
}

func (d *DBMeta) GetTableColumns(tableName string) (columns []string) {
	for _, column := range d.tableMap[tableName].Columns {
		columns = append(columns, column.Name)
	}
	return
}

func (d *DBMeta) GetTableNameList() []string {
	return lo.Keys(d.tableMap)
}
