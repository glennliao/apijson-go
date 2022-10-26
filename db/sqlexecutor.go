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

	Columns []string
	Order   string
	Group   string

	WithEmptyResult bool // 是否最终为空结果
}

func NewSqlExecutor(ctx context.Context, tableName string, accessVerify bool) (*SqlExecutor, error) {

	m := g.DB().Model(tableName)

	return &SqlExecutor{
		ctx:     ctx,
		Table:   tableName,
		m:       m,
		builder: m.Builder(),
		Columns: nil,
		Order:   "",
		Group:   "",
	}, nil
}

func (e *SqlExecutor) ParseCondition(conditions g.MapStrAny) error {

	for k, condition := range conditions {
		switch {
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

func (e *SqlExecutor) ParseCtrl(m g.Map) error {

	for k, v := range m {
		switch k {

		case "@order":
			order := strings.Replace(gconv.String(v), "-", " desc", -1)
			order = strings.Replace(order, "+", " ", -1)
			e.Order = order

		case "@column":
			columns := gconv.String(v)
			columns = strings.Replace(columns, ";", ", ", -1)
			columns = strings.Replace(columns, ":", " as ", -1)
			e.Columns = strings.Split(columns, ",")

		case "@group":
			e.Group = gconv.String(v)
		}
	}

	return nil
}

func (e *SqlExecutor) build() *gdb.Model {
	m := e.m.Clone()

	if e.Order != "" {
		m = m.Order(e.Order)
	}

	m = m.Where(e.builder)

	if e.Group != "" {
		m = m.Group(e.Group)
	}

	return m
}

func (e *SqlExecutor) List(page int, count int, needTotal bool) (list []g.Map, total int, err error) {

	if e.WithEmptyResult {
		return nil, 0, err
	}

	m := e.build()

	if needTotal {
		total, err = m.Fields("*").Count()
		if err != nil {
			return nil, 0, err
		}
	}

	// 无需下一步查询
	if needTotal && total == 0 {
		return nil, 0, err
	}

	if e.Columns != nil {
		m = m.Fields(e.Columns)
	}

	m = m.Page(page, count)
	all, err := m.All()
	if err != nil {
		return nil, 0, err
	}

	for _, item := range all.List() {
		list = append(list, item)
	}

	return
}

func (e *SqlExecutor) One() (g.Map, error) {
	if e.WithEmptyResult {
		return nil, nil
	}

	m := e.build()

	if e.Columns != nil {
		m = m.Fields(e.Columns)
	}

	one, err := m.One()

	return one.Map(), err
}
