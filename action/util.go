package action

import (
	"github.com/glennliao/apijson-go/db"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func checkTag(req g.Map, method string) (*db.Request, error) {
	_tag, ok := req["tag"]
	if !ok {
		return nil, gerror.New("tag 缺失")
	}

	tag := gconv.String(_tag)

	request, err := db.GetRequest(tag, method, -1)
	if err != nil {
		return nil, err
	}

	return request, nil
}
