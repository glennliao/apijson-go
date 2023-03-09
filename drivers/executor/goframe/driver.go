package goframe

import (
	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/drivers/executor_goframe"
)

func init() {
	executor.RegQueryExecutor("default", executor_goframe.New)
	executor.RegActionExecutor("default", &executor_goframe.ActionExecutor{})
}
