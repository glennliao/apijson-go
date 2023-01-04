package tests

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQueryExample(t *testing.T) {
	Convey("QueryExample", t, func() {
		// ============================================================
		Convey("total", func() {
			ctx := gctx.New()

			cnt, err := g.DB().Model("notice").Ctx(ctx).Count()
			So(err, ShouldBeNil)
			SoMsg("数据库中不能没有数据,不然无法测试", cnt, ShouldBeGreaterThan, 0)

			req := `
				{
					"Notice[]": {},
					"total@":"Notice[]/total"
				}
			`
			out, err := queryByJsonStr(ctx, req)

			So(err, ShouldBeNil)
			So(hasKey(out, "total"), ShouldBeTrue)

			total := gconv.Int(out["total"])
			So(total, ShouldEqual, cnt)
		})

	})
}

func TestTodoGet(t *testing.T) {
	Convey("TodoGet", t, func() {
		ctx := iAmWM()
		userId := UserIdWM

		req := fmt.Sprintf(`
				{
					"Todo":{
							"userId": "%s"
					}
				}
		`, userId)

		out, err := queryByJsonStr(ctx, req)
		So(err, ShouldBeNil)
		So(hasKey(out, "Todo"), ShouldBeTrue)

		data := out["Todo"].(g.Map)
		So(data["userId"], ShouldEqual, userId)
	})
}

func TestTodoList(t *testing.T) {
	Convey("TodoList", t, func() {
		ctx := iAmWM()

		userId := UserIdWM

		cnt, err := g.DB().Model("t_todo").Where("user_id", userId).Count()
		So(err, ShouldBeNil)

		req := fmt.Sprintf(`
			{
				"[]": {
					"Todo":{
						"userId": "%s"
					}
				}
			}
		`, userId)

		out, err := queryByJsonStr(ctx, req)
		So(err, ShouldBeNil)
		data := out["[]"]
		So(data, ShouldNotBeNil)

		list := data.([]g.Map)
		// 默认分页查询: limit 0,10
		if cnt > 10 {
			cnt = 10
		}
		So(len(list), ShouldEqual, cnt)

		for _, item := range list {
			elem := item["Todo"].(g.Map)
			So(elem["userId"], ShouldEqual, userId)
		}
	})
}

// TestListTodoWithPage 分页列表查询
func TestListTodoWithPage(t *testing.T) {
	Convey("ListTodoWithPage", t, func() {

		req := `
		{
			"[]": {
				"Todo": {
				  "@column": "id:todoId;title"
				},
				"page": 2,
				"count": 2
			}
		}
`
		out, err := queryByJsonStr(iAmUnKnow(), req)
		So(err, ShouldBeNil)

		list := out["[]"].([]g.Map)
		So(len(list), ShouldBeGreaterThan, 0)

		for _, item := range list {
			todo := item["Todo"].(g.Map)

			So(len(todo), ShouldEqual, 2)
			So(hasKey(todo, "todoId"), ShouldBeTrue)
			So(hasKey(todo, "title"), ShouldBeTrue)
		}
	})
}
