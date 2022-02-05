package gin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

func (c *Context) ISetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     key,
		Value:    val,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return c
}

func (c *Context) ISetStatus(code int) IResponse {
	c.Writer.WriteHeader(code)
	return c
}

func (c *Context) ISetStatusOk() IResponse {
	c.Writer.WriteHeader(http.StatusOK)
	return c
}

func (c *Context) ISetHeader(key, val string) IResponse {
	c.Writer.Header().Add(key, val)
	return c
}

func (c *Context) IJson(v interface{}) IResponse {
	respBody, err := json.Marshal(v)
	if err != nil {
		log.Printf("failed to marshal data(type: %s), error: %s", reflect.TypeOf(v), err.Error())
		return c.ISetStatus(http.StatusInternalServerError)
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	_, err = c.Writer.Write(respBody)
	if err != nil {
		log.Printf("[Error] send response error, error: %s", err.Error())
	}
	return c
}

func (c *Context) IJsonp(v interface{}) IResponse {
	callbackFn, _ := c.DefaultQueryString("callback", "callback_function")
	c.ISetHeader("Content-Type", "application/javascript")
	callback := template.JSEscapeString(callbackFn)

	_, err := c.Writer.Write([]byte(callback))
	if err != nil {
		log.Printf("failed to write callback-func to response of %s, callback-func='%s', error: %s",
			c.RequestString(), callbackFn, err.Error())
		return c
	}
	// 写入左括号'('
	_, err = c.Writer.Write([]byte("("))
	if err != nil {
		log.Printf("failed to write '(' to response of %s, error: %s",
			c.RequestString(), err.Error())
		return c
	}
	// 写入数据函数参数
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("failed to marshal data(%+v) will be wrote to response of %s, error: %s",
			v, c.RequestString(), err.Error())
		return c
	}
	_, err = c.Writer.Write(data)
	if err != nil {
		log.Printf("failed to write data(%+v) to response of %s, error: %s",
			v, c.RequestString(), err.Error())
		return c
	}
	// 写入右括号')'
	_, err = c.Writer.Write([]byte(")"))
	if err != nil {
		log.Printf("failed to write '(' to response of %s, error: %s",
			c.RequestString(), err.Error())
		return c
	}

	return c
}

func (c *Context) IXml(v interface{}) IResponse {
	data, err := xml.Marshal(v)
	if err != nil {
		return c.ISetStatus(http.StatusInternalServerError)
	}
	c.ISetHeader("Content-Type", "application/xml")
	c.Writer.Write(data)
	return c
}

func (c *Context) IHTML(file string, v interface{}) IResponse {
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		log.Printf("failed to parse file '%s' to html template of %s, error: %s", file, c.RequestString(), err.Error())
		return c
	}
	if err = t.Execute(c.Writer, v); err != nil {
		log.Printf("failed to write output to html template of %s, error: %s", c.RequestString(), err.Error())
		return c
	}
	c.ISetHeader("Content-Type", "application/json")
	return c
}

func (c *Context) IText(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	c.ISetHeader("Content-Type", "application/text")
	c.Writer.Write([]byte(out))
	return c
}

func (c *Context) IRedirect(path string) IResponse {
	http.Redirect(c.Writer, c.Request, path, http.StatusMovedPermanently)
	return c
}
