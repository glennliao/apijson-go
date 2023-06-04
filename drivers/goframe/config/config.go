package config

import (
	"context"
	"strings"

	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// 设置 _access/_request 自定义表名
var (
	TableAccess  = "_access"
	TableRequest = "_request"
	ProviderName = "db"
)

func RequestListProvider(ctx context.Context) []config.Request {
	var requestList []config.Request
	err := g.DB().Model(TableRequest).OrderAsc("version").Scan(&requestList)
	if err != nil {
		panic(err)
	}

	for i, item := range requestList {
		item := item

		if item.Structure == nil {
			item.Structure = make(map[string]*config.Structure)
		}

		// provider处理
		//if strings.ToLower(tag) != tag {
		//	// 本身大写, 如果没有外层, 则套一层
		//	if _, ok := item.Structure[tag]; !ok {
		//		item.Structure = map[string]any{
		//			tag: item.Structure,
		//		}
		//	}
		//}

		for k, v := range item.Structure {
			structure := config.Structure{}
			err := gconv.Scan(v, &structure)
			if err != nil {
				panic(err)
			}

			if structure.Must != nil {
				structure.Must = strings.Split(structure.Must[0], ",")
			}
			if structure.Refuse != nil {
				structure.Refuse = strings.Split(structure.Refuse[0], ",")
			}

			item.Structure[k] = &structure
		}

		if len(item.ExecQueue) > 0 {
			item.ExecQueue = strings.Split(item.ExecQueue[0], ",")
		}

		requestList[i] = item

	}

	return requestList
}

func DbMetaProvider(ctx context.Context) []config.Table {
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
}

func AccessListDBProvider(ctx context.Context) []config.AccessConfig {
	// access
	var accessList []config.AccessConfig

	db := g.DB()

	err := db.Model(TableAccess).Scan(&accessList)
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
}
