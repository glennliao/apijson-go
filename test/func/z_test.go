package _func

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"reflect"
	"testing"
)

// 1
type Func1 struct {
	ParamList     []string
	ParamTypeList []string
	Handler       func(ctx context.Context, param model.Map) (res interface{}, err error)
}

type Func2 struct {
	Handler any
}

func basic(ctx context.Context, user string) (string, error) {
	return user + " hi", nil
}

var f1 Func1
var f2 Func2
var ctx = gctx.New()

var funcValue reflect.Value

func init() {
	f1 = Func1{
		ParamList:     []string{"req", "xx"},
		ParamTypeList: []string{"string", "string"}, // or any
		Handler: func(ctx context.Context, param model.Map) (res interface{}, err error) {
			return param["user"].(string) + " :hi", nil
		},
	}

	f2 = Func2{Handler: func(ctx context.Context, user string) (string, error) {
		return user + " hi", nil
	}}

	// 将函数包装为反射值对象
	funcValue = reflect.ValueOf(f2.Handler)
}

func TestName(t *testing.T) {
	ret, err := f1.Handler(ctx, model.Map{"user": 1})
	if err != nil {
		panic(err)
	}
	g.Dump(ret)
}

func TestName2(t *testing.T) {
	// 将函数包装为反射值对象
	funcValue := reflect.ValueOf(f2.Handler)
	// 构造函数参数, 传入两个整型值
	paramList := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf("luc")}
	// 反射调用函数
	retList := funcValue.Call(paramList)
	// 获取第一个返回值, 取整数值
	fmt.Println(retList[0].String())
}

func BenchmarkBasic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := basic(ctx, "user1")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := f1.Handler(ctx, model.Map{"user": "user1"})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkName23(b *testing.B) {

	for i := 0; i < b.N; i++ {
		// 构造函数参数, 传入两个整型值
		paramList := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf("luc")}
		// 反射调用函数
		_ = funcValue.Call(paramList)
		// 获取第一个返回值, 取整数值
		//fmt.Println(retList[0].String())
	}
}

func BenchmarkName2(b *testing.B) {

	for i := 0; i < b.N; i++ {
		funcValue := reflect.ValueOf(f2.Handler)
		// 构造函数参数, 传入两个整型值
		paramList := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf("luc")}
		// 反射调用函数
		_ = funcValue.Call(paramList)
		// 获取第一个返回值, 取整数值
		//fmt.Println(retList[0].String())
	}
}
