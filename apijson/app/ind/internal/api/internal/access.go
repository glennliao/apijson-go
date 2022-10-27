package internal

import (
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	AccessListReq struct {
		g.Meta `method:"GET" path:"access" tags:"权限配置" summary:"权限配置列表"`
		model.AccessListIn
	}

	AccessListRes struct {
		List []model.AccessListResult `json:"list"`
		Total
	}
)

type (
	AccessGetReq struct {
		g.Meta `method:"GET" path:"access/get" tags:"权限配置" summary:"获取单条权限配置"`
		Id     string
	}

	AccessGetRes struct {
		*model.AccessGetOut
	}
)

type (
	AccessAddReq struct {
		g.Meta `method:"POST" path:"access" tags:"权限配置" summary:"新增权限配置"`
		model.AccessAddIn
	}

	AccessAddRes struct {
		model.AccessAddOut
	}
)

type AccessUpdateReq struct {
	g.Meta `method:"PUT" path:"access" tags:"权限配置" summary:"修改权限配置"`
	model.AccessUpdateIn
}

type AccessDeleteReq struct {
	g.Meta `method:"DELETE" path:"access" tags:"权限配置" summary:"删除权限配置"`
	Id     string
}

// #1093608713ee8f272c9eb4b64933fb6f09c6d352:021026:52

type (
	AccessRoleListReq struct {
		g.Meta `method:"GET" path:"access/role" tags:"权限配置" summary:"权限配置角色列表"`
	}

	AccessRoleListRes struct {
		List []string `json:"list"`
	}
)
