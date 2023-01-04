package tests

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestFunctionsQuery(t *testing.T) {
	ctx := iAmWM()
	Convey("functions in query", t, func() {

		// ===================================================================
		Convey("sayHello", func() {

			req := `
				{
					"User": {},
					"hello()":"sayHello"
				}
			`
			out, err := queryByJsonStr(ctx, req)

			So(err, ShouldBeNil)
			So(hasKey(out, "hello"), ShouldBeTrue)
			So(gconv.String(out["hello"]) == "world", ShouldBeTrue)
		})

		Convey("sayHello()", func() {
			req := `
				{
					"User": {},
					"hello()":"sayHello()"
				}
			`
			out, err := queryByJsonStr(ctx, req)

			So(err, ShouldBeNil)
			So(hasKey(out, "hello"), ShouldBeTrue)
			So(gconv.String(out["hello"]) == "world", ShouldBeTrue)
		})

		Convey("sayHi(realname)", func() {
			req := `
				{
					"User": {
						"hello()":"sayHi(realname)"
					}
				}
			`
			out, err := queryByJsonStr(ctx, req)

			//g.Dump(out)

			So(err, ShouldBeNil)
			So(hasKey(out, "User"), ShouldBeTrue)

			user := gconv.Map(out["User"])

			So(hasKey(user, "hello"), ShouldBeTrue)
			So(gconv.String(user["hello"]) == "你好:"+gconv.String(user["realname"]), ShouldBeTrue)
		})

		Convey("sayHi(realname) in List", func() {
			req := `
				{
					"User[]": {
						"hello()":"sayHi(realname)"
					}
				}
			`
			out, err := queryByJsonStr(ctx, req)

			//g.Dump(out)

			So(err, ShouldBeNil)
			userList := gconv.Maps(out["User[]"])
			So(hasKey(userList[0], "hello"), ShouldBeTrue)
			So(gconv.String(userList[0]["hello"]) == "你好:"+gconv.String(userList[0]["realname"]), ShouldBeTrue)
		})

	})
}

func TestFunctionsInAction(t *testing.T) {
	ctx := iAmSQ()
	Convey("functions in action", t, func() {
		// ===================================================================
		Convey("check", func() {

			req := `
				{
					"Todo": {
						"title": "去找林云喝茶"
					},
					"tag": "Todo",
					"version": 1
				}
			`

			_, err := actionByJsonStr(ctx, req, http.MethodPost)

			So(err, ShouldNotBeNil)

		})

		Convey("replace", func() {

			req := `
				{
					"Todo": {
						"title": "去找林云逛街"
					},
					"tag": "Todo",
					"version": 1
				}
			`

			out, err := actionByJsonStr(ctx, req, http.MethodPost)
			//g.Dump(out)
			So(err, ShouldBeNil)

			// 删除测试数据
			todo := out["Todo"].(g.Map)
			todoId := todo["todoId"].(string)

			g.DB().Model("t_todo").Unscoped().Delete(g.Map{"todo_id": todoId})

		})
	})
}
