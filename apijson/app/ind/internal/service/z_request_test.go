package service_test

import (
	"context"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/service"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestRequestService_Add(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	addIn := model.RequestAddIn{
		Debug:     Pointer(int8(0)),
		Version:   Pointer(int8(1)),
		Method:    Pointer(consts.MethodPost),
		Tag:       Pointer("test_" + grand.S(4)),
		Structure: gjson.New(g.Map{}),
		Detail:    nil,
	}

	_, err := service.Request().Add(ctx, addIn)
	a.Nil(err)

	_, err = service.Request().Add(ctx, addIn)
	a.NotNil(err)
}

func TestRequestService_List(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	list, total, err := service.Request().List(ctx, model.RequestListIn{
		ReqList: model.ReqList{
			PageNum:   0,
			PageSize:  0,
			CreatedAt: nil,
		},
	})
	a.Nil(err)

	a.True(len(list) == total)

	row := list[rand.Intn(total/2)]

	row.Version++

	ret, err := service.Request().Update(ctx, model.RequestUpdateIn{
		Id:      row.Id,
		Version: &row.Version,
	})
	a.Nil(err)

	affected, _ := ret.RowsAffected()
	a.True(affected == 1)

}

func TestRequestService_Get(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	row, err := service.Request().Get(ctx, model.RequestGetIn{
		Id: "5",
	})

	a.Nil(err)

	a.True(row != nil)
}

func TestRequestService_Delete(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	addIn := model.RequestAddIn{
		Debug:     Pointer(int8(0)),
		Version:   Pointer(int8(1)),
		Method:    Pointer(consts.MethodPost),
		Tag:       Pointer("test_" + grand.S(4)),
		Structure: gjson.New(g.Map{}),
		Detail:    nil,
	}

	out, err := service.Request().Add(ctx, addIn)
	a.Nil(err)

	row, err := service.Request().Get(ctx, model.RequestGetIn{Id: out.Id})
	a.Nil(err)

	a.True(row != nil)
	g.Dump(row)

	ret, err := service.Request().Delete(ctx, []string{out.Id})
	a.Nil(err)

	g.Dump(ret)
}
