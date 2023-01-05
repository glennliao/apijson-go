package app

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	// 设置 rowKey 生成策略
	config.RowKeyGenFunc("time", func(ctx context.Context, genParam g.Map, table string, data g.Map) (g.Map, error) {
		return g.Map{
			"rowKey": gtime.Now().Layout("20060102150405"),
		}, nil
	})
}
