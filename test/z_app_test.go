package main

import (
	"context"
	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/tables"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/table-sync/tablesync"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type User struct {
	Id        uint32 `ddl:"primaryKey"`
	Username  string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Todo struct {
	Id        uint32 `ddl:"primaryKey"`
	UserId    uint32
	Content   string
	CreatedAt *time.Time
}

func init() {
	config.RegAccessListProvider("db", func(ctx context.Context) []config.AccessConfig {
		return []config.AccessConfig{
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
	})
}

func App(ctx context.Context, a *apijson.ApiJson) {

	syncer := tablesync.Syncer{
		Tables: []tablesync.Table{
			User{}, Todo{},
			tables.Access{}, tables.Request{},
		},
	}
	err := syncer.Sync(ctx, g.DB())
	if err != nil {
		panic(err)
	}

	a.Config().Functions.Bind("test", config.Func{
		Handler: func(ctx context.Context, param model.Map) (res any, err error) {
			return "你好", nil
		},
	})

	//a.Config().AccessListProvider = "custom"

}
