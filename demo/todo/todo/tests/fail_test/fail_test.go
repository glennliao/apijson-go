package fail_test

import (
	"github.com/glennliao/apijson-go/handlers"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestRefCircle 循环依赖 死锁 - 依赖的节点依赖自己
func TestRefCircle(t *testing.T) {
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
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)

	a := assert.New(t)
	a.NotNil(err)

	if err != nil {
		//panic(err)
	}
	g.Dump(out)
}

// TestRefCircle 循环依赖 死锁 - 依赖的节点最终构成一个圈
func TestRefCircle2(t *testing.T) {
	req := `
{
	"User":{
		"user_id@":"Todo/user_id"
	},
	"Todo":{
		"user_id@":"Notice/user_id"
	},
	"Notice":{
		"id@":"User/id"
	}
}
`
	reqMap := gjson.New(req).Map()
	out, err := handlers.Get(ctx, reqMap)

	a := assert.New(t)
	a.NotNil(err)

	if err != nil {
		panic(err)
	}
	g.Dump(out)
}
