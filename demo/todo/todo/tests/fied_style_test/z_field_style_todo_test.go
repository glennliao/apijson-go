package fied_style_test

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/handlers"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

func init() {
	config.SetDbFieldStyle(config.CaseSnake)

	config.SetJsonFieldStyle(config.CaseSnake)
	//config.SetJsonFieldStyle(consts.CASE_CAMEL_UPPER)
}

// TestCaseCameTodoList 列表查询
func TestCaseCameTodoList(t *testing.T) {
	req := `
	{
		"Todo[]":{
			"query":"1"
		   }
	}
`
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestCaseCameListTodoWithPage 分页列表查询
func TestCaseCameListTodoWithPage(t *testing.T) {
	req := `
	{
    	"[]": {
			"Todo": {
			  "@column": "id:todoId;title;userId:userIdCaseCamelLower"
			},
			"page": 2,
			"count": 2
  		}
	}
`

	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestCaseCameTodoWithUser 两表关联查询
func TestCaseCameTodoWithUser(t *testing.T) {
	req := `
	{
    	"Todo": {
			
  		},
		"User":{
			"userId@":"Todo/userId"
		}
	}
`

	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestCaseCameTodoListWithUser 两表关联查询
func TestCaseCameTodoListWithUser(t *testing.T) {
	config.SetJsonFieldStyle(config.CaseCamel)

	req := `
	{
	  "[]": {
		"Todo": {
	
		},
		"User": {
		  "userId@": "[]/Todo/userId"
		}
	  }
	}
`

	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestCaseCameTodoListByUser 两表关联查询
func TestCaseCameTodoListByUser(t *testing.T) {
	req := `
	{
	  "User": {},
	  "[]": {
		"Todo": {
		  "userId@": "User/userId"
		}
	  }
	}
`

	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestCaseCameTodoRef 两表关联查询
func TestCaseCameTodoRef(t *testing.T) {
	req := `
	{
	  "Todo": {
		"@column": "id,userId",
		"userId@":"User/userId"
	  },
	  "User": {
		"@column": "userId"
	  },
	  "[]": {
		"Todo": {
		  "userId@": "Todo/userId",
		  "@column": "id,userId"
		},
		"User": {
		  "userId@": "/Todo/userId",
		  "@column": "userId"
		}
	  }
	}
`

	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestCaseCameTodoOneMany 列表中一对多
func TestCaseCameTodoOneMany(t *testing.T) {
	req := `
{
	"[]":{
		"User":{
			
		},
		"Todo[]":{
			"userId@":"/User/userId"
		},
		"query": "1"
	}
}
`
	reqMap := gjson.New(req).Map()

	out, err := handlers.Get(ctx, reqMap)
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}
