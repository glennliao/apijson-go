package model

import (
	"github.com/gogf/gf/v2/encoding/gjson"
)

type (
	RequestListIn struct {
		ReqList
	}

	RequestListResult struct {
		Id        string      `json:"id" des:"唯一标识"`
		Debug     int8        `json:"debug" des:"是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。"`
		Version   int8        `json:"version" des:"GET,HEAD可用任意结构访问任意开放内容，不需要这个字段。其它的操作因为写入了结构和内容，所以都需要，按照不同的version选择对应的structure。"`
		Method    string      `json:"method" des:"只限于GET,HEAD外的操作方法。"`
		Tag       string      `json:"tag" des:"标签"`
		Structure *gjson.Json `json:"structure" des:"结构。"`
		Detail    string      `json:"detail" des:"详细说明"`
	}
)

type (
	RequestGetIn struct {
		Id string
	}

	RequestGetOut struct {
		Id        string      `json:"id" des:"唯一标识"`
		Debug     int8        `json:"debug" des:"是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。"`
		Version   int8        `json:"version" des:"GET,HEAD可用任意结构访问任意开放内容，不需要这个字段。其它的操作因为写入了结构和内容，所以都需要，按照不同的version选择对应的structure。"`
		Method    string      `json:"method" des:"只限于GET,HEAD外的操作方法。"`
		Tag       string      `json:"tag" des:"标签"`
		Structure *gjson.Json `json:"structure" des:"结构。"`
		Detail    string      `json:"detail" des:"详细说明"`
	}
)

type (
	RequestAddIn struct {
		Debug     *int8       `des:"是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。" d:"0"`
		Version   *int8       `des:"用于GET,HEAD外的操作方法访问结构和内容" d:"1"`
		Method    *string     `des:"只限于GET,HEAD外的操作方法。" v:"required"`
		Tag       *string     `des:"标签" v:"required"`
		Structure *gjson.Json `des:"结构。" v:"required"`
		Detail    *string     `des:"详细说明"`
	}

	RequestAddOut struct {
		Id string
	}
)

type RequestUpdateIn struct {
	Id        string      `des:"唯一标识"`
	Debug     *int8       `des:"是否为 DEBUG 调试数据，只允许在开发环境使用，测试和线上环境禁用：0-否，1-是。"`
	Version   *int8       `des:"GET,HEAD可用任意结构访问任意开放内容，不需要这个字段。其它的操作因为写入了结构和内容，所以都需要，按照不同的version选择对应的structure。"`
	Method    *string     `des:"只限于GET,HEAD外的操作方法。"`
	Tag       *string     `des:"标签"`
	Structure *gjson.Json `des:"结构。"`
	Detail    *string     `des:"详细说明"`
}

// #1db24bf040a42d23d35f34fb0ffa032bfddafb37:021026:52
