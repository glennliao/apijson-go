package app

import (
	"context"
	"github.com/glennliao/apijson-go/action"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	// 设置 rowKey 生成策略
	action.RowKeyGenFunc(func(ctx context.Context, table string, data g.Map) (g.Map, error) {

		if table == "t_todo" {
			t := gtime.Now().Layout("20060102150405")
			return g.Map{
				"todo_id": t,
			}, nil
		}

		return nil, nil

	})
}
