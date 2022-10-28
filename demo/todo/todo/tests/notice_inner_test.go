package tests

import (
	"github.com/glennliao/apijson-go/handlers"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

// TestNotice 查看公共公告
func TestNoticeInner(t *testing.T) {
	req := `
	{
		"NoticeInner": {}
	}
`
	out, err := handlers.Get(ctx, gjson.New(req).Map())
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}

// TestNotice 查看公共公告
func TestNoticeInnerList(t *testing.T) {
	req := `
	{
		"NoticeInner[]": {},
		"total@":"NoticeInner[]/total"
	}
`
	out, err := handlers.Get(ctx, gjson.New(req).Map())
	if err != nil {
		panic(err)
	}
	g.Dump(out)
}
