package tests

import (
	"github.com/glennliao/apijson-go/config"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoginExample(t *testing.T) {

	oriAccessVerify := config.AccessVerify

	Convey("NoticeInner (登录可见资源)", t, func() {

		config.AccessVerify = true

		req := `
			{
				"NoticeInner": {},
"User":{}
			}
		`
		cnt, err := g.DB().Model("notice_inner").Count()
		So(err, ShouldBeNil)
		SoMsg("数据库中不能没有数据,不然无法测试", cnt, ShouldBeGreaterThan, 0)

		// ============================================================
		Convey("未登录用户访问", func() {
			ctx := iAmUnKnow()

			out, err := queryByJsonStr(ctx, req)
			So(err, ShouldNotBeNil)

			notice := gconv.Map(out["NoticeInner"])
			//g.Dump(notice)

			SoMsg("未获取到数据", len(lo.Keys(notice)), ShouldEqual, 0)

		})

		// ============================================================
		Convey("登录用户访问", func() {
			ctx := iAmWM()

			out, err := queryByJsonStr(ctx, req)
			So(err, ShouldBeNil)

			notice := gconv.Map(out["NoticeInner"])
			//g.Dump(notice)

			SoMsg("未获取到数据", len(lo.Keys(notice)), ShouldBeGreaterThan, 0)

		})

	})

	config.AccessVerify = oriAccessVerify

}
