package tests

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/db"
	"github.com/glennliao/apijson-go/handlers"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"todo/app"
)

var ctx = context.TODO()

const (
	UserIdWM = "10001"
	UserIdSQ = "10002"
)

func init() {

	config.DefaultRoleFunc = app.Role
	config.AccessConditionFunc = app.AccessCondition
	config.AccessVerify = false // 全局配置验证权限开关

	g.Log().SetLevelStr("all")
	g.Log().SetLevelStr("info") // 需要显示debug时将本句注释即可

	logger := g.Log("db")      // 使用独立的Logger控制sql日志
	logger.SetLevelStr("info") // 不打印db.Init初始化的日志
	g.DB().SetLogger(logger)

	db.Init()

	logger.SetLevelStr("all")
	g.DB().SetLogger(logger)

	config.SetDbFieldStyle(config.CaseSnake)
	config.SetJsonFieldStyle(config.CaseCamel)
}

// iAmWM 使用汪淼账号
func iAmWM() {
	ctx = context.WithValue(context.TODO(), config.UserIdKey, &app.CurrentUser{UserId: UserIdWM})
}

// iAmSQ 使用史强账号
func iAmSQ() {
	ctx = context.WithValue(context.TODO(), config.UserIdKey, &app.CurrentUser{UserId: UserIdSQ})
}

// 未登录用户
func iAmUnKnow() {
	ctx = context.TODO()
}

func queryByJsonStr(req string) (res g.Map, err error) {
	reqMap := gjson.New(req).Map()
	return handlers.Get(ctx, reqMap)
}

func countTodoByUser(userId string) int {
	m := g.Model("todo").Ctx(ctx)
	if userId != "" {
		m = m.Where(g.Map{"user_id": userId})
	}
	cnt, err := m.Count()
	if err != nil {
		return -1
	}
	return cnt
}
