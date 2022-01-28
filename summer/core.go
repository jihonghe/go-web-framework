package summer

import (
	"net/http"
)

// Core 框架核心结构体
type Core struct{}

// NewCore 初始化 Core 对象
func NewCore() *Core {
	return &Core{}
}

func (c Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: hello
}
