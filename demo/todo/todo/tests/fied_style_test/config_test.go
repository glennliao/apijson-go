package fied_style_test

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/db"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"todo/todo"
)

var ctx = context.TODO()

func init() {
	db.Init()
	config.DefaultRoleFunc = todo.Role
	config.AccessConditionFunc = todo.AccessCondition
	config.AccessVerify = true // 全局配置验证权限开关

	g.Log().SetLevelStr("all")

	iamWangmiao()
	iamShiqiang()
}

// iamWangmiao 使用汪淼账号
func iamWangmiao() {
	ctx = context.WithValue(context.TODO(), config.UserIdKey, &todo.CurrentUser{UserId: "10001"})
}

// iamShiqiang 使用史强账号
func iamShiqiang() {
	ctx = context.WithValue(context.TODO(), config.UserIdKey, &todo.CurrentUser{UserId: "10002"})
}
