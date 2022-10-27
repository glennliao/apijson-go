package service_test

import (
	"context"
	"fmt"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/glennliao/apijson-go/apijson/app/ind/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTableService_Add(t *testing.T) {
	ctx := context.TODO()

	in := model.TableAddIn{
		TableName: "t_test",
		Columns:   nil,
		Comment:   "test create table",
	}

	in.Columns = []model.Column{
		{
			Field:         "id",
			Type:          "int",
			Comment:       "主键id",
			Default:       "",
			NotNull:       true,
			AutoIncrement: true,
			Index: model.Index{
				IndexType: service.PrimaryIndex,
				IndexName: "",
			},
		}, {
			Field:         "str1_uni",
			Type:          "varchar(16)",
			Comment:       "str1",
			Default:       "",
			NotNull:       false,
			AutoIncrement: false,
			Index: model.Index{
				IndexType: service.UniqueIndex,
				IndexName: "str1_str2",
			},
		}, {
			Field:         "str2_uni",
			Type:          "varchar(16)",
			Comment:       "主键id",
			Default:       "",
			NotNull:       false,
			AutoIncrement: false,
			Index: model.Index{
				IndexType: service.UniqueIndex,
				IndexName: "str1_str2",
			},
		}, {
			Field:         "str4",
			Type:          "varchar(16)",
			Comment:       "主键id",
			Default:       "",
			NotNull:       false,
			AutoIncrement: false,
			Index: model.Index{
				IndexType: service.NormalIndex,
				IndexName: "str4",
			},
		},
	}

	service.Table().Add(ctx, in)
}

func TestSync(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)

	out, err := service.Table().Sync(ctx, model.TableSyncIn{
		TableExistIn: model.TableExistIn{TableName: "t_no_exist"},
	})
	a.NotNil(err, err)

	fmt.Println(out, err)

	table := "t_test" + strings.ToLower(grand.Letters(10))
	sql := fmt.Sprintf("CREATE TABLE %s (`id`int(10) NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`));", table)

	g.DB().Exec(ctx, sql)

	out, err = service.Table().Sync(ctx, model.TableSyncIn{
		TableExistIn: model.TableExistIn{TableName: table},
	})

	a.Nil(err, err)

	g.DB().Exec(ctx, "DROP TABLE "+table)
}
