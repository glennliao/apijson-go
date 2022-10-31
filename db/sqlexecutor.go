package db

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"regexp"
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

	var columns []string
	for _, column := range tableMap[tableName].Columns {
		columns = append(columns, column.Name)
	}

	return &SqlExecutor{
		ctx:     ctx,
		Table:   tableName,
		m:       m,
		builder: m.Builder(),
		Columns: columns,
		Order:   "",
		Group:   "",
	}, nil
}

func (e *SqlExecutor) ParseCondition(conditions g.MapStrAny) error {

	for k, condition := range conditions {
		k = config.ToDbField(k) // 将请求字段转化为数据库字段风格

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

var exp = regexp.MustCompile(`^[\s\w][\w()]+`) // 匹配 field, COUNT(field)

func (e *SqlExecutor) ParseCtrl(ctrl g.Map) error {

	for k, v := range ctrl {
		// https://github.com/Tencent/APIJSON/blob/master/Document.md
		// 应该用分号 ; 隔开 SQL 函数，改为 "@column":"store_id;sum(amt):totAmt"）
		fieldStr := strings.ReplaceAll(gconv.String(v), ";", Separator)

		fieldList := strings.Split(fieldStr, ",")
		for i, item := range fieldList {
			fieldList[i] = exp.ReplaceAllStringFunc(item, config.ToDbField) // 将请求字段转化为数据库字段风格
		}

		fieldStr = strings.Join(fieldList, Separator)

		switch k {

		case "@order":
			order := strings.ReplaceAll(fieldStr, "-", DESC)
			order = strings.ReplaceAll(order, "+", " ")
			e.Order = order

		case "@column":
			e.Columns = fieldList

		case "@group":
			e.Group = fieldStr
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
		total, err = m.Count()
		if err != nil || total == 0 {
			return nil, 0, err
		}
	}

	m = m.Fields(e.JsonFields())

	m = m.Page(page, count)
	all, err := m.All()
	if err != nil {
		return nil, 0, err
	}

	return all.List(), total, nil
}

func (e *SqlExecutor) One() (g.Map, error) {
	if e.WithEmptyResult {
		return nil, nil
	}

	m := e.build()

	m = m.Fields(e.JsonFields())

	one, err := m.One()

	return one.Map(), err
}

// JsonFields 返回 config.ToJsonField 指定的字段格式.
func (e *SqlExecutor) JsonFields() []string {

	var fields = make([]string, 0, len(e.Columns))
	for _, column := range e.Columns {
		column = strings.ReplaceAll(column, ":", AS)
		if !strings.Contains(column, AS) {
			field := config.ToJsonField(column)
			if field != column {
				column = column + AS + field
			}
		}

		fields = append(fields, column)
	}

	return fields
}
