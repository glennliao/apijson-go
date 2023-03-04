package framework_goframe

import "github.com/gogf/gf/v2/container/gmap"

type Mode = func(data *gmap.ListMap, meta *gmap.ListMap) gmap.ListMap

func SpreadMode(data *gmap.ListMap, meta *gmap.ListMap) gmap.ListMap {

	res := gmap.ListMap{}
	for _, k := range data.Keys() {
		res.Set(k, data.Get(k))
	}
	for _, k := range meta.Keys() {
		res.Set(k, meta.Get(k))
	}

	return res
}

func InDataMode(data *gmap.ListMap, meta *gmap.ListMap) gmap.ListMap {
	res := gmap.ListMap{}
	res.Set("data", data)
	for _, k := range meta.Keys() {
		res.Set(k, meta.Get(k))
	}
	return res
}
