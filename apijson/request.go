package apijson

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"my-apijson/apijson/db"
	"my-apijson/apijson/util"
	"strings"
)

func checkByRequest(req g.Map, method string) (reqMap g.Map, err error) {
	tag, ok := req["tag"]
	if !ok {
		return nil, gerror.New("tag 缺失")
	}

	request, ok := db.RequestMap[method+"@"+gconv.String(tag)]

	if !ok {
		return nil, gerror.New("tag错误: 404")
	}

	delete(req, "tag")

	for k, v := range request.Structure {
		if reqV, ok := req[k]; !ok {
			return nil, gerror.New("structure错误: 400, 缺少" + k)
		} else {

			kStructure := v.(map[string]any)
			_reqV := reqV.(map[string]any)

			for opeK, _opeV := range kStructure {
				switch opeK {
				case "MUST":
					mustKeys := strings.Split(gconv.String(_opeV), ",")
					for _, key := range mustKeys {
						if _, ok := _reqV[key]; !ok {
							return nil, gerror.New("structure错误: 400, 缺少" + k + "." + key)
						}
					}
				case "REFUSE":

					if gconv.String(_opeV) == "!" {
						if opeV, ok := kStructure["MUST"]; ok {
							mustKeys := strings.Split(gconv.String(opeV), ",")
							if len(mustKeys) == 0 {
								return nil, gerror.New("structure错误: 400, REFUSE为!时必须指定MUST" + k)
							}

							for reqK, _ := range _reqV {

								if util.Contains(mustKeys, reqK) {
									return nil, gerror.New("structure错误: 400, 不能包含" + k + "." + reqK)
								}
							}

						} else {
							return nil, gerror.New("structure错误: 400, REFUSE为!时必须指定MUST" + k)
						}
					} else {
						keys := strings.Split(gconv.String(_opeV), ",")
						for _, key := range keys {
							if _, ok := _reqV[key]; ok {
								return nil, gerror.New("structure错误: 400, 不能包含" + k + "." + key)
							}
						}
					}

				}
			}

		}
	}

	return req, nil
}
