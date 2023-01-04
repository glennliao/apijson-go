package tests

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestTodoWithUser 两表关联查询
func TestTodoWithUser(t *testing.T) {
	Convey("TodoWithUser", t, func() {
		req := `
		{
			"Todo": {
				
			},
			"User":{
				"userId@":"Todo/userId"
			}
		}
`
		ctx := iAmUnKnow()
		out, err := queryByJsonStr(ctx, req)
		So(err, ShouldBeNil)

		todo, user := out["Todo"].(g.Map), out["User"].(g.Map)
		So(todo, ShouldNotBeNil)
		So(user, ShouldNotBeNil)

		So(user["userId"], ShouldNotBeEmpty)

		row, err := g.DB().Model("t_user").Ctx(ctx).Where("user_id", user["userId"]).One()
		So(err, ShouldBeNil)

		//field := config.GetDbFieldStyle()(ctx, "t_todo", "userId")
		So(todo["userId"], ShouldEqual, row["user_id"].String())
	})
}

// TestTodoListWithUser 两表关联查询
func TestTodoListWithUser(t *testing.T) {
	Convey("", t, func() {
		req := `
		{
			"[]": {
				"Todo": {
				
				},
				"User": {
					"userId@": "[]/Todo/userId"
				}
			}
		}
`
		ctx := iAmUnKnow()
		out, err := queryByJsonStr(ctx, req)
		So(err, ShouldBeNil)

		data := out["[]"]
		So(data, ShouldNotBeNil)
		list := data.([]g.Map)

		for _, item := range list {
			todo, user := item["Todo"].(g.Map), item["User"].(g.Map)
			So(todo, ShouldNotBeNil)
			So(user, ShouldNotBeNil)

			So(user["userId"], ShouldNotBeEmpty)
			So(user["userId"], ShouldEqual, todo["userId"])
		}
		fmt.Println()

	})
}

// TestTodoListByUser 两表关联查询
func TestTodoListByUser(t *testing.T) {
	Convey("TodoListByUser", t, func() {
		req := `
		{
			"User": {
			
			},
			"[]": {
				"Todo": {
					"userId@": "User/userId"
				}
			}
		}
`
		out, err := queryByJsonStr(iAmUnKnow(), req)
		So(err, ShouldBeNil)

		data, user := out["[]"], out["User"]
		So(data, ShouldNotBeNil)
		So(user, ShouldNotBeNil)

		userId := user.(g.Map)["userId"]
		list := data.([]g.Map)
		So(len(list), ShouldBeGreaterThan, 0)

		for _, item := range list {
			row := item["Todo"].(g.Map)
			So(row["userId"], ShouldEqual, userId)
		}
	})
}

// TestTodoRef 两表关联查询
func TestTodoRef(t *testing.T) {
	Convey("TodoRef", t, func() {
		req := `
		{
			"Todo": {
				"@column": "id,userId",
				"userId@":"User/userId"
			},
			"User": {
				"@column": "userId"
			},
			"[]": {
				"Todo": {
					"userId@": "Todo/userId",
					"@column": "id,userId"
				},
				"User": {
					"user_id@": "/Todo/userId",
					"@column": "userId"
				}
			}
		}
`
		out, err := queryByJsonStr(iAmUnKnow(), req)
		So(err, ShouldBeNil)

		{
			todo := out["Todo"]
			So(todo, ShouldNotBeNil)

			row := todo.(g.Map)
			So(len(row), ShouldEqual, 2)
			So(hasKey(row, "id"), ShouldBeTrue)
			So(hasKey(row, "userId"), ShouldBeTrue)
		}

		{
			user := out["User"]
			So(user, ShouldNotBeNil)

			row := user.(g.Map)
			So(len(row), ShouldEqual, 1)
			So(hasKey(row, "userId"), ShouldBeTrue)
		}

		{
			data := out["[]"]
			list := data.([]g.Map)
			So(len(list), ShouldBeGreaterThan, 0)

			for _, item := range list {
				todo := item["Todo"]
				So(hasKey(item, "Todo"), ShouldBeTrue)
				So(len(item), ShouldEqual, 1)

				row := todo.(g.Map)

				So(len(row), ShouldEqual, 2)
				So(hasKey(row, "id"), ShouldBeTrue)
				So(hasKey(row, "userId"), ShouldBeTrue)
			}
		}
	})

}

// TestTodoOneMany 列表中一对多
func TestTodoOneMany(t *testing.T) {
	Convey("TodoOneMany", t, func() {

		req := `
		{
			"[]":{
				"User":{
		
				},
				"Todo[]":{
					"userId@":"/User/userId"
				}
			}
		}
`
		out, err := queryByJsonStr(iAmUnKnow(), req)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 1)

		So(hasKey(out, "[]"), ShouldBeTrue)
		list := out["[]"].([]g.Map)
		So(len(list), ShouldBeGreaterThan, 0)

		for _, item := range list {
			So(len(item), ShouldBeGreaterThan, 0)

			So(hasKey(item, "User"), ShouldBeTrue)
			user := item["User"].(g.Map)
			userId := user["userId"]
			So(userId, ShouldNotBeNil)

			if hasKey(item, "Todo[]") {
				todoList := item["Todo[]"].([]g.Map)
				So(len(todoList), ShouldBeGreaterThan, 0)

				for _, todo := range todoList {
					So(todo["userId"], ShouldEqual, userId)
				}
			}
		}
	})
}
