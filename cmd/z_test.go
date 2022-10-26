package main

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/db"
	"github.com/glennliao/apijson-go/handlers"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"testing"
)

func init() {
	db.Init()
}

func TestList(t *testing.T) {
	req := `
{
"User":{

   }


}
`
	db.Init()
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// 此处形成了依赖循坏
func TestRefRef(t *testing.T) {
	req := `
{
	"User":{
	"user_id@":"Todo/user_id"
	   },
	"Todo":{
		"user_id@":"User/user_id"
	}
}
`
	db.Init()
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

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
	config.AccessVerify = false

	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
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
		   "user_id@":"/Todo/user_id"
	   },
		"Todo":{
		
		}
	},
	"Todo[]":{}
	
	}
`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

func TestTotal(t *testing.T) {
	req := `
	{
	 "[]": {
	   "Todo": {
	
	     "@column": "id,user_id"
	   },
	   "User": {
	     "user_id@": "[]/Todo/user_id",
	     "@column": "id"
	   }
	 },
			"total@":"[]/total"
	}
	`

	config.AccessVerify = false

	ctx := context.TODO()

	out, err := handlers.Get(ctx, gjson.New(req).Map())
	if err != nil {
		panic(err)
	}

	g.Dump(out)

	var reqListMap gmap.ListMap
	var resListMap gmap.ListMap

	err = gjson.New(req).Scan(&reqListMap)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	for _, k := range reqListMap.Keys() {
		toKey := k.(string)
		if strings.HasSuffix(toKey, "@") {
			toKey = toKey[0 : len(toKey)-1]
		}
		resListMap.Set(toKey, out[k.(string)])
	}

	g.Dump(resListMap.String())

}

// TestRefOnRef
func TestRefOnRef(t *testing.T) {
	req := `
{

 "User": {
   "id@": "Todo/userId",
   "@column": "id"
 },
 "[]": {
   "Todo": {
     "userId@": "Todo/userId",
     "@column": "id,userId"
   },
   "User": {
     "id@": "/Todo/userId",
     "@column": "id"
   },
   "Credential": {
     "id@": "/User/id",
     "@column": "id"
   }
 },
"Todo": {
   "id": 1627794043692,
   "@column": "id,userId"
 }
}

`
	ctx := context.TODO()
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
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

	out, err := handlers.Post(ctx, reqMap)
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
	ctx = context.WithValue(ctx, config.RoleKey, []string{consts.LOGIN, consts.OWNER})

	//config.AccessCondition = accessCondition
	config.AccessVerify = true

	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

//func accessCondition(ctx context.Context, table string, req g.Map, reqRole string, needRole []string) (g.Map, error) {
//
//	userRole := ctx.Value(config.RoleKey).([]string)
//
//	// 可改成switch方式
//
//	if lo.Contains(needRole, consts.UNKNOWN) {
//		return nil, nil
//	}
//
//	if lo.Contains(needRole, consts.LOGIN) && lo.Contains(userRole, consts.LOGIN) { // 登录后公开资源
//		return nil, nil
//	}
//
//	if lo.Contains(needRole, consts.OWNER) && lo.Contains(userRole, consts.OWNER) {
//		if table == "User" {
//			return g.Map{
//				"id": ctx.Value("ajg.userId"),
//			}, nil
//		} else {
//			return g.Map{
//				"userId": ctx.Value("ajg.userId"),
//			}, nil
//		}
//	}
//
//	return nil, nil
//}
