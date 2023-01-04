package tests

import (
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func hasKey(m g.Map, k string) bool {
	_, exists := m[k]
	return exists
}

func TestFieldStyle(t *testing.T) {
	oriDbStyle := config.GetDbFieldStyle()
	oriJsonStyle := config.GetJsonFieldStyle()
	defer func() {
		config.SetDbFieldStyle(oriDbStyle)
		config.SetJsonFieldStyle(oriJsonStyle)
	}()

	// 手动临时启用
	SkipConvey("TestFieldStyle", t, func() {
		// ============================================================
		Convey("json use CaseCamel, db use CaseSnake", func() {

			req := `
					{
						"Todo":{
							"@column":"user_id;userId:userId2",
							"user_id":"10001",
							"userId$":"1000%"
						}
			}
			`
			out, err := queryByJsonStr(iAmUnKnow(), req)
			//g.Dump(out)

			So(err, ShouldBeNil)

			todo := gconv.Map(out["Todo"])
			SoMsg("数据库中不能没有数据,不然无法测试", len(lo.Keys(todo)), ShouldBeGreaterThan, 0)

			So(hasKey(todo, "userId"), ShouldBeTrue)
			So(hasKey(todo, "userId2"), ShouldBeTrue)
			So(todo["userId"] == todo["userId2"], ShouldBeTrue)
			So(todo["userId"] == "10001", ShouldBeTrue)
		})

		// ============================================================
		Convey("json use nil, db use nil", func() {
			config.SetJsonFieldStyle(nil)
			config.SetDbFieldStyle(nil)

			req := `
					{
						"Todo":{
							"@column":"user_id;user_Id:userId2",
							"user_id":"10001"
						}
			}
			`
			out, err := queryByJsonStr(iAmUnKnow(), req)
			//g.Dump(out)

			So(err, ShouldBeNil)

			todo := gconv.Map(out["Todo"])
			SoMsg("数据库中不能没有数据,不然无法测试", len(lo.Keys(todo)), ShouldBeGreaterThan, 0)

			So(hasKey(todo, "user_id"), ShouldBeTrue)
			So(hasKey(todo, "userId"), ShouldBeFalse)
			So(hasKey(todo, "userId2"), ShouldBeTrue)

			So(todo["user_id"] == todo["userId2"], ShouldBeTrue)
			So(todo["user_id"] == "10001", ShouldBeTrue)
		})
	})
}
