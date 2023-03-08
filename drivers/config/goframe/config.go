package goframe

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

func init() {
	config.RegAccessListProvider("db", func(ctx context.Context) []config.AccessConfig {
		// access
		var accessList []config.AccessConfig

		db := g.DB()

		err := db.Model("_access").Scan(&accessList)
		if err != nil {
			panic(err)
		}

		for _, access := range accessList {
			if len(access.Get) == 1 {
				access.Get = strings.Split(access.Get[0], ",")
			}
			if len(access.Head) == 1 {
				access.Head = strings.Split(access.Head[0], ",")
			}
			if len(access.Gets) == 1 {
				access.Gets = strings.Split(access.Gets[0], ",")
			}
			if len(access.Heads) == 1 {
				access.Heads = strings.Split(access.Heads[0], ",")
			}
			if len(access.Post) == 1 {
				access.Post = strings.Split(access.Post[0], ",")
			}
			if len(access.Put) == 1 {
				access.Put = strings.Split(access.Put[0], ",")
			}
			if len(access.Delete) == 1 {
				access.Delete = strings.Split(access.Delete[0], ",")
			}
			accessList = append(accessList, access)
		}

		return accessList
	})

	config.RegRequestListProvider("db", func(ctx context.Context) []config.Request {
		var requestList []config.Request
		err := g.DB().Model("_request").OrderAsc("version").Scan(&requestList)
		if err != nil {
			panic(err)
		}
		return requestList
	})

	config.RegDbMetaProvider("db", func(ctx context.Context) []config.Table {
		var _tables []config.Table

		db := g.DB()
		tables, err := db.Tables(ctx)
		if err != nil {
			panic(err)
		}

		for _, table := range tables {
			fields, err := db.TableFields(ctx, table)
			if err != nil {
				panic(err)
			}

			var columns []config.Column
			for field, _ := range fields {
				columns = append(columns, config.Column{Name: field})
			}

			_tables = append(_tables, config.Table{
				Name:    table,
				Columns: columns,
			})

		}

		return _tables
	})
}
