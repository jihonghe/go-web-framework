package summer

type IResponse interface {
	// 返回值设置为IResponse的好处：方便使用链式调用，链式调用能够提高代码可读性

	// region response header info: cookie, status-code, and other settings
	SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse
	SetStatus(code int) IResponse
	SetStatusOk() IResponse
	SetHeader(key, val string) IResponse
	// end-region response header info

	// region response body: json, jsonp, xml, html
	Json(v interface{}) IResponse
	Jsonp(v interface{}) IResponse
	Xml(v interface{}) IResponse
	HTML(file string, v interface{}) IResponse
	Text(format string, values ...interface{}) IResponse
	Redirect(path string) IResponse
	// end-region response body
}
