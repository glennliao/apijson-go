package config

import (
	"context"
	"fmt"
	"reflect"

	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/container/gvar"
)

const (
	ParamTypeInt    = "int"
	ParamTypeString = "string"
	FromRes         = "res" // 只从响应的数据字段中取, 不从用户传递的数据取
)

type ParamItem struct {
	Type    string
	Name    string
	Desc    string
	Default any
	From    string // 指定参数从何处取值
	V       string // 参数校验规则
}

type Func struct {
	Desc string // 描述
	// 参数可直接读取函数参数传递过来的， ''括起来
	ParamList []ParamItem // 参数列表 // fixme 限制参数来源，强制用户传递的无法覆盖内部的，减免权限的重复判断，  参数校验限制 ， v （最大值，最小值，默认值， 自定义校验。 使用gvaild）
	Handler   func(ctx context.Context, param model.FuncParam) (res any, err error)
}

type functions struct {
	funcMap map[string]*Func
}

func (f *functions) Bind(name string, _func Func) {
	if _, exists := f.funcMap[name]; exists {
		panic(fmt.Errorf(" function %s has exists", name))
	}
	f.funcMap[name] = &_func
}

func (f *functions) BindFunc(name string, function interface{}) {
	funcValue := reflect.ValueOf(function)
	// todo check func param and return

	reqType := funcValue.Type().In(1)

	f.Bind(name, Func{
		Handler: func(ctx context.Context, param model.FuncParam) (res any, err error) {
			//if funcInfo.Type().In(1).Kind() == reflect.Ptr {
			//} else {
			//	inputObject = reflect.New(funcInfo.Type().In(1).Elem()).Elem()
			//	err = param.Scan(inputObject.Addr().Interface())
			//}

			inputObject := reflect.New(reqType.Elem())
			err = param.Scan(inputObject.Interface())
			if err != nil {
				return nil, err
			}

			inputValues := []reflect.Value{
				reflect.ValueOf(ctx),
				inputObject,
			}

			results := funcValue.Call(inputValues)

			if _err, ok := results[0].Interface().(error); ok {
				err = _err
			}

			return results[0], err
		},
		Desc: "",
	})
}

func (f *functions) Call(ctx context.Context, name string, param model.Map) (any, error) {
	params := map[string]model.Var{}
	for k, v := range param {
		params[k] = gvar.New(v)
	}

	return f.funcMap[name].Handler(ctx, params)
}

// functions 可能提供的功能
// 1. 增加响应字段 -> 该字段需要与系统中别的数据结合处理,如果只是静态处理(去空格,与常量拼接等可直接前端处理即可) 目前会不受_access_ext 中field_get控制, 需处理. 响应字段修改(脱敏、加密、字典转换) 不提供前端控制, 由_access_ext处理
// 2. 通过func节点获取一些系统信息
// 3. actions 中 自定义校验参数、自定义校验权限, 请求体修改(批量字段替换处理?)
// 4. 其他需要自定义的地方 (在action中可看成是hook的替代)

// functions 可用于 field_get 使用, 用于修改请求、响应
