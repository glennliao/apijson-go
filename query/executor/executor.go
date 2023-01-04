package executor

import (
	"github.com/gogf/gf/v2/frame/g"
)

type QueryExecutor interface {
	ParseCondition(conditions g.MapStrAny, accessVerify bool) error
	ParseCtrl(ctrl g.Map) error
	List(page int, count int, needTotal bool) (list []g.Map, total int64, err error)
	One() (g.Map, error)
	EmptyResult()
}
