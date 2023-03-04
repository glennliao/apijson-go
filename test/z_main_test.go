package main

import (
	"github.com/glennliao/apijson-go"
	_ "github.com/glennliao/apijson-go/drivers/executor/goframe" // need import for executor with goframe
	"github.com/glennliao/apijson-go/drivers/framework_goframe"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/query"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2" // need import for sqlite
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

var a *apijson.ApiJson

func init() {
	a = apijson.Load(App)
}

// notice: import section
func TestServer(t *testing.T) {
	s := framework_goframe.New(a)
	s.Run()
	// then test in test.http
}

func TestQuery(t *testing.T) {

	ctx := gctx.New()
	q := query.New(ctx, model.Map{
		"t_user": model.Map{
			"id":      "123",
			"id{}":    []string{"123", "456"},
			"id>":     "222",
			"@column": "id,userId",
		},
		"t_user[]": model.Map{
			//"userId": "123",
		},
		//"t_todo":  model.Map{},
		//"_access": model.Map{},
	})

	q.NoAccessVerify = true //config.AccessVerify
	q.Access = a.Config().Access

	result, err := q.Result()

	if err != nil {
		panic(err)
	}

	g.Dump(result)

}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := gctx.New()

		q := query.New(ctx, model.Map{
			"Todo": model.Map{},
			"User": model.Map{},
		})

		_, err := q.Result()
		if err != nil {
			panic(err)
		}
	}
}
