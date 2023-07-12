package goframe

import (
	"context"

	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
	gfConfig "github.com/glennliao/apijson-go/drivers/goframe/config"
	gfExecutor "github.com/glennliao/apijson-go/drivers/goframe/executor"
	"github.com/glennliao/apijson-go/query"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {

	config.RegAccessListProvider(gfConfig.ProviderName, gfConfig.AccessListDBProvider)
	config.RegRequestListProvider(gfConfig.ProviderName, gfConfig.RequestListProvider)

	config.RegDbMetaProvider(gfConfig.ProviderName, gfConfig.DbMetaProvider)

	query.RegExecutor("default", gfExecutor.New)

	action.RegExecutor("default", &gfExecutor.ActionExecutor{
		DbResolver: func(ctx context.Context) gdb.DB {
			return g.DB()
		},
	})

	action.RegTransactionResolver(func(ctx context.Context, req *action.Action) action.TransactionHandler {
		return func(ctx context.Context, action func(ctx context.Context) error) error {
			return g.DB().Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				return action(ctx)
			})
		}
	})

}
