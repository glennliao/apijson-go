package service

import (
	"context"
	"fmt"
	. "github.com/glennliao/apijson-go/apijson/app/ind/internal/model"
	"github.com/glennliao/apijson-go/apijson/internal/dao"
	"github.com/glennliao/apijson-go/apijson/internal/do"
	"github.com/glennliao/apijson-go/apijson/util"
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"strings"
)

type tableService struct{}

var _table = tableService{}

func Table() *tableService {
	return &_table
}

const (
	NormalIndex  = "INDEX"
	UniqueIndex  = "UNIQUE INDEX"
	PrimaryIndex = "PRIMARY KEY"
)

func (s *tableService) Add(ctx context.Context, in TableAddIn) (out TableAddOut, err error) {
	var columns = make([]string, 0, len(in.Columns))
	var index = map[Index][]string{}

	for _, c := range in.Columns {
		notNull := ""
		if c.NotNull {
			notNull = "NOT NULL"
		}

		def := ""
		if c.Default != "" {
			def = "DEFAULT " + c.Default
		}

		autoIncrement := ""
		if c.AutoIncrement && c.Type == "int" {
			autoIncrement = "AUTO_INCREMENT"
		}

		column := fmt.Sprintf("%s %s %s %s %s COMMENT '%s'", c.Field, c.Type, notNull, def, autoIncrement, c.Comment)
		columns = append(columns, column)

		index[c.Index] = append(index[c.Index], c.Field)
	}

	for idx, fields := range index {
		idxName := ""
		switch idx.IndexType {
		case NormalIndex:
			idxName = "idx_" + idx.IndexName
		case UniqueIndex:
			idxName = "udx_" + idx.IndexName
		}

		column := fmt.Sprintf("%s %s (%s)", idx.IndexType, idxName, strings.Join(fields, ","))
		columns = append(columns, column)
	}

	sql := fmt.Sprintf("CREATE TABLE %s (\n"+
		"%s"+
		") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='%s'", in.TableName, strings.Join(columns, ","), in.Comment)

	fmt.Println(sql)

	//ret, err := g.DB().Exec(ctx, sql)
	//return TableAddOut{Result: ret}, err

	return
}

func (s *tableService) Exist(ctx context.Context, in TableExistIn) (out TableExistOut, err error) {
	var (
		db = g.DB().Schema("INFORMATION_SCHEMA")
		m  = db.Model("TABLES").Safe().Ctx(ctx)
	)

	if err := util.Exist(m.Where("table_name", in.TableName), "表不存在:"+in.TableName); err != nil {
		return out, err
	}

	return out, nil
}

func (s *tableService) Sync(ctx context.Context, in TableSyncIn) (out TableSyncOut, err error) {
	if _, err = s.Exist(ctx, in.TableExistIn); err != nil {
		return out, err
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		title := gstr.CaseCamel(in.TableName)
		tag, alias := in.Tag, in.TableAlias
		if tag == "" {
			tag = title
		}
		if alias == "" {
			alias = title
		}

		{
			var (
				m, c    = dao.Request.MC(ctx)
				methods = []string{consts.MethodPost, consts.MethodDelete, consts.MethodPut}
			)

			if method, err := m.Fields(c.Method).
				Where(do.Request{Version: 1, Tag: tag}).
				WhereIn(c.Method, methods).Value(); err != nil {
				return err
			} else if !method.IsNil() {
				return gerror.New("Method Already Exist: " + method.String())
			}

			var data = make([]do.Request, 0, len(methods))
			for _, method := range methods {
				data = append(data, do.Request{
					Debug:     0,
					Version:   1,
					Method:    method,
					Tag:       tag,
					Structure: g.Map{},
				})
			}

			_, err := m.Insert(data)
			if err != nil {
				return err
			}
		}
		{
			var (
				m, _ = dao.Access.MC(ctx)
			)

			if err := util.NotExist(m.
				Where(do.Access{Name: in.TableName}).
				WhereOr(do.Access{Alias: alias}),
			); err != nil {
				return err
			}

			var data = do.Access{
				Debug: 0,
				Name:  in.TableName,
				Alias: alias,
			}

			_, err := m.Insert(data)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return out, err
}
