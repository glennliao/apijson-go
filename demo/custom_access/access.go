package main

import (
	"context"

	"github.com/glennliao/apijson-go/config"
)

func useAccessFromCode() {
	config.RegAccessListProvider("custom", func(ctx context.Context) []config.AccessConfig {
		return []config.AccessConfig{
			{
				Name:   "user", // 数据库表名
				Alias:  "User", // 实际访问的名称
				Get:    []string{"UNKNOWN"},
				RowKey: "id", // 主键
				FieldsGet: map[string]*config.FieldsGetValue{
					"default": {
						In: nil,
						// 配置可查看到的字段, key 为数据库内字段名, value 暂保留为空串
						Out: map[string]string{
							"id":       "",
							"username": "",
						},
					},
				},
			},
			{
				Name:   "todo", // 数据库表名
				Alias:  "Todo", // 实际访问的名称
				Get:    []string{"UNKNOWN"},
				RowKey: "id", // 主键
				FieldsGet: map[string]*config.FieldsGetValue{
					"default": {
						In: nil,
						// 配置可查看到的字段, key 为数据库内字段名, value 暂保留为空串
						Out: map[string]string{
							"content": "",
						},
					},
				},
			},
		}
	})

	config.RegRequestListProvider("custom", func(ctx context.Context) []config.RequestConfig {
		// 此处还未使用到请求, 故返回空
		return []config.RequestConfig{}
	})
}
