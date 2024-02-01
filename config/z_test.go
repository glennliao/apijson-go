package config

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
)

type req struct{}

func (r *req) XxXx() {
}

type Req struct {
	Name   string
	UserId string
}

func (r *Req) XxXx() {
}

type Res struct {
	Username string
}

func TestCall(t *testing.T) {
	f := functions{
		funcMap: map[string]*Func{},
	}

	var fn any

	fn = func(ctx context.Context, req *Req) (res any, err error) {
		res = &Res{
			Username: req.Name + "_" + req.UserId,
		}

		return
	}

	f.BindFunc("test", fn)

	ret, err := f.Call(context.TODO(), "test", g.Map{
		"userId": "123",
		"name":   "123",
		"age":    123,
		"sex":    "123",
		"sex2":   "123",
		"sex3":   "",
	})

	ret, err = f.Call(context.TODO(), "test", g.Map{
		"userId": "12311111",
		"name":   "123",
		"age":    123,
		"sex":    "123",
		"sex2":   "123",
		"sex3":   "",
	})

	g.Dump(ret, err)
}

func BenchmarkName(b *testing.B) {
	f := functions{
		funcMap: map[string]*Func{},
	}

	f.BindFunc("test", func(ctx context.Context, req *Req) (res *Res, err error) {
		res = &Res{
			Username: req.Name + "_" + req.UserId,
		}

		return
	})

	for i := 0; i < b.N; i++ {
		f.Call(context.TODO(), "test", g.Map{
			"userId": "123",
			"name":   "123",
			"age":    123,
			"sex":    "123",
			"sex2":   "123",
			"sex3":   "",
		})
	}
}
