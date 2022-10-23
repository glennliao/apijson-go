package db

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

type SqlExecutor struct {
	ctx     context.Context
	Table   string
	m       *gdb.Model
	builder *gdb.WhereBuilder
	Page    int
	Count   int

	Columns []string
	Order   string
	Group   string

	isList bool
}

func NewSqlExecutor(ctx context.Context, table string, isList bool) *SqlExecutor {
	// table = gstr.CaseSnake(table)

	m := g.DB().Model(table)

	return &SqlExecutor{
		ctx:     ctx,
		Table:   table,
		m:       m,
		builder: m.Builder(),
		Page:    1,
		Count:   10,
		Columns: nil,
		Order:   "",
		Group:   "",
		isList:  isList,
	}
}

func (e *SqlExecutor) ParseCondition(conditions g.MapStrAny) error {

	for k, condition := range conditions {
		if strings.HasPrefix(k, "//") { // for debug, 如果字段//开头, 则忽略, 用于json"注释"
			continue
		}

		// k = gstr.CaseSnake(k)

		switch {
		case k == "page":
			e.parseCtrlCondition(k, condition)

		case k == "count":
			e.parseCtrlCondition(k, condition)

		case strings.HasPrefix(k, "@"):
			e.parseCtrlCondition(k, condition)

		case strings.HasSuffix(k, "{}"):
			e.parseMultiCondition(k[0:len(k)-2], condition)

		case strings.HasSuffix(k, "$"):
			e.builder = e.builder.WhereLike(k[0:len(k)-1], gconv.String(condition))
		case strings.HasSuffix(k, "~"):
			e.builder = e.builder.Where(k[0:len(k)-1]+" REGEXP ", gconv.String(condition))

		default:
			e.builder = e.builder.Where(k, condition)
		}
	}
	return nil
}

func (e *SqlExecutor) parseMultiCondition(k string, condition any) {

	var conditions [][]string

	if _str, ok := condition.(string); ok {
		for _, s := range strings.Split(_str, ",") {

			var item []string

			ops := []string{"<=", "<", ">=", ">"}
			isEq := true
			for _, op := range ops {
				if strings.HasPrefix(s, op) {
					item = append(item, op, s[len(op):])
					isEq = false
					break
				}
			}
			if isEq {
				item = append(item, " = ", s)
			}

			conditions = append(conditions, item)
		}

	}

	getK := func(k string) string {
		return k[0 : len(k)-1]
	}

	switch k[len(k)-1] {
	case '&':
		b := e.m.Builder()
		for _, c := range conditions {
			b = b.Where(getK(k)+" "+c[0], c[1])
		}
		e.builder = e.builder.Where(b)
	case '|':
		b := e.m.Builder()
		for _, c := range conditions {
			b = b.WhereOr(getK(k)+" "+c[0], c[1])
		}
		e.builder = e.builder.Where(b)

	case '!':
		e.builder = e.builder.WhereNotIn(getK(k), condition)
	default:
		e.builder = e.builder.WhereIn(k, condition)
	}

}

func (e *SqlExecutor) parseCtrlCondition(k string, condition any) error {

	switch k {
	case "count":
		e.Count = gconv.Int(condition)
	case "page":
		e.Page = gconv.Int(condition)
	case "@order":
		order := strings.Replace(gconv.String(condition), "-", " desc", -1)
		order = strings.Replace(order, "+", " ", -1)
		e.Order = order
	case "@column":
		columns := gconv.String(condition)
		columns = strings.Replace(columns, ";", ", ", -1)
		columns = strings.Replace(columns, ":", " as ", -1)
		e.Columns = strings.Split(columns, ",")
	case "@group":
		e.Group = gconv.String(condition)
	}
	return nil
}

func (e *SqlExecutor) Fetch() (any, error) {

	if e.Columns != nil {
		e.m = e.m.Fields(e.Columns)
	}

	if e.Order != "" {
		e.m = e.m.Order(e.Order)
	}

	e.m = e.m.Where(e.builder)

	if e.Group != "" {
		e.m = e.m.Group(e.Group)
	}

	if e.isList {
		e.m = e.m.Page(e.Page, e.Count)
		return e.m.All()
	}

	return e.m.One()
}
