package summer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/spf13/cast"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

func (c *Context) QueryParams() map[string][]string {
	if c.request != nil {
		return c.request.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) FormParams() map[string][]string {
	if c.request != nil {
		c.request.ParseForm()
		return c.request.PostForm
	}
	return map[string][]string{}
}

// # region implement interface IRequest

func (c *Context) Query(key string) (interface{}, bool) {
	values, ok := c.QueryParams()[key]
	if ok && len(values) > 0 {
		return values[0], true
	}
	return nil, false
}

func (c *Context) QueryInt(key string, def int) (int, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToInt(values[0]), true
	}

	return def, false
}

func (c *Context) QueryInt64(key string, def int64) (int64, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToInt64(values[0]), true
	}
	return def, false
}

func (c *Context) QueryFloat32(key string, def float32) (float32, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat32(values[0]), true
	}
	return def, false
}

func (c *Context) QueryFloat64(key string, def float64) (float64, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat64(values[0]), true
	}
	return def, false
}

func (c *Context) QueryBool(key string, def bool) (bool, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToBool(values[0]), true
	}
	return def, false
}

func (c *Context) QueryString(key string, def string) (string, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToString(values[0]), true
	}
	return def, false
}

func (c *Context) QueryStrSlice(key string, def []string) ([]string, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return values, true
	}
	return def, false
}

func (c *Context) QueryIntSlice(key string, def []int) ([]int, bool) {
	if values, ok := c.QueryParams()[key]; ok && len(values) > 0 {
		return cast.ToIntSlice(values), true
	}
	return def, false
}

func (c *Context) Param(key string) (interface{}, bool) {
	if c.params != nil {
		if val, ok := c.params[key]; ok {
			return val, true
		}
	}
	return nil, false
}

func (c *Context) ParamInt(key string, def int) (int, bool) {
	val, ok := c.Param(key)
	if ok {
		return cast.ToInt(val), true
	}
	return def, false
}

func (c *Context) ParamInt64(key string, def int64) (int64, bool) {
	val, ok := c.Param(key)
	if ok {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (c *Context) ParamFloat32(key string, def float32) (float32, bool) {
	val, ok := c.Param(key)
	if ok {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (c *Context) ParamFloat64(key string, def float64) (float64, bool) {
	val, ok := c.Param(key)
	if ok {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (c *Context) ParamBool(key string, def bool) (bool, bool) {
	val, ok := c.Param(key)
	if ok {
		return cast.ToBool(val), true
	}
	return def, false
}

func (c *Context) ParamString(key string, def string) (string, bool) {
	val, ok := c.Param(key)
	if ok {
		return cast.ToString(val), true
	}
	return def, false
}

func (c *Context) Form(key string) (interface{}, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return nil, false
}

func (c *Context) FormInt(key string, def int) (int, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToInt(values[0]), true
	}
	return def, false
}

func (c *Context) FormInt64(key string, def int64) (int64, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToInt64(values[0]), true
	}
	return def, false
}

func (c *Context) FormFloat32(key string, def float32) (float32, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat32(values[0]), true
	}
	return def, false
}

func (c *Context) FormFloat64(key string, def float64) (float64, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToFloat64(values[0]), true
	}
	return def, false
}

func (c *Context) FormBool(key string, def bool) (bool, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToBool(values[0]), true
	}
	return def, false
}

func (c *Context) FormString(key string, def string) (string, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToString(values[0]), true
	}
	return def, false
}

func (c *Context) FormStrSlice(key string, def []string) ([]string, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToStringSlice(values), true
	}
	return def, false
}

func (c *Context) FormIntSlice(key string, def []int) ([]int, bool) {
	if values, ok := c.FormParams()[key]; ok && len(values) > 0 {
		return cast.ToIntSlice(values), true
	}
	return def, false
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if c.request.MultipartForm == nil {
		return nil, c.request.ParseMultipartForm(defaultMultipartMemory)
	}
	f, fh, err := c.request.FormFile(key)
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

func (c *Context) GetRawData() ([]byte, error) {
	return c.requestBody()
}

func (c *Context) Method() string {
	return c.request.Method
}

func (c *Context) Uri() string {
	return c.request.RequestURI
}

func (c *Context) Host() string {
	return c.request.URL.Host
}

func (c *Context) ClientIp() string {
	req := c.request
	ipAddr := req.Header.Get("X-Real-Ip")
	if ipAddr == "" {
		ipAddr = req.Header.Get("X-Forwarded-For")
	}
	if ipAddr == "" {
		ipAddr = req.RemoteAddr
	}
	return ipAddr
}

func (c *Context) Cookie(key string) (string, bool) {
	if val, ok := c.Cookies()[key]; ok {
		return val, true
	}
	return "", false
}

func (c *Context) Cookies() map[string]string {
	cookies := c.request.Cookies()
	res := map[string]string{}
	for _, cookie := range cookies {
		res[cookie.Name] = cookie.Value
	}
	return res
}

func (c *Context) Header(key string) (string, bool) {
	if values, ok := c.Headers()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (c *Context) Headers() map[string][]string {
	return c.request.Header
}

// # end-region implement interface IRequest
