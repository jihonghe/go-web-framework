package summer

import (
	"log"
	"net/http"
)

// Core 框架核心结构体
type Core struct {
	router map[string]ControllerHandler
}

// NewCore 初始化 Core 对象
func NewCore() *Core {
	return &Core{
		router: make(map[string]ControllerHandler),
	}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start executing summer.Core.ServeHTTP")
	ctx := NewContext(r, w)

	routerHandler := c.router["foo"]
	if routerHandler == nil {
		return
	}
	log.Println("matched routerHandler for url 'foo'")

	routerHandler(ctx)
}
