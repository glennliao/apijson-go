package config_test

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"log"
	"regexp"
	"strings"
	"testing"
)

const (
	countAs = "count(userId) AS CnT,id,username, userId"
	count   = "count(userId) : CnT,id,username, userId"
)

func TestRegexp(t *testing.T) {
	//(\w|\(|\)|,)

	exp := regexp.MustCompile(`^\s+[\w()]+`)

	list := strings.Split(countAs, ",")

	for _, s := range list {
		s2 := exp.ReplaceAllStringFunc(s, func(s string) string {
			return gstr.CaseSnake(s)
		})

		log.Println(s, ":", s2)
	}
}

func TestRegexpCount(t *testing.T) {
	separatorExp := regexp.MustCompile(`[^,;]+`)
	wordExp := regexp.MustCompile(`^[\s\w][\w()]+`)

	ret := separatorExp.ReplaceAllStringFunc(count, func(s1 string) string {

		s2 := wordExp.ReplaceAllStringFunc(s1, func(s2 string) string {
			return gstr.CaseSnake(s2)
		})
		log.Println(s1, ":", s2)

		return s2
	})

	g.Dump("ret", ret)

}
