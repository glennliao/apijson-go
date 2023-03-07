package executor_goframe

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/executor"
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/apijson-go/util"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

type SqlExecutor struct {
	ctx context.Context

	Role string

	//保存where条件 [ ["user_id",">", 123], ["user_id","<=",345] ]
	Where           [][]any
	accessCondition model.Map

	Columns []string
	Order   string
	Group   string

	// 是否最终为空结果, 用于node中中断数据获取
	WithEmptyResult bool

	config *config.ExecutorConfig
}

func New(ctx context.Context, config *config.ExecutorConfig) (executor.QueryExecutor, error) {

	return &SqlExecutor{
		ctx:             ctx,
		Where:           [][]any{},
		Columns:         nil,
		Order:           "",
		Group:           "",
		WithEmptyResult: false,
		config:          config,
	}, nil
}

// ParseCondition 解析查询条件
// accessVerify 内部调用时, 不校验是否可使用该种查询方式
func (e *SqlExecutor) ParseCondition(conditions model.MapStrAny, accessVerify bool) error {

	for key, condition := range conditions {
		switch {
		case strings.HasSuffix(key, "{}"):
			e.parseMultiCondition(util.RemoveSuffix(key, "{}"), condition)

		case strings.HasSuffix(key, consts.OpLike):
			e.Where = append(e.Where, []any{key[0 : len(key)-1], consts.SqlLike, gconv.String(condition)})

		case strings.HasSuffix(key, consts.OpRegexp):
			e.Where = append(e.Where, []any{key[0 : len(key)-1], consts.SqlRegexp, gconv.String(condition)})

		case key == consts.Raw && !accessVerify:
			e.accessCondition = condition.(model.Map)

		default:
			e.Where = append(e.Where, []any{key, consts.SqlEqual, condition})
		}
	}

	if !accessVerify {
		return nil
	}

	if e.config.NoVerify { // 可任意字段搜索
		return nil
	}

	inFieldsMap := e.config.GetFieldsGetInByRole()

	dbStyle := e.config.DbFieldStyle

	tableName := e.config.TableName()

	for _, where := range e.Where {
		k := dbStyle(e.ctx, tableName, where[0].(string))
		if val, exists := inFieldsMap[k]; exists {

			if val[0] == "*" {
				continue
			}

			op := where[1].(string)
			if op == consts.SqlLike {
				condition := where[2].(string)
				op = consts.OpLike
				if strings.HasPrefix(condition, "%") {
					op = "%" + op
				}
				if strings.HasSuffix(condition, "%") {
					op = op + "%"
				}
			}

			if !lo.Contains(val, op) {
				panic("不允许使用" + where[0].(string) + "的搜索方式:" + op)
			}

		} else {
			panic("不允许使用" + where[0].(string) + "搜索")
		}
	}

	return nil
}

// ParseCondition 解析批量查询条件
func (e *SqlExecutor) parseMultiCondition(k string, condition any) {

	var conditions [][]string
	var value = condition

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
		value = conditions
	}

	getK := func(k string) string {
		return k[0 : len(k)-1]
	}

	switch k[len(k)-1] {
	case '&', '|', '!':
		e.Where = append(e.Where, []any{getK(k), k[len(k)-1], value})
	default:
		e.Where = append(e.Where, []any{k, "in", value})

	}

}

var exp = regexp.MustCompile(`^[\s\w][\w()]+`) // 匹配 field, COUNT(field)

// ParseCtrl 解析 @column,@group等控制类
func (e *SqlExecutor) ParseCtrl(ctrl model.Map) error {

	fieldStyle := e.config.DbFieldStyle
	tableName := e.config.TableName()
	for k, v := range ctrl {
		// 使用;分割字段
		fieldStr := strings.ReplaceAll(gconv.String(v), ";", ",")

		fieldList := strings.Split(fieldStr, ",")

		for i, item := range fieldList {
			fieldList[i] = exp.ReplaceAllStringFunc(item, func(field string) string {
				return fieldStyle(e.ctx, tableName, field)
			}) // 将请求字段转化为数据库字段风格
		}

		switch k {

		case "@column":
			e.Columns = fieldList

		case "@order":
			fieldStr = strings.Join(fieldList, ",")
			order := strings.ReplaceAll(fieldStr, "-", " DESC")
			order = strings.ReplaceAll(order, "+", " ")
			e.Order = order

		case "@group":
			fieldStr = strings.Join(fieldList, ",")
			e.Group = fieldStr
		}
	}

	return nil
}

func (e *SqlExecutor) build() *gdb.Model {
	tableName := e.config.TableName()
	m := g.DB().Model(tableName).Ctx(e.ctx)

	if e.Order != "" {
		m = m.Order(e.Order)
	}

	whereBuild := m.Builder()

	fieldStyle := e.config.DbFieldStyle

	for _, whereItem := range e.Where {
		key := fieldStyle(e.ctx, tableName, whereItem[0].(string))
		op := whereItem[1]
		value := whereItem[2]

		if conditions, ok := value.([][]string); ok { // multiCondition

			switch op {
			case '&':
				b := m.Builder()
				for _, c := range conditions {
					b = b.Where(key+" "+c[0], c[1])
				}
				whereBuild = whereBuild.Where(b)

			case '|':
				b := m.Builder()
				for _, c := range conditions {
					b = b.WhereOr(key+" "+c[0], c[1])
				}
				whereBuild = whereBuild.Where(b)

			case '!':
				whereBuild = whereBuild.WhereNotIn(key, conditions)

			default:
				whereBuild = whereBuild.WhereIn(key, conditions)
			}
		} else {

			switch op {
			case consts.SqlLike:
				whereBuild = whereBuild.WhereLike(key, value.(string))
			case consts.SqlRegexp:
				whereBuild = whereBuild.Where(key+" "+consts.SqlRegexp, value.(string))
			case "in":
				whereBuild = whereBuild.WhereIn(key, value)
			case consts.SqlEqual:
				whereBuild = whereBuild.Where(key, value)
			}

		}
	}

	m = m.Where(whereBuild)
	if e.accessCondition != nil {
		m = m.Where(e.accessCondition)
	}

	if e.Group != "" {
		m = m.Group(e.Group)
	}

	return m
}

func (e *SqlExecutor) column() []string {

	outFields := e.config.GetFieldsGetOutByRole()

	tableName := e.config.TableName()

	var columns []string

	if e.Columns != nil {
		columns = e.Columns
	} else {
		columns = e.config.TableColumns()
	}

	var fields = make([]string, 0, len(columns))

	fieldStyle := e.config.JsonFieldStyle
	dbStyle := e.config.DbFieldStyle

	for _, column := range columns {
		fieldName := column
		column = strings.ReplaceAll(column, ":", " AS ")
		if !strings.Contains(column, " AS ") {
			field := fieldStyle(e.ctx, tableName, column)
			if field != column {
				column = "`" + column + "`" + " AS " + field
			} else {
				column = "`" + column + "`"
			}
		} else {
			fieldName = strings.TrimSpace(strings.Split(fieldName, "AS")[0])
		}

		// 过滤可访问字段
		if e.config.NoVerify || lo.Contains(outFields, dbStyle(e.ctx, tableName, fieldName)) ||
			len(outFields) == 0 /* 数据库中未设置, 则看成全部可访问 */ {
			fields = append(fields, column)
		}
	}

	return fields
}

func (e *SqlExecutor) SetEmptyResult() {
	e.WithEmptyResult = true
}

func (e *SqlExecutor) Count() (total int64, err error) {
	m := e.build()
	_total, err := m.Count()
	if err != nil || _total == 0 {
		return 0, err
	} else {
		total = int64(_total)
	}

	return total, nil

}

func (e *SqlExecutor) List(page int, count int) (list []model.Map, err error) {

	if e.WithEmptyResult {
		return nil, err
	}

	m := e.build()

	m = m.Fields(e.column())

	m = m.Page(page, count)
	all, err := m.All()

	if err != nil {
		return nil, err
	}

	for _, item := range all.List() {
		list = append(list, item)
	}

	return list, nil
}

func (e *SqlExecutor) One() (model.Map, error) {
	if e.WithEmptyResult {
		return nil, nil
	}

	m := e.build()

	m = m.Fields(e.column())

	one, err := m.One()

	return one.Map(), err
}
