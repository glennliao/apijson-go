package main

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	_ "github.com/glennliao/apijson-go/drivers/executor/goframe" // need import for executor with goframe
	"github.com/glennliao/apijson-go/drivers/framework_goframe"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/query"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2" // need import for sqlite
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"log"
	"strings"
	"testing"
)

var a *apijson.ApiJson

func init() {
	a = apijson.Load(App, func(ctx context.Context, a *apijson.ApiJson) {
		// access
		var accessList []config.AccessConfig

		db := g.DB()

		err := db.Model("_access").Scan(&accessList)
		if err != nil {
			panic(err)
		}

		for _, access := range accessList {
			name := access.Alias
			if name == "" {
				name = access.Name
			}

			if len(access.Get) > 0 {
				access.Get = strings.Split(access.Get[0], ",")
			}
			if len(access.Head) > 0 {
				access.Head = strings.Split(access.Head[0], ",")
			}
			if len(access.Gets) > 0 {
				access.Gets = strings.Split(access.Gets[0], ",")
			}
			if len(access.Heads) > 0 {
				access.Heads = strings.Split(access.Heads[0], ",")
			}
			if len(access.Post) > 0 {
				access.Post = strings.Split(access.Post[0], ",")
			}
			if len(access.Put) > 0 {
				access.Put = strings.Split(access.Put[0], ",")
			}
			if len(access.Delete) > 0 {
				access.Delete = strings.Split(access.Delete[0], ",")
			}

			accessList = append(accessList, access)
		}

		a.Config().AccessList = accessList
		a.Config().AccessList = []config.AccessConfig{
			{
				Name:   "user",
				Alias:  "User",
				Get:    []string{"UNKNOWN"},
				RowKey: "id",
			},
			{
				Name:   "user",
				Alias:  "User2",
				Get:    []string{"UNKNOWN"},
				RowKey: "id",
			},
		}

		// request
		var requestList []config.Request
		err = g.DB().Model("_request").OrderAsc("version").Scan(&requestList)
		if err != nil {
			panic(err)
		}
		a.Config().RequestConfig = config.NewRequestConfig(requestList)

		// db meta

		var _tables []config.Table

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

		a.Config().DbMeta = config.NewDbMeta(_tables)

	})
}

// notice: import section
func TestServer(t *testing.T) {
	s := framework_goframe.New(a)
	s.Run()
	// then test in test.http
}

func TestQuery(t *testing.T) {

	ctx := gctx.New()
	q := query.New(ctx, model.Map{
		"User": model.Map{
			//"id":      "123",
			//"id{}":    []string{"123", "456"},
			//"id>":     "222",
			//"@column": "id",
		},
		//"User[]": model.Map{
		//	"@column": "id",
		//	//"userId": "123",
		//},
		//"user2": model.Map{},
		"a@": "User/username",
		"b": model.Map{
			"User": model.Map{
				"id": 1,
			},
			"c@": "/User/username",
		},
		"say()": "test()",
	})

	q.NoAccessVerify = true
	q.Access = a.Config().Access
	q.Access.NoVerify = true
	q.Config = a.Config()
	q.Functions = a.Config().Functions

	result, err := q.Result()

	if err != nil {
		log.Fatalf("%+v", err)
	}

	g.Dump(result)

}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := gctx.New()
		q := query.New(ctx, model.Map{
			"User": model.Map{
				//"id":      "123",
				//"id{}":    []string{"123", "456"},
				//"id>":     "222",
				//"@column": "id",
			},
			"User[]": model.Map{
				"@column": "id",
				//"userId": "123",
			},
			"user2": model.Map{},
			"a@":    "User/username",
			"b": model.Map{
				"User": model.Map{
					"id": 1,
				},
				"c@": "/User/username",
			},
			"say()": "test()",
		})

		q.NoAccessVerify = true
		q.Access = a.Config().Access
		q.Access.NoVerify = true
		q.Config = a.Config()
		q.Functions = a.Config().Functions

		_, err := q.Result()

		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
