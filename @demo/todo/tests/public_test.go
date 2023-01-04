package tests

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPublicExample(t *testing.T) {

	Convey("Notice (公开资源)", t, func() {
		req := `
			{
				"Notice": {}
			}
		`
		cnt, err := g.DB().Model("notice").Count()
		So(err, ShouldBeNil)
		SoMsg("数据库中不能没有数据,不然无法测试", cnt, ShouldBeGreaterThan, 0)

		// ============================================================
		Convey("未登录用户访问", func() {
			ctx := iAmUnKnow()

			out, err := queryByJsonStr(ctx, req)
			So(err, ShouldBeNil)

			notice := gconv.Map(out["Notice"])
			//g.Dump(notice)

			SoMsg("未获取到数据", len(lo.Keys(notice)), ShouldBeGreaterThan, 0)

		})

		// ============================================================
		Convey("登录用户访问", func() {
			ctx := iAmWM()

			out, err := queryByJsonStr(ctx, req)
			So(err, ShouldBeNil)

			notice := gconv.Map(out["Notice"])
			//g.Dump(notice)

			SoMsg("未获取到数据", len(lo.Keys(notice)), ShouldBeGreaterThan, 0)

		})

	})
}
