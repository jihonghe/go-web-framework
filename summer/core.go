package summer

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Core 框架核心结构体
type Core struct {
	router map[string]*Trie
}

// NewCore 初始化 Core 对象
func NewCore() *Core {
	return &Core{
		router: map[string]*Trie{
			http.MethodGet:    NewTrie(),
			http.MethodPost:   NewTrie(),
			http.MethodPut:    NewTrie(),
			http.MethodDelete: NewTrie(),
		},
	}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	if err := c.router[http.MethodGet].AddRouter(url, handler); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodGet, url, err.Error())
	}
	log.Printf("register route: [%s] '%s'", http.MethodGet, url)
}

func (c *Core) Post(url string, handler ControllerHandler) {
	if err := c.router[http.MethodPost].AddRouter(url, handler); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodPost, url, err.Error())
	}
	log.Printf("register route: [%s] '%s'", http.MethodPost, url)
}

func (c *Core) Put(url string, handler ControllerHandler) {
	if err := c.router[http.MethodPut].AddRouter(url, handler); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodPut, url, err.Error())
	}
	log.Printf("register route: [%s] '%s'", http.MethodPut, url)
}

func (c *Core) Delete(url string, handler ControllerHandler) {
	if err := c.router[http.MethodDelete].AddRouter(url, handler); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodDelete, url, err.Error())
	}
	log.Printf("register route: [%s] '%s'", http.MethodDelete, url)
}

func (c *Core) FindRouteByRequest(r *http.Request) ControllerHandler {
	method := strings.ToUpper(r.Method)
	uri := strings.ToLower(r.URL.Path)
	if trie, ok := c.router[method]; ok {
		handler := trie.FindHandler(uri)
		return handler
	}
	return nil
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start executing summer.Core.ServeHTTP")
	ctx := NewContext(r, w)

	routerHandler := c.FindRouteByRequest(r)
	if routerHandler == nil {
		ctx.Json(http.StatusNotFound, fmt.Sprintf("route '[%s] %s' not found", r.Method, r.URL.Path))
		return
	}

	if err := routerHandler(ctx); err != nil {
		ctx.Json(http.StatusInternalServerError, "inner error")
		return
	}
}

func (c *Core) Group(prefix string) *Group {
	return NewGroup(c, prefix)
}
