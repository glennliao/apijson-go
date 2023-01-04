package tests

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFailExample(t *testing.T) {
	Convey("TestFailExample this case must be err", t, func() {
		// ============================================================
		Convey("循环依赖(死锁) - 依赖的节点依赖自己", func() {
			req := `
				{
					"User":{
						"userId@":"Todo/userId"
					},
					"Todo":{
						"userId@":"User/userId"
					}
				}
			`
			_, err := queryByJsonStr(iAmUnKnow(), req)
			So(err, ShouldNotBeNil)
		})
		// ============================================================
		Convey("循环依赖(死锁) - 依赖的节点最终构成一个圈", func() {
			req := `
				{
					"User":{
						"userId@":"Todo/userId"
					},
					"Todo":{
						"userId@":"Notice/userId"
					},
					"Notice":{
						"id@":"User/id"
					}
				}
			`
			_, err := queryByJsonStr(iAmUnKnow(), req)
			So(err, ShouldNotBeNil)
		})
	})
}
