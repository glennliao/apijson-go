package main

import (
	"context"
	"github.com/glennliao/apijson-go"
	_ "github.com/glennliao/apijson-go/drivers/config/goframe"   // need import for executor with goframe
	_ "github.com/glennliao/apijson-go/drivers/executor/goframe" // need import for executor with goframe
	"github.com/glennliao/apijson-go/drivers/framework_goframe"
	"github.com/glennliao/apijson-go/model"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2" // need import for sqlite
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"log"
	"testing"
)

var a *apijson.ApiJson

func init() {
	a = apijson.Load(App)
}

// notice: import section
func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	s := framework_goframe.New(a)
	s.Run()
	// then test in test.http
}

func TestQuery(t *testing.T) {

	ctx := gctx.New()

	q := a.NewQuery(ctx, model.Map{
		"User": model.Map{
			//"id":      "123",
			//"id{}":    []string{"123", "456"},
			//"id>":     "222",
			//"@column": "id",
		},
		"User[]": model.Map{
			"@column": "id",
			//"userId": "123",
		},
		//"user2": model.Map{},
		"a@": "User/username",
		"b": model.Map{
			"User": model.Map{
				"id": 1,
			},
			"c@": "/User/username",
		},
		"say()": "test()",
	})

	q.NoAccessVerify = true

	result, err := q.Result()

	if err != nil {
		log.Fatalf("%+v", err)
	}

	g.Dump(result)

}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		q := a.NewQuery(ctx, model.Map{
			//"User": model.Map{
			//	"id": 1,
			//	//"id":      "123",
			//	//"id{}":    []string{"123", "456"},
			//	//"id>":     "222",
			//	//"@column": "id",
			//},
			//"User[]": model.Map{
			//	"@column": "id",
			//	//"userId": "123",
			//},
			//"a@": "User/username",
			//"b": model.Map{
			//	"User": model.Map{
			//		"id": 1,
			//	},
			//	//"c@": "/User/username",
			//},
			"say()": "test()",
		})

		q.NoAccessVerify = true

		_, err := q.Result()

		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
