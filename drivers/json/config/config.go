package config

import (
	"context"

	"github.com/glennliao/apijson-go/config"
	"github.com/gogf/gf/v2/util/gconv"
)

func RequestListProvider(ctx context.Context, jsonStr string) config.RequestListProvider {

	return func(ctx context.Context) []config.RequestConfig {
		var requestList []config.RequestConfig
		err := gconv.Scan(jsonStr, &requestList)
		if err != nil {
			panic(err)
		}
		for i, request := range requestList {
			if _, ok := request.Structure[request.Tag]; !ok {
				requestList[i].Structure = map[string]*config.Structure{
					request.Tag: {
						Must:    nil,
						Refuse:  nil,
						Unique:  nil,
						Insert:  nil,
						Update:  nil,
						Replace: nil,
						Remove:  nil,
					},
				}
			}
		}
		return requestList
	}

}

func AccessListProvider(ctx context.Context, jsonStr string) config.AccessListProvider {
	return func(ctx context.Context) []config.AccessConfig {
		var accessList []config.AccessConfig
		err := gconv.Scan(jsonStr, &accessList)
		if err != nil {
			panic(err)
		}
		return accessList
	}
}
