package gin

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/spf13/cast"
)

func (c *Context) QueryParams() map[string][]string {
	if c.Request != nil {
		return c.Request.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) FormParams() map[string][]string {
	if c.Request != nil {
		c.Request.ParseForm()
		return c.Request.PostForm
	}
	return map[string][]string{}
}

// # region implement interface IRequest

func (c *Context) DefaultQueryInt(key string, def int) (int, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToInt(values[0]), true
	}

	return def, false
}

func (c *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToInt64(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultQueryFloat32(key string, def float32) (float32, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat32(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultQueryFloat64(key string, def float64) (float64, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat64(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToBool(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultQueryString(key string, def string) (string, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToString(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultQueryStrSlice(key string, def []string) ([]string, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return values, true
	}
	return def, false
}

func (c *Context) DefaultQueryIntSlice(key string, def []int) ([]int, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToIntSlice(values), true
	}
	return def, false
}

func (c *Context) DefaultParam(key string) (interface{}, bool) {
	if c.params != nil {
		if val, ok := c.params.Get(key); ok {
			return val, true
		}
	}
	return nil, false
}

func (c *Context) DefaultParamInt(key string, def int) (int, bool) {
	val, ok := c.DefaultParam(key)
	if ok {
		return cast.ToInt(val), true
	}
	return def, false
}

func (c *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	val, ok := c.DefaultParam(key)
	if ok {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (c *Context) DefaultParamFloat32(key string, def float32) (float32, bool) {
	val, ok := c.DefaultParam(key)
	if ok {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (c *Context) DefaultParamFloat64(key string, def float64) (float64, bool) {
	val, ok := c.DefaultParam(key)
	if ok {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (c *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	val, ok := c.DefaultParam(key)
	if ok {
		return cast.ToBool(val), true
	}
	return def, false
}

func (c *Context) DefaultParamString(key string, def string) (string, bool) {
	val, ok := c.DefaultParam(key)
	if ok {
		return cast.ToString(val), true
	}
	return def, false
}

func (c *Context) DefaultForm(key string) (interface{}, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return nil, false
}

func (c *Context) DefaultFormInt(key string, def int) (int, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToInt(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToInt64(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultFormFloat32(key string, def float32) (float32, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat32(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultFormFloat64(key string, def float64) (float64, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat64(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToBool(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultFormString(key string, def string) (string, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToString(values[0]), true
	}
	return def, false
}

func (c *Context) DefaultFormStrSlice(key string, def []string) ([]string, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToStringSlice(values), true
	}
	return def, false
}

func (c *Context) DefaultFormIntSlice(key string, def []int) ([]int, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToIntSlice(values), true
	}
	return def, false
}

func (c *Context) DefaultFormFile(key string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		return nil, c.Request.ParseMultipartForm(defaultMultipartMemory)
	}
	f, fh, err := c.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

// BindJson 将body文本解析到结构体 v 中
func (c *Context) BindJson(v interface{}) error {
	reqBody, err := c.requestBody()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(reqBody, v); err != nil {
		return errors.New(fmt.Sprintf("failed to unmarshal body of %s to type '%s' by json tag, body='%s', error: %s",
			c.RequestString(), reflect.TypeOf(v).String(), string(reqBody), err.Error()))
	}
	return nil
}

func (c *Context) BindXml(v interface{}) error {
	reqBody, err := c.requestBody()
	if err != nil {
		return err
	}
	if err = xml.Unmarshal(reqBody, v); err != nil {
		return errors.New(fmt.Sprintf("failed to unmarshal body of %s to type '%s' by xml tag, body='%s', error: %s",
			c.RequestString(), reflect.TypeOf(v).String(), string(reqBody), err.Error()))
	}
	return nil
}

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) Uri() string {
	return c.Request.RequestURI
}

func (c *Context) Host() string {
	return c.Request.URL.Host
}

func (c *Context) ClientIp() string {
	req := c.Request
	ipAddr := req.Header.Get("X-Real-Ip")
	if ipAddr == "" {
		ipAddr = req.Header.Get("X-Forwarded-For")
	}
	if ipAddr == "" {
		ipAddr = req.RemoteAddr
	}
	return ipAddr
}

func (c *Context) Cookies() map[string]string {
	cookies := c.Request.Cookies()
	res := map[string]string{}
	for _, cookie := range cookies {
		res[cookie.Name] = cookie.Value
	}
	return res
}

func (c *Context) Headers() map[string][]string {
	return c.Request.Header
}

// # end-region implement interface IRequest
