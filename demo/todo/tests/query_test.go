package tests

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQueryExample(t *testing.T) {
	Convey("QueryExample", t, func() {
		// ============================================================
		Convey("total", func() {

			cnt, err := g.DB().Model("notice").Count()
			So(err, ShouldBeNil)
			SoMsg("数据库中不能没有数据,不然无法测试", cnt, ShouldBeGreaterThan, 0)

			req := `
				{
					"Notice[]": {},
					"total@":"Notice[]/total"
				}
			`
			out, err := queryByJsonStr(req)

			So(err, ShouldBeNil)
			So(hasKey(out, "total"), ShouldBeTrue)

			total := gconv.Int(out["total"])
			So(total, ShouldEqual, cnt)
		})

	})
}
