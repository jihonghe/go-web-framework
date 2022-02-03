package summer

import (
	"mime/multipart"
)

type IRequest interface {
	// region request param info
	// 返回值说明：第2个返回值表示返回的第一个参数值的来源，true: 从参数中获取；false: 返回的是默认值def

	// 1. query param, example: '/user/list?name=Bob&age=17&friends=["a", "b"]'
	Query(key string) (interface{}, bool)
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStrSlice(key string, def []string) ([]string, bool)
	QueryIntSlice(key string, def []int) ([]int, bool)

	// 2. path param, example: '/user/:id/name'
	Param(key string) (interface{}, bool)
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)

	// 3. req-body param, example: '{"name": "Bob"}'
	Form(key string) (interface{}, bool)
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStrSlice(key string, def []string) ([]string, bool)
	FormIntSlice(key string, def []int) ([]int, bool)
	FormFile(key string) (*multipart.FileHeader, error)

	// 4. other(json, xml, get raw-data...)
	BindJson(v interface{}) error
	BindXml(v interface{}) error
	GetRawData() ([]byte, error)
	// end-region request param info

	// region request header info

	// 1. basic info(method, path, remote-ip, request-host)
	Method() string
	Uri() string
	Host() string
	ClientIp() string

	// 2. cookie
	Cookie(key string) (string, bool)
	Cookies() map[string]string

	// 3. header
	Header(key string) (string, bool)
	Headers() map[string][]string
	// end-region request header info
}
