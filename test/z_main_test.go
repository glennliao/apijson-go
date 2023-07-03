package main

import (
	"context"
	"log"
	"testing"

	"github.com/glennliao/apijson-go"
	_ "github.com/glennliao/apijson-go/drivers/goframe"
	"github.com/glennliao/apijson-go/drivers/goframe/web"
	"github.com/glennliao/apijson-go/model"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2" // need import for sqlite
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
	s := web.New(a)
	s.Run()
	// then test in test.http
}

func TestQuery(t *testing.T) {

	ctx := gctx.New()

	q := a.NewQuery(ctx, model.Map{
		"User": model.Map{
			// "id":      "123",
			// "id{}":    []string{"123", "456"},
			// "id>":     "222",
			// "@column": "id",
		},
		"User[]": model.Map{
			"@column":      "id,username",
			"concatTest()": "concat(username,c)",
			// "userId": "123",
		},
		// "user2": model.Map{},
		// "a@": "User/username",
		// "b": model.Map{
		// 	"User": model.Map{
		// 		"id": 1,
		// 	},
		// 	"c@": "/User/username",
		// },
		// "say()":        "test()",
		// "a":            "12",
		// "c":            "34",
		// "concatTest()": "concat(/User/username,c)",
		// "concatTest()": "concat(User/username,c)",
	})

	q.NoAccessVerify = true

	result, err := q.Result()

	if err != nil {
		log.Fatalf("%+v", err)
	}

	g.Dump(result)

}

func TestAlias(t *testing.T) {

	ctx := gctx.New()

	q := a.NewQuery(ctx, model.Map{
		"User[]": model.Map{
			"@column": "id,password:username",
			// "userId": "123",
		},
	})

	q.NoAccessVerify = false

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
			// "User": model.Map{
			//	"id": 1,
			//	//"id":      "123",
			//	//"id{}":    []string{"123", "456"},
			//	//"id>":     "222",
			//	//"@column": "id",
			// },
			// "User[]": model.Map{
			//	"@column": "id",
			//	//"userId": "123",
			// },
			// "a@": "User/username",
			// "b": model.Map{
			//	"User": model.Map{
			//		"id": 1,
			//	},
			//	//"c@": "/User/username",
			// },
			"say()": "test()",
		})

		q.NoAccessVerify = true

		_, err := q.Result()

		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
