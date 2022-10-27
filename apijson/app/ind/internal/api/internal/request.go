package internal

import (
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	RequestListReq struct {
		g.Meta `method:"GET" path:"request" tags:"请求参数校验配置" summary:"请求参数校验配置"`
		model.RequestListIn
	}

	RequestListRes struct {
		List []model.RequestListResult `json:"list"`
		Total
	}
)

type (
	RequestGetReq struct {
		g.Meta `method:"GET" path:"request/get" tags:"请求参数校验配置" summary:"获取单条请求参数校验配置"`
		model.RequestGetIn
	}

	RequestGetRes struct {
		*model.RequestGetOut
	}
)

type (
	RequestAddReq struct {
		g.Meta `method:"POST" path:"request" tags:"请求参数校验配置" summary:"新增请求参数校验配置"`
		model.RequestAddIn
	}

	RequestAddRes struct {
		model.RequestAddOut
	}
)

type RequestUpdateReq struct {
	g.Meta `method:"PUT" path:"request" tags:"请求参数校验配置" summary:"修改请求参数校验配置"`
	model.RequestUpdateIn
}

type RequestDeleteReq struct {
	g.Meta `method:"DELETE" path:"request" tags:"请求参数校验配置" summary:"删除请求参数校验配置"`
	Id     string
}

// #26899633f8b49f465ae678e36137a1f4065990a6:021026:52

type (
	RequestMethodListReq struct {
		g.Meta `method:"GET" path:"request/method" tags:"请求参数校验配置" summary:"请求方法列表"`
	}

	RequestMethodListRes struct {
		List []string `json:"list"`
	}
)
