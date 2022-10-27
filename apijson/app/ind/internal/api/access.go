package api

import (
	"context"
	. "github.com/glennliao/apijson-go/apijson/app/ind/internal/api/internal"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/service"
	"github.com/glennliao/apijson-go/config"
)

type accessApi struct{}

var Access = accessApi{}

// List 列表
func (a *accessApi) List(ctx context.Context, req *AccessListReq) (res *AccessListRes, err error) {
	list, total, err := service.Access().List(ctx, req.AccessListIn)
	return &AccessListRes{
		List:  list,
		Total: Total{Total: total},
	}, err
}

// Get 获取单条记录
func (a *accessApi) Get(ctx context.Context, req *AccessGetReq) (res *AccessGetRes, err error) {
	out, err := service.Access().Get(ctx, req.Id)
	return &AccessGetRes{
		AccessGetOut: out,
	}, err
}

// Add 添加记录
func (a *accessApi) Add(ctx context.Context, req *AccessAddReq) (res *AccessAddRes, err error) {
	out, err := service.Access().Add(ctx, req.AccessAddIn)
	return &AccessAddRes{
		AccessAddOut: out,
	}, err
}

// Update 更新记录
func (a *accessApi) Update(ctx context.Context, req *AccessUpdateReq) (res *RowsRes, err error) {
	ret, err := service.Access().Update(ctx, req.AccessUpdateIn)
	return &RowsRes{Result: ret}, err
}

// Delete 删除单条记录
func (a *accessApi) Delete(ctx context.Context, req *AccessDeleteReq) (res *RowsRes, err error) {
	ret, err := service.Access().Delete(ctx, []string{req.Id})
	return &RowsRes{Result: ret}, err
}

// #9fcc3ff4d6c5a84ea54fd4d8cf07836ade28dc59:021026:52

func (a *accessApi) RoleList(ctx context.Context, req *AccessRoleListReq) (res *AccessRoleListRes, err error) {
	list := config.RoleList()
	return &AccessRoleListRes{List: list}, nil
}
