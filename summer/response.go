package summer

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

func (c *Context) SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.response, &http.Cookie{
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

func (c *Context) SetStatus(code int) IResponse {
	c.response.WriteHeader(code)
	return c
}

func (c *Context) SetStatusOk() IResponse {
	c.response.WriteHeader(http.StatusOK)
	return c
}

func (c *Context) SetHeader(key, val string) IResponse {
	c.response.Header().Add(key, val)
	return c
}

func (c *Context) Json(v interface{}) IResponse {
	respBody, err := json.Marshal(v)
	if err != nil {
		log.Printf("failed to marshal data(type: %s), error: %s", reflect.TypeOf(v), err.Error())
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.response.Header().Set("Content-Type", "application/json")
	_, err = c.response.Write(respBody)
	if err != nil {
		log.Printf("[Error] send response error, error: %s", err.Error())
	}
	return c
}

func (c *Context) Jsonp(v interface{}) IResponse {
	callbackFn, _ := c.QueryString("callback", "callback_function")
	c.SetHeader("Content-Type", "application/javascript")
	callback := template.JSEscapeString(callbackFn)

	_, err := c.response.Write([]byte(callback))
	if err != nil {
		log.Printf("failed to write callback-func to response of %s, callback-func='%s', error: %s",
			c.RequestString(), callbackFn, err.Error())
		return c
	}
	// 写入左括号'('
	_, err = c.response.Write([]byte("("))
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
	_, err = c.response.Write(data)
	if err != nil {
		log.Printf("failed to write data(%+v) to response of %s, error: %s",
			v, c.RequestString(), err.Error())
		return c
	}
	// 写入右括号')'
	_, err = c.response.Write([]byte(")"))
	if err != nil {
		log.Printf("failed to write '(' to response of %s, error: %s",
			c.RequestString(), err.Error())
		return c
	}

	return c
}

func (c *Context) Xml(v interface{}) IResponse {
	data, err := xml.Marshal(v)
	if err != nil {
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.SetHeader("Content-Type", "application/xml")
	c.response.Write(data)
	return c
}

func (c *Context) HTML(file string, v interface{}) IResponse {
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		log.Printf("failed to parse file '%s' to html template of %s, error: %s", file, c.RequestString(), err.Error())
		return c
	}
	if err = t.Execute(c.response, v); err != nil {
		log.Printf("failed to write output to html template of %s, error: %s", c.RequestString(), err.Error())
		return c
	}
	c.SetHeader("Content-Type", "application/json")
	return c
}

func (c *Context) Text(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	c.SetHeader("Content-Type", "application/text")
	c.response.Write([]byte(out))
	return c
}

func (c *Context) Redirect(path string) IResponse {
	http.Redirect(c.response, c.request, path, http.StatusMovedPermanently)
	return c
}
