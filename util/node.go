package util

import (
	"github.com/glennliao/apijson-go/consts"
	"github.com/gogf/gf/v2/errors/gerror"
	"path/filepath"
	"strings"
)

func ParseNodeKey(inK string) (k string, isList bool) {
	k = inK

	if strings.HasSuffix(k, consts.ListKeySuffix) {
		isList = true
		k = k[0 : len(k)-len(consts.ListKeySuffix)]
	}
	return
}

// ParseRefCol 解析引用字段
// 将 "id@":"[]/User/userId"  解析出引用信息
func ParseRefCol(refStr string) (refPath string, refCol string) {
	refCol = filepath.Base(refStr)                  // userId
	refPath = refStr[0 : len(refStr)-len(refCol)-1] // []/User
	return refPath, refCol
}

// AnalysisOrder 使用拓扑排序 分析节点fetch优先级
func AnalysisOrder(prerequisites [][]string) ([]string, error) {

	var pointMap = make(map[string]bool)
	for _, prerequisite := range prerequisites {
		pointMap[prerequisite[0]] = true
		pointMap[prerequisite[1]] = true
	}

	var pointNum = len(pointMap)
	var edgesMap = make(map[string][]string)
	var inDeg = make(map[string]int)
	var result []string

	for _, prerequisite := range prerequisites {
		edgesMap[prerequisite[1]] = append(edgesMap[prerequisite[1]], prerequisite[0])
		inDeg[prerequisite[0]]++
	}

	var queue []string

	for point, _ := range pointMap {
		if inDeg[point] == 0 {
			queue = append(queue, point)
		}
	}

	for len(queue) > 0 {
		var first string
		first, queue = queue[0], queue[1:]
		result = append(result, first)
		for _, point := range edgesMap[first] {
			inDeg[point]--
			if inDeg[point] == 0 {
				queue = append(queue, point)
			}
		}
	}

	if len(result) != pointNum {
		return nil, gerror.New("依赖循环, 请检查请求")
	}

	return result, nil

}
