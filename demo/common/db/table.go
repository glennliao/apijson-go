package db

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/util/grand"

	"github.com/go-faker/faker/v4"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/database/gdb"

	"github.com/glennliao/table-sync/tablesync"
)

type User struct {
	Id        uint32 `ddl:"primaryKey" faker:"unique"`
	Username  string `ddl:"size:32;comment:用户名" faker:"username,unique"`
	Nickname  string `ddl:"size:32;comment:昵称" faker:"name"`
	Password  string `ddl:"size:128;comment:密码" faker:"password"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Todo struct {
	Id        uint32 `ddl:"primaryKey"`
	UserId    uint32
	Content   string
	CreatedAt *time.Time
}

func InitTable(ctx context.Context, db gdb.DB) {
	syncer := tablesync.Syncer{
		Tables: []tablesync.Table{
			User{},
			Todo{},
		},
	}
	err := syncer.Sync(ctx, db)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	GenRandomData(ctx, db)
}

func GenRandomData(ctx context.Context, db gdb.DB) {
	for i := 0; i < 10; i++ {
		user := User{}
		_ = faker.FakeData(&user)
		db.Model("user").Insert(user)

		num := grand.N(3, 8)
		for i := 0; i < num; i++ {
			todo := Todo{}
			_ = faker.FakeData(&todo)
			todo.UserId = user.Id
			db.Model("todo").Insert(todo)
		}
	}
}
