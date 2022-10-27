package service_test

import (
	"context"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/service"
	"github.com/glennliao/apijson-go/apijson/internal/dao"
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestAccessService_Add(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	list, total, err := service.Access().List(ctx, model.AccessListIn{})
	a.Nil(err)
	a.Len(list, total)

	out, err := service.Access().Add(ctx, model.AccessAddIn{
		Debug: Pointer(int8(0)),
		Name:  Pointer("test_" + grand.S(4)),
		Alias: Pointer("test_" + grand.S(4)),
		Get:   Pointer([]string{"TEST_OWNER", "TEST_ADMIN"}),
	})

	a.Nil(err)
	g.Dump(out)

	list2, total2, err := service.Access().List(ctx, model.AccessListIn{})
	a.Nil(err)

	a.Len(list2, total2)
	a.Equal(total2, total+1)
}

func TestAccessService_Get(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	list, total, err := service.Access().List(ctx, model.AccessListIn{})
	a.Nil(err)
	a.Len(list, total)

	item := list[rand.Intn(total)]

	row, err := service.Access().Get(ctx, item.Id)
	a.Nil(err)
	a.NotNil(row)

	a.Equal(row.Name, item.Name)

	g.Dump(row)
}

func TestAccessService_Update(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	var (
		m, c = dao.Access.MC(ctx)
		list []model.AccessListResult
	)

	err := m.WhereLike(c.Name, "test_%").Scan(&list)
	a.Nil(err)

	rand.Seed(time.Now().UnixNano())

	item := list[rand.Intn(len(list))]

	updateIn := model.AccessUpdateIn{Id: item.Id}
	ret, err := service.Access().Update(ctx, updateIn)
	a.Nil(err)

	affected, err := ret.RowsAffected()
	a.Nil(err)
	a.True(affected == 0)

	roleList := Pointer(config.RoleList())
	rand.Shuffle(len(*roleList), func(i, j int) { (*roleList)[i], (*roleList)[j] = (*roleList)[j], (*roleList)[i] })

	switch rand.Intn(8) {
	case 0:
		updateIn.Get = roleList
	case 1:
		updateIn.Head = roleList
	case 2:
		updateIn.Get = roleList
	case 3:
		updateIn.Gets = roleList
	case 4:
		updateIn.Heads = roleList
	case 5:
		updateIn.Post = roleList
	case 6:
		updateIn.Put = roleList
	case 7:
		updateIn.Delete = roleList
	}

	ret, err = service.Access().Update(ctx, updateIn)
	a.Nil(err)

	affected, err = ret.RowsAffected()
	a.Nil(err)
	a.True(affected == 1)
}

func TestAccessService_Delete(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	out, err := service.Access().Add(ctx, model.AccessAddIn{
		Debug: Pointer(int8(0)),
		Name:  Pointer("test_" + grand.S(4)),
		Alias: Pointer("test_" + grand.S(4)),
		Get:   Pointer([]string{"TEST_OWNER", "TEST_ADMIN"}),
	})
	a.Nil(err)

	row, err := service.Access().Get(ctx, out.Id)
	a.Nil(err)
	a.NotNil(row)

	g.Dump(row)

	_, err = service.Access().Delete(ctx, []string{out.Id})
	a.Nil(err)
}
