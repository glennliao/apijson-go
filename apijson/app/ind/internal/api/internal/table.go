package internal

import (
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	TableSyncReq struct {
		g.Meta `method:"POST" path:"table/sync" tags:"表设置" summary:"新增表设置"`
		model.TableSyncIn
	}

	TableSyncRes struct {
		model.TableSyncOut
	}
)
