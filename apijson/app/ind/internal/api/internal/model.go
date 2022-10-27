package internal

import (
	"database/sql"
)

type RowsRes struct {
	Result sql.Result
}

type Total struct {
	Total int `json:"total" des:"总数"`
}
