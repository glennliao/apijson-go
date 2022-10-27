package api

import (
	"context"
	. "github.com/glennliao/apijson-go/apijson/app/ind/internal/api/internal"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/service"
	"github.com/glennliao/apijson-go/config"
	"strings"
)

type requestApi struct{}

var Request = requestApi{}

// List 列表
func (a *requestApi) List(ctx context.Context, req *RequestListReq) (res *RequestListRes, err error) {
	list, total, err := service.Request().List(ctx, req.RequestListIn)
	return &RequestListRes{
		List:  list,
		Total: Total{Total: total},
	}, err
}

// Get 获取单条记录
func (a *requestApi) Get(ctx context.Context, req *RequestGetReq) (res *RequestGetRes, err error) {
	out, err := service.Request().Get(ctx, req.RequestGetIn)
	return &RequestGetRes{
		RequestGetOut: out,
	}, err
}

// Add 添加记录
func (a *requestApi) Add(ctx context.Context, req *RequestAddReq) (res *RequestAddRes, err error) {
	*req.Method = strings.ToUpper(*req.Method)

	out, err := service.Request().Add(ctx, req.RequestAddIn)
	return &RequestAddRes{
		RequestAddOut: out,
	}, err
}

// Update 更新记录
func (a *requestApi) Update(ctx context.Context, req *RequestUpdateReq) (res *RowsRes, err error) {
	*req.Method = strings.ToUpper(*req.Method)

	ret, err := service.Request().Update(ctx, req.RequestUpdateIn)
	return &RowsRes{Result: ret}, err
}

// Delete 删除单条记录
func (a *requestApi) Delete(ctx context.Context, req *RequestDeleteReq) (res *RowsRes, err error) {
	ret, err := service.Request().Delete(ctx, []string{req.Id})
	return &RowsRes{Result: ret}, err
}

// #183f2c630cb500fe7ca7864a59324ea200f977d4:021026:52

func (a *requestApi) MethodList(ctx context.Context, req *RequestMethodListReq) (res *RequestMethodListRes, err error) {
	list := config.MethodList()
	return &RequestMethodListRes{List: list}, nil
}
