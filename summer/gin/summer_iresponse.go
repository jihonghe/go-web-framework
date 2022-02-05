package gin

type IResponse interface {
	// 返回值设置为IResponse的好处：方便使用链式调用，链式调用能够提高代码可读性

	// region response header info: cookie, status-code, and other settings
	ISetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse
	ISetStatus(code int) IResponse
	ISetStatusOk() IResponse
	ISetHeader(key, val string) IResponse
	// end-region response header info

	// region response body: json, jsonp, xml, html
	IJson(v interface{}) IResponse
	IJsonp(v interface{}) IResponse
	IXml(v interface{}) IResponse
	IHTML(file string, v interface{}) IResponse
	IText(format string, values ...interface{}) IResponse
	IRedirect(path string) IResponse
	// end-region response body
}
