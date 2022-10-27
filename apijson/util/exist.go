package util

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"strings"
)

func Exist(m *gdb.Model, msg ...string) error {
	one, err := m.Fields("1 AS v").One()
	if err != nil {
		return err
	}
	if one.IsEmpty() {
		return gerror.New("Not Found: " + strings.Join(msg, ", "))
	}

	return nil
}

func NotExist(m *gdb.Model, msg ...string) error {
	one, err := m.Fields("1 AS v").One()
	if err != nil {
		return err
	}
	if !one.IsEmpty() {
		return gerror.New("Duplication: " + strings.Join(msg, ", "))
	}

	return nil
}
