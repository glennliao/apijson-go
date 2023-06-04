package goframe

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/executor"
	gfConfig "github.com/glennliao/apijson-go/drivers/goframe/config"
	gfExecutor "github.com/glennliao/apijson-go/drivers/goframe/executor"
)

func init() {

	config.RegAccessListProvider(gfConfig.ProviderName, gfConfig.AccessListDBProvider)
	config.RegRequestListProvider(gfConfig.ProviderName, gfConfig.RequestListProvider)

	config.RegDbMetaProvider(gfConfig.ProviderName, gfConfig.DbMetaProvider)

	executor.RegQueryExecutor("default", gfExecutor.New)
	executor.RegActionExecutor("default", &gfExecutor.ActionExecutor{})

	// todo 添加未配置driver时的报错信息, 而不是 invalid memory address or nil pointer dereference
}
