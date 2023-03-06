package config

import (
	"context"
	"github.com/gogf/gf/v2/text/gstr"
)

type FieldStyle func(ctx context.Context, table string, column string) string

func CaseCamel(ctx context.Context, table string, column string) string {
	return gstr.CaseCamelLower(column)
}

func CaseCamelUpper(ctx context.Context, table string, column string) string {
	return gstr.CaseCamel(column)
}

func CaseSnake(ctx context.Context, table string, column string) string {
	return gstr.CaseSnake(column)
}

func Ori(ctx context.Context, table string, column string) string {
	return column
}
