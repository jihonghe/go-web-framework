package summer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request     *http.Request
	response    http.ResponseWriter
	ctx         context.Context
	handlerPipe *handlerPipeline

	params     map[string]string // 路由匹配参数
	hasTimeout bool              // 请求链路是否已超时
	writerMux  *sync.Mutex       // 写保护机制
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

func (c *Context) requestBody() ([]byte, error) {
	if c.request != nil {
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to read %s body, error: %s", c.RequestString(), err.Error()))
		}
		c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return body, nil
	}
	return nil, errors.New("context.request is empty")
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

func (c *Context) SetParams(params map[string]string) {
	c.params = params
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
