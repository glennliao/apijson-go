package main

import (
	"context"

	"github.com/glennliao/apijson-go/drivers/goframe"

	"github.com/glennliao/apijson-go"
	"github.com/glennliao/apijson-go/demo/common/db"
	"github.com/glennliao/apijson-go/model"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"  // need import for mysql
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2" // need import for sqlite
	"github.com/gogf/gf/v2/frame/g"
)

var a *apijson.ApiJson

func init() {
	ctx := context.Background()
	db.InitTable(ctx, g.DB())

	initApiJson(ctx)
}

func initApiJson(ctx context.Context) {
	// 启动goFrame driver
	goframe.Enable()

	useAccessFromCode()

	a = apijson.New()

	a.Config().AccessListProvider = "custom"
	a.Config().RequestListProvider = "custom"

	err := a.Load()
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
}

func main() {
	ctx := context.Background()

	query := a.NewQuery(ctx, model.Map{
		"User[]": model.Map{},
	})

	ret, err := query.Result()
	g.Dump(ret, err)
}
