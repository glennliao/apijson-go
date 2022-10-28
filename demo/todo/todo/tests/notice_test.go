package tests

import (
	"github.com/glennliao/apijson-go/handlers"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

// TestNotice 查看公共公告
func TestNotice(t *testing.T) {
	req := `
	{
		"Notice": {}
	}
`
	out, err := handlers.Get(ctx, gjson.New(req).Map())
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestNotice 查看公共公告
func TestNoticeList(t *testing.T) {
	req := `
	{
		"Notice[]": {},
		"total@":"Notice[]/total"
	}
`
	out, err := handlers.Get(ctx, gjson.New(req).Map())
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}
