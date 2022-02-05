package gin

import (
	"mime/multipart"
)

type IRequest interface {
	// region request param info
	// 返回值说明：第2个返回值表示返回的第一个参数值的来源，true: 从参数中获取；false: 返回的是默认值def

	// 1. query param, example: '/user/list?name=Bob&age=17&friends=["a", "b"]'
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat32(key string, def float32) (float32, bool)
	DefaultQueryFloat64(key string, def float64) (float64, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStrSlice(key string, def []string) ([]string, bool)
	DefaultQueryIntSlice(key string, def []int) ([]int, bool)

	// 2. path param, example: '/user/:id/name'
	DefaultParam(key string) (interface{}, bool)
	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat32(key string, def float32) (float32, bool)
	DefaultParamFloat64(key string, def float64) (float64, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	DefaultParamString(key string, def string) (string, bool)

	// 3. req-body param, example: '{"name": "Bob"}'
	DefaultForm(key string) (interface{}, bool)
	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat32(key string, def float32) (float32, bool)
	DefaultFormFloat64(key string, def float64) (float64, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormString(key string, def string) (string, bool)
	DefaultFormStrSlice(key string, def []string) ([]string, bool)
	DefaultFormIntSlice(key string, def []int) ([]int, bool)
	DefaultFormFile(key string) (*multipart.FileHeader, error)

	// 4. other(json, xml, get raw-data...)
	BindJson(v interface{}) error
	BindXml(v interface{}) error
	// end-region request param info

	// region request header info

	// 1. basic info(method, path, remote-ip, request-host)
	Method() string
	Uri() string
	Host() string
	ClientIp() string

	// 2. cookie
	Cookies() map[string]string

	// 3. header
	Headers() map[string][]string
	// end-region request header info
}
