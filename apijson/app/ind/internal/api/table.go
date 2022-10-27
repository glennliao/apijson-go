package api

import (
	"context"
	. "github.com/glennliao/apijson-go/apijson/app/ind/internal/api/internal"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/service"
)

type tableApi struct{}

var Table = tableApi{}

func (a *tableApi) Sync(ctx context.Context, req *TableSyncReq) (res *TableSyncRes, err error) {
	out, err := service.Table().Sync(ctx, req.TableSyncIn)
	return &TableSyncRes{TableSyncOut: out}, err
}
