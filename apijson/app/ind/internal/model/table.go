package model

import (
	"database/sql"
)

type (
	Column struct {
		Field         string
		Type          string
		Comment       string
		Default       string
		NotNull       bool
		AutoIncrement bool // 仅整型可用
		Index
	}

	Index struct {
		IndexType string // 主键, 唯一, 普通
		IndexName string // 索引名, 相同索引名代表复合索引
		//Index         string // 索引结构, 默认 B+树
	}

	TableAddIn struct {
		TableName string
		Columns   []Column
		Comment   string
	}

	TableAddOut struct {
		Result sql.Result
	}
)

type (
	TableExistIn struct {
		TableName string `v:"required"`
	}

	TableExistOut struct{}
)

type (
	TableSyncIn struct {
		TableExistIn
		TableAlias string
		Tag        string
	}

	TableSyncOut struct {
		//
	}
)
