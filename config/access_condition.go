package config

import (
	"github.com/glennliao/apijson-go/consts"
)

type ConditionRet struct {
	condition    map[string]any
	rawCondition map[string][]any
}

func NewConditionRet() *ConditionRet {
	c := ConditionRet{
		condition:    map[string]any{},
		rawCondition: map[string][]any{},
	}
	return &c
}

func (c *ConditionRet) Add(k string, v any) {
	c.condition[k] = v
}

func (c *ConditionRet) AddRaw(k string, v ...any) {
	c.rawCondition[k] = v
}

func (c *ConditionRet) AllWhere() map[string]any {
	if len(c.rawCondition) > 0 {
		c.condition[consts.Raw] = c.rawCondition
	}
	return c.condition
}
