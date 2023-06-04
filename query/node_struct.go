package query

import (
	"github.com/glennliao/apijson-go/consts"
	"github.com/glennliao/apijson-go/model"
	"github.com/gogf/gf/v2/errors/gerror"
	"path/filepath"
	"strings"
)

type structNode struct {
	node *Node
}

func newStructNode(n *Node) *structNode {
	return &structNode{node: n}
}

func (h *structNode) parse() {
	n := h.node

	for _, childNode := range n.children {
		childNode.parse()
	}

	if n.isList { // []节点

		hasPrimary := false // 是否存在主查询表
		for _, child := range n.children {

			if child.err != nil {
				n.err = child.err
				return
			}

			if child.primaryTableKey != "" {

				if hasPrimary {
					panic(gerror.Newf("node must only one primary table: (%s)", n.Path))
				}

				hasPrimary = true
				n.primaryTableKey = filepath.Base(child.Key)
				child.page = n.page

			}
		}

		if n.Key == consts.ListKeySuffix && !hasPrimary {
			panic(gerror.Newf("node must have  primary table: (%s)", n.Path))
		}
	}

}

func (h *structNode) fetch() {
	n := h.node
	// 目前结构节点组装数据在result, 如果被依赖的是组装后的, 则无法查询。 如遇到此情况再看
	if n.isList && n.needTotal {
		n.total = n.children[n.primaryTableKey].total
	}
}

func (h *structNode) result() {
	n := h.node
	if n.isList {
		var retList []model.Map

		var primaryList []model.Map

		if n.children[n.primaryTableKey].ret != nil {
			primaryList = n.children[n.primaryTableKey].ret.([]model.Map)

			for i := 0; i < len(primaryList); i++ {

				pItem := primaryList[i]

				item := model.Map{
					n.primaryTableKey: pItem,
				}

				// 遍历组装数据, 后续考虑使用别的方案优化 (暂未简单使用map的id->item ,主要考虑多字段问题)
				for childK, childNode := range n.children {
					if childNode.primaryTableKey == "" {
						if childNode.ret != nil {

							var resultList []model.Map

							for _, depRetItem := range childNode.ret.([]model.Map) {
								match := true
								for refK, refNode := range childNode.refKeyMap {
									if pItem[refNode.column] != depRetItem[refK] {
										match = false
										break
									}
								}
								if match {
									resultList = append(resultList, depRetItem)
								}

							}
							if len(resultList) > 0 {
								if strings.HasSuffix(childK, consts.ListKeySuffix) {
									item[childK] = resultList
								} else {
									item[childK] = resultList[0]
								}
							}

						}

					}
				}

				retList = append(retList, item)
			}
		}

		n.ret = retList

		if len(n.ret.([]model.Map)) == 0 {
			n.ret = []model.Map{}
		}
	} else {

		retMap := model.Map{}
		for k, node := range n.children {
			var err error

			if strings.HasSuffix(k, consts.RefKeySuffix) {
				k = k[0 : len(k)-1]
			}
			if strings.HasSuffix(k, consts.FunctionsKeySuffix) {
				k = k[0 : len(k)-2]
			}

			// todo 增加alias ？用来重命名返回的key，避免前端调整取值
			if node.req["@alias"] != nil {
				k = node.req["@alias"].(string)
			}

			retMap[k], err = node.Result()
			if node.Type == NodeTypeFunc && retMap[k] == nil {
				delete(retMap, k)
			}

			if err != nil {
				n.err = err
			}
		}
		n.ret = retMap
	}
}

func (h *structNode) nodeType() int {
	return NodeTypeStruct
}
