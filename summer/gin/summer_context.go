package gin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
)

func (c *Context) BaseContext() context.Context {
	return c.Request.Context()
}

func (c *Context) RequestString() string {
	return fmt.Sprintf("request{ method: %s, url: %s}", c.Method(), c.Request.URL.Path)
}

func (c *Context) requestBody() ([]byte, error) {
	if c.Request != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to read %s body, error: %s", c.RequestString(), err.Error()))
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return body, nil
	}
	return nil, errors.New("context.request is empty")
}