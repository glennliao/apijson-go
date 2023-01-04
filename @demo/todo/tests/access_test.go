package tests

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// TestAccessExt 字段RightLike查询
func TestAccessExtFiledRightLike(t *testing.T) {
	Convey("AccessExt", t, func() {

		rows, err := g.DB().Model("t_todo").All()
		So(err, ShouldBeNil)

		const layout = "2006-01-02"

		var createdAtList = map[string]int{}
		for _, row := range rows {
			createdAt := row["created_at"]
			if createdAt.IsNil() {
				continue
			}

			createdAtList[createdAt.Time().Format(layout)]++
		}

		outByDate := func(date string) g.Map {
			req := fmt.Sprintf(`
			{
				"[]":{
					"User":{
						"userId@":"/Todo/userId"
					},
					"Todo":{
						"@role":"PARTNER",
						"createdAt$":"%s%%"
					}
				}
			}
			`, date)

			out, err := queryByJsonStr(gctx.New(), req)
			So(err, ShouldBeNil)
			return out
		}

		// 制造无效数据
		now := time.Now()
		for i := 100; i < 1000; i++ {
			date := now.AddDate(0, 0, -i).Format(layout)
			if _, ok := createdAtList[date]; !ok {
				createdAtList[date] = 0
				break
			}
		}

		for date, cnt := range createdAtList {
			out := outByDate(date)
			list, ok := out["[]"].([]g.Map)

			// 默认分页：limit 0,10
			if cnt > 10 {
				cnt = 10
			}

			if ok {
				So(len(list), ShouldEqual, cnt)
			} else {
				So(out["[]"], ShouldBeEmpty)
			}
		}
	})
}
