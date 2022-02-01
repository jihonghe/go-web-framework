package summer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request     *http.Request
	response    http.ResponseWriter
	ctx         context.Context
	handlerPipe *handlerPipeline

	// 请求链路是否已超时
	hasTimeout bool
	// 写保护机制
	writerMux *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:     r,
		response:    w,
		ctx:         r.Context(),
		handlerPipe: newHandlerPipeline(),
		writerMux:   &sync.Mutex{},
	}
}

func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

// # region implement handler pipeline

func (c *Context) Next() error {
	return c.handlerPipe.next(c)
}

func (c *Context) SetHandlers(handlers []ControllerHandler) {
	c.handlerPipe.add(handlers)
}

// # end-region implement handler pipeline

// # region base func

func (c *Context) WriterMux() *sync.Mutex {
	return c.writerMux
}

func (c *Context) RequestString() string {
	return fmt.Sprintf("request:{method: %s, url-path: %s}", c.request.Method, c.request.URL.Path)
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}

func (c *Context) GetResponse() http.ResponseWriter {
	return c.response
}

func (c *Context) SetHasTimeout() {
	c.hasTimeout = true
}

func (c *Context) HasTimeout() bool {
	return c.hasTimeout
}

// # end-region base func

// # region implement interface context.Context

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.BaseContext().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.BaseContext().Done()
}

func (c *Context) Err() error {
	return c.BaseContext().Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.BaseContext().Value(key)
}

// # end-region implement interface context.Context

// # region implement request

func (c *Context) QueryParams() map[string][]string {
	if c.request != nil {
		return c.request.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) getQueryParamValue(key string) string {
	params := c.QueryParams()
	values, ok := params[key]
	if !ok || len(values) <= 0 {
		return ""
	}
	return values[len(values)-1]
}

func (c *Context) getQueryParamValues(key string) []string {
	return c.QueryParams()[key]
}

func (c *Context) QueryInt(key string, def int) int {
	paramVal := c.getQueryParamValue(key)
	if paramVal == "" {
		return def
	}

	val, err := strconv.Atoi(paramVal)
	if err != nil {
		return def
	}

	return val
}

func (c *Context) QueryString(key string, def string) string {
	paramVal := c.getQueryParamValue(key)
	if paramVal == "" {
		return def
	}
	return paramVal
}

// QueryBool returns the boolean value represented by the string.
// The value it accepts is the same as what strconv.ParseBool accepts.
func (c *Context) QueryBool(key string, def bool) bool {
	paramVal := c.getQueryParamValue(key)
	if paramVal == "" {
		return def
	}

	val, err := strconv.ParseBool(paramVal)
	if err != nil {
		return def
	}
	return val
}

func (c *Context) QueryStrSlc(key string, def []string) []string {
	paramValues := c.getQueryParamValues(key)
	if len(paramValues) == 0 {
		return def
	}
	return paramValues
}

// # end-region implement request

// # region form post

func (c *Context) FormParams() map[string][]string {
	if c.request != nil {
		return c.request.PostForm
	}
	return map[string][]string{}
}

func (c *Context) getFormParamValue(key string) string {
	params := c.FormParams()
	values, ok := params[key]
	if !ok || len(values) <= 0 {
		return ""
	}
	return values[len(values)-1]
}

func (c *Context) getFormParamValues(key string) []string {
	return c.FormParams()[key]
}

func (c *Context) FormString(key string, def string) string {
	formValue := c.getFormParamValue(key)
	if formValue == "" {
		return def
	}
	return formValue
}

func (c *Context) FormInt(key string, def int) int {
	formValue := c.getFormParamValue(key)
	if formValue == "" {
		return def
	}

	val, err := strconv.Atoi(formValue)
	if err != nil {
		return def
	}
	return val
}

func (c *Context) FormBool(key string, def bool) bool {
	formValue := c.getFormParamValue(key)
	if formValue == "" {
		return def
	}

	val, err := strconv.ParseBool(formValue)
	if err != nil {
		return def
	}
	return val
}

func (c *Context) FormStrSlc(key string, def []string) []string {
	paramValues := c.getFormParamValues(key)
	if len(paramValues) == 0 {
		return def
	}
	return paramValues
}

func (c *Context) UnmarshalReqForm(v interface{}) error {
	if c.request == nil {
		return errors.New("request body is empty")
	}

	reqBody, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body, error: %s", err.Error())
	}

	err = json.Unmarshal(reqBody, v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal request body to type %s, error: %s", reflect.TypeOf(v).String(), err.Error())
	}

	return nil
}

// # end-region form post

// # region response

func (c *Context) Json(status int, data interface{}) error {
	c.WriterMux().Lock()
	defer c.WriterMux().Unlock()

	if c.HasTimeout() {
		return nil
	}
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status)

	respBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data(type: %s), error: %s", reflect.TypeOf(data), err.Error())
	}
	c.response.Write(respBody)
	return nil
}

func (c *Context) HTML(status int, data interface{}, template string) error {
	return nil
}

func (c *Context) Text(status int, data string) error {
	return nil
}

// # end-region response
