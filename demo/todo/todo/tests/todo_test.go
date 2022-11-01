package tests

import (
	"github.com/glennliao/apijson-go/handlers"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

// TestTodoList 列表查询
func TestTodoList(t *testing.T) {
	req := `
	{
		"Todo":{
		
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

// TestListTodoWithPage 分页列表查询
func TestListTodoWithPage(t *testing.T) {
	req := `
	{
    	"[]": {
			"Todo": {
			  "@column": "id:todoId;title"
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

// TestTodoWithUser 两表关联查询
func TestTodoWithUser(t *testing.T) {
	req := `
	{
    	"Todo": {
			
  		},
		"User":{
			"user_id@":"Todo/user_id"
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

// TestTodoListWithUser 两表关联查询
func TestTodoListWithUser(t *testing.T) {
	req := `
	{
	  "[]": {
		"Todo": {
	
		},
		"User": {
		  "user_id@": "[]/Todo/user_id"
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

// TestTodoListByUser 两表关联查询
func TestTodoListByUser(t *testing.T) {
	req := `
	{
	  "User": {},
	  "[]": {
		"Todo": {
		  "user_id@": "User/user_id"
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

// TestTodoRef 两表关联查询
func TestTodoRef(t *testing.T) {
	req := `
	{
	  "Todo": {
		"@column": "id,user_id",
		"user_id@":"User/userId"
	  },
	  "User": {
		"@column": "user_id"
	  },
	  "[]": {
		"Todo": {
		  "user_id@": "Todo/userId",
		  "@column": "id,user_id"
		},
		"User": {
		  "user_id@": "/Todo/userId",
		  "@column": "user_id"
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

// TestTodoOneMany 列表中一对多
func TestTodoOneMany(t *testing.T) {
	req := `
{
	"[]":{
		"User":{

		},
		"Todo[]":{
			"userId@":"/User/userId"
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
