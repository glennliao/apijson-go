package config

import (
	"context"
	"github.com/gogf/gf/v2/text/gstr"
)

// =========================  字段配置 =======================

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

// jsonFieldStyle 数据库返回的字段
var jsonFieldStyleFunc FieldStyle = CaseCamel

// dbFieldStyle 数据库字段命名风格 请求传递到数据库中
var dbFieldStyleFunc = CaseSnake

// SetJsonFieldStyle 设置返回的 json字段风格,
func SetJsonFieldStyle(style FieldStyle) {
	jsonFieldStyleFunc = style
}

// SetDbFieldStyle 设置数据库的字段风格
func SetDbFieldStyle(style FieldStyle) {
	dbFieldStyleFunc = style
}

// GetJsonFieldStyle 设置返回的 json字段风格,
func GetJsonFieldStyle() FieldStyle {
	if jsonFieldStyleFunc == nil {
		return Ori
	}
	return jsonFieldStyleFunc
}

// GetDbFieldStyle 设置数据库的字段风格
func GetDbFieldStyle() FieldStyle {
	if dbFieldStyleFunc == nil {
		return Ori
	}
	return dbFieldStyleFunc
}
