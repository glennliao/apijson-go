package action

import (
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
)

func (a *Action) parse() {
	structures := a.tagRequestConfig.Structure

	a.children = make(map[string]*Node)
	a.keyNode = make(map[string]*Node)

	for key, v := range a.req {

		if lo.Contains([]string{consts.Tag, consts.Version}, key) {
			continue
		}

		structure, ok := structures[key]
		if !ok {
			a.err = consts.NewStructureKeyNoFoundErr(key)
			return
		}

		val := gconv.Map(v)

		node := newNode(a, key, val, structure)

		a.keyNode[key] = node
		a.children[key] = node

		err := node.parse()
		if err != nil {
			a.err = consts.NewStructureKeyNoFoundErr(key)
			return
		}
	}
}
