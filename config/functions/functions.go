package functions

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

type Func struct {
	Handler func(ctx context.Context, param g.Map) (res any, err error)
}

var funcMap = make(map[string]Func)

func Reg(name string, f Func) {
	if _, exists := funcMap[name]; exists {
		panic(fmt.Errorf(" function %s has exists", name))
	}
	funcMap[name] = f
}

func Call(ctx context.Context, name string, param g.Map) (any, error) {
	return funcMap[name].Handler(ctx, param)
}

// functions 提供的功能
// 1. 增加响应字段 -> 该字段需要与系统中别的数据结合处理,如果只是静态处理(去空格,与常量拼接等可直接前端处理即可) 目前会不受_access_ext 中field_get控制, 需处理. 响应字段修改(脱敏、加密、字典转换) 不提供前端控制, 由_access_ext处理
// 2. 通过func节点获取一些系统信息
// 3. actions 中 自定义校验参数、自定义校验权限, 请求体修改(批量字段替换处理?)
// 4. 其他需要自定义的地方 (在action中可看成是hook的替代)

// functions 可用于 field_get 使用, 用于修改请求、响应
