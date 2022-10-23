package apijson

import (
	"context"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"my-apijson/apijson/db"
	"my-apijson/apijson/query"
	"testing"
)

func TestTwoTableGet(t *testing.T) {
	req := `
{
 "User":{
        "id@":"Todo/userId"
    },
    "Todo":{
        "id":1627794043692
    }
   
}
`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

func TestTowTableGetList(t *testing.T) {
	req := `
{
 	"[]":{
	"User":{
        "id@":"/Todo/userId"
    },
    "Todo":{
       
    }
	}
   
}
`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

func TestCheckRequest(t *testing.T) {

	db.Init()

	req := `
{
    "Todo":{
        
		"title":"asdasda"
    },
	"tag":"Todo"
}
`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()

	out, err := Post(ctx, reqMap)
	if err != nil {
		g.Dump(err)
	}
	g.Dump(out)
}

func TestAccess(t *testing.T) {

	db.Init()

	req := `
{
 "User":{
        "id@":"Todo/userId"
    },
    "Todo":{
        "id":1
    }
   
}
`
	ctx := context.TODO()

	ctx = context.WithValue(ctx, "ajg.userId", "2")
	ctx = context.WithValue(ctx, "ajg.role", []string{query.LOGIN, query.OWNER})

	reqMap := gjson.New(req).Map()
	out, err := Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}
