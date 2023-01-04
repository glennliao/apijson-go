package tests

import (
	"fmt"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSome(t *testing.T) {
	iAmWM()
	m := g.DB().Model("t_todo").Ctx(ctx)

	Convey("单表单条数据操作", t, func() {
		Convey("新增", func() {

			cnt1, err := m.Clone().Count(g.Map{
				"user_id": UserIdWM,
			})
			So(err, ShouldBeNil)

			req := `
				{
					"Todo": {
						"title": "去找林云喝茶 ♪(^∇^*)"
					},
					"TodoLog":{
						"log":"created by one"
					},
					"TodoLog[]":[
						{"log":"created by list[0]"},
						{"log":"created by list[1]"}
					],
					"tag": "Todo",
					"version": 2
				}
			`

			out, err := actionByJsonStr(req, consts.MethodPost)
			So(err, ShouldBeNil)

			//g.Dump(out)
			todo := out["Todo"].(g.Map)
			todoId := todo["todoId"].(string)

			cnt2, err := m.Clone().Count(g.Map{
				"user_id": UserIdWM,
			})
			So(err, ShouldBeNil)
			So(cnt2-cnt1, ShouldEqual, 1)

			cnt, err := g.DB().Model("t_todo_log").Ctx(ctx).Count(g.Map{
				"todo_id": todoId,
				"log":     "created by one",
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 1)

			cnt, err = g.DB().Model("t_todo_log").Ctx(ctx).WhereLike("log", "created by list%").Count(g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 2)

		})
	})

}

func TestActionOneTableOneLine(t *testing.T) {
	iAmWM()
	todoId := ""
	Convey("单表单条数据操作", t, func() {

		m := g.DB().Model("t_todo").Ctx(ctx)

		// ===================================================================
		Convey("新增", func() {

			cnt1, err := m.Clone().Count(g.Map{
				"user_id": UserIdWM,
			})
			So(err, ShouldBeNil)

			req := `
				{
					"Todo": {
						"title": "去找林云喝茶"
					},
					"tag": "Todo",
					"version": 1
				}
			`

			out, err := actionByJsonStr(req, consts.MethodPost)
			So(err, ShouldBeNil)

			//g.Dump(out)
			todo := out["Todo"].(g.Map)
			todoId = todo["todoId"].(string)

			cnt2, err := m.Clone().Count(g.Map{
				"user_id": UserIdWM,
			})
			So(err, ShouldBeNil)

			So(cnt2-cnt1, ShouldEqual, 1)

		})

		// ===================================================================
		Convey("修改", func() {

			req := `
				{
					"tag":"Todo",
					"Todo":{
						"todoId":"%s",
						"title":"去找林云喝茶, 把史强的预约先取消"
					}
				}
`

			_, err := actionByJsonStr(fmt.Sprintf(req, todoId), consts.MethodPut)
			So(err, ShouldBeNil)

			one, err := m.Clone().One(g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			So(one.Map()["title"], ShouldEqual, "去找林云喝茶, 把史强的预约先取消")

		})

		// ===================================================================
		Convey("删除", func() {
			req := `
				{
					"tag":"Todo",
					"Todo":{
						"todoId":"%s"
					}
				}
`

			_, err := actionByJsonStr(fmt.Sprintf(req, todoId), consts.MethodDelete)
			So(err, ShouldBeNil)

			one, err := m.Clone().One(g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			So(one.IsEmpty(), ShouldBeTrue)

			// 物理删除测试数据
			g.DB().Model("t_todo").Unscoped().Delete(g.Map{"todo_id": todoId})

		})
	})
}

func TestActionMoreTableMoreLine(t *testing.T) {
	iAmWM()
	todoId := ""
	Convey("多表多数据操作", t, func() {

		m := g.DB().Model("t_todo").Ctx(ctx)

		// ===================================================================
		Convey("新增", func() {

			cnt1, err := m.Clone().Count(g.Map{
				"user_id": UserIdWM,
			})
			So(err, ShouldBeNil)

			req := `
				{
					"Todo": {
						"title": "去找林云喝茶 ♪(^∇^*)"
					},
					"TodoLog":{
						"log":"created by one"
					},
					"TodoLog[]":[
						{"log":"created by list[0]"},
						{"log":"created by list[1]"}
					],
					"tag": "Todo",
					"version": 2
				}
			`

			out, err := actionByJsonStr(req, consts.MethodPost)
			So(err, ShouldBeNil)

			//g.Dump(out)
			todo := out["Todo"].(g.Map)
			todoId = todo["todoId"].(string)

			cnt2, err := m.Clone().Count(g.Map{
				"user_id": UserIdWM,
			})
			So(err, ShouldBeNil)
			So(cnt2-cnt1, ShouldEqual, 1)

			cnt, err := g.DB().Model("t_todo_log").Ctx(ctx).Count(g.Map{
				"todo_id": todoId,
				"log":     "created by one",
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 1)

			cnt, err = g.DB().Model("t_todo_log").Ctx(ctx).WhereLike("log", "created by list%").Count(g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 2)

		})

		// ===================================================================
		Convey("修改", func() {
			oneId, err := g.DB().Model("t_todo_log").Ctx(ctx).Value("id", g.Map{
				"todo_id": todoId,
				"log":     "created by one",
			})
			So(err, ShouldBeNil)

			manyId, err := g.DB().Model("t_todo_log").Ctx(ctx).WhereLike("log", "created by list%").Array("id", g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			many0 := manyId[0]
			many1 := manyId[1]

			allIdStr := gjson.New([]int{oneId.Int(), many0.Int(), many1.Int()}).MustToJsonString()

			req := `
							{
								"tag":"TodoLog[]",
								"TodoLog":{
									"id{}":%s,
									"remark":"update all"
								},
								"TodoLog[]":[
									{"log":"update by one","id":"%d"},
									{"log":"update by list[0]","id":"%d"},
									{"log":"update by list[1]","id":"%d"}
								]
							}
			`

			_, err = actionByJsonStr(fmt.Sprintf(req, allIdStr, oneId.Int(), many0.Int(), many1.Int()), consts.MethodPut)
			So(err, ShouldBeNil)

			cnt, err := g.DB().Model("t_todo_log").Ctx(ctx).Count(g.Map{
				"todo_id": todoId,
				"remark":  "update all",
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 3)

			cnt, err = g.DB().Model("t_todo_log").Ctx(ctx).Count(g.Map{
				"todo_id": todoId,
				"log":     "update by one",
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 1)

			cnt, err = g.DB().Model("t_todo_log").Ctx(ctx).WhereLike("log", "update by list%").Count(g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 2)

		})

		// ===================================================================
		Convey("删除", func() {

			oneId, err := g.DB().Model("t_todo_log").Ctx(ctx).Value("id", g.Map{
				"todo_id": todoId,
				"log":     "update by one",
			})
			So(err, ShouldBeNil)

			manyId, err := g.DB().Model("t_todo_log").Ctx(ctx).WhereLike("log", "update by list%").Array("id", g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			many0 := manyId[0]
			many1 := manyId[1]

			allIdStr := gjson.New([]int{oneId.Int(), many0.Int(), many1.Int()}).MustToJsonString()

			req := `
							{
								"tag":"TodoLog",
								"TodoLog":{
									"id{}":%s
								}
							}
			`

			_, err = actionByJsonStr(fmt.Sprintf(req, allIdStr), consts.MethodDelete)
			So(err, ShouldBeNil)

			cnt, err := g.DB().Model("t_todo_log").Ctx(ctx).Count(g.Map{
				"todo_id": todoId,
			})
			So(err, ShouldBeNil)
			So(cnt, ShouldEqual, 0)

		})
	})
}

func TestActionDEmptyRowKey(t *testing.T) {
	iAmWM()
	Convey("条件为空的情况", t, func() {

		// ===================================================================
		Convey("修改", func() {

			req := `
				{
					"tag":"Todo",
					"Todo":{
						"todoId":"",
						"title":"去找林云喝茶, 把史强的预约先取消"
					}
				}
`

			_, err := actionByJsonStr(req, consts.MethodPut)
			So(err, ShouldNotBeNil)

		})

		// ===================================================================
		Convey("删除", func() {
			req := `
				{
					"tag":"Todo",
					"Todo":{
						"todoId":""
					}
				}
`

			_, err := actionByJsonStr(req, consts.MethodDelete)
			So(err, ShouldNotBeNil)

		})
	})
}
