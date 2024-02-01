package config

const (
	OpEq = iota
	OpNotEq
	OpIn
	OpRaw
)

type WhereItem struct {
	Args   interface{}
	Column string
	Op     int
}

type ConditionRet struct {
	conditionList []WhereItem
	builder       []*ConditionRet
	isEmptyResult bool
}

func NewConditionRet() *ConditionRet {
	c := ConditionRet{}
	return &c
}

func (c *ConditionRet) Eq(k string, v any) {
	c.where(OpEq, k, v)
}

func (c *ConditionRet) NotEq(k string, v any) {
	c.where(OpNotEq, k, v)
}

func (c *ConditionRet) In(k string, v any) {
	c.where(OpIn, k, v)
}

func (c *ConditionRet) where(op int, k string, v any) {
	item := WhereItem{
		Op:     op,
		Args:   v,
		Column: k,
	}
	c.conditionList = append(c.conditionList, item)
}

//func (c *ConditionRet) orWhere(op int, k string, v any) {
//	// TODO OR
//	item := WhereItem{
//		Op:     op,
//		Args:   v,
//		Column: k,
//	}
//	c.conditionList = append(c.conditionList, item)
//}

func (c *ConditionRet) Raw(k string, v ...any) {
	item := WhereItem{
		Column: k,
		Args:   v,
	}
	c.conditionList = append(c.conditionList, item)
}

func (c *ConditionRet) OrRaw(k string, v ...any) {
	prefix := " OR "
	if len(c.conditionList) == 0 {
		prefix = ""
	}
	item := WhereItem{
		Column: prefix + k,
		Args:   v,
	}
	c.conditionList = append(c.conditionList, item)
}

// SetEmptyResult 设置结果为空,不进行实际查询(例如sql)
func (c *ConditionRet) SetEmptyResult() {
	c.isEmptyResult = true
}

func (c *ConditionRet) IsEmptyResult() bool {
	return c.isEmptyResult
}

func (c *ConditionRet) AllCondition() ([]WhereItem, []*ConditionRet) {
	return c.conditionList, c.builder
}

func (c *ConditionRet) NewBuilder() *ConditionRet {
	b := &ConditionRet{}
	c.builder = append(c.builder, b)
	return b
}
