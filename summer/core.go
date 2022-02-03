package summer

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Core 框架核心结构体
type Core struct {
	router      map[string]*Trie
	middlewares *handlerPipeline
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
		middlewares: newHandlerPipeline(),
	}
}

func (c *Core) Use(handlers ...ControllerHandler) {
	c.middlewares.add(handlers)
}

// # region http method wrap

func (c *Core) Get(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares.handlers, handlers...)
	if err := c.router[http.MethodGet].AddRouter(url, allHandlers); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodGet, url, err.Error())
	}

	var handlerPipeInfo string
	for _, h := range allHandlers {
		handlerPipeInfo += runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name() + " -> "
	}
	handlerPipeInfo = strings.TrimRight(handlerPipeInfo, " -> ")
	log.Printf("register route: [%s] '%s', handler-pipeline: %s", http.MethodGet, url, handlerPipeInfo)
}

func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares.handlers, handlers...)
	if err := c.router[http.MethodPost].AddRouter(url, allHandlers); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodPost, url, err.Error())
	}
	var handlerPipeInfo string
	for _, h := range allHandlers {
		handlerPipeInfo += runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name() + " -> "
	}
	handlerPipeInfo = strings.TrimRight(handlerPipeInfo, " -> ")
	log.Printf("register route: [%s] '%s', handler-pipeline: %s", http.MethodPost, url, handlerPipeInfo)
}

func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares.handlers, handlers...)
	if err := c.router[http.MethodPut].AddRouter(url, allHandlers); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodPut, url, err.Error())
	}
	var handlerPipeInfo string
	for _, h := range allHandlers {
		handlerPipeInfo += runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name() + " -> "
	}
	handlerPipeInfo = strings.TrimRight(handlerPipeInfo, " -> ")
	log.Printf("register route: [%s] '%s', handler-pipeline: %s", http.MethodPut, url, handlerPipeInfo)
}

func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares.handlers, handlers...)
	if err := c.router[http.MethodDelete].AddRouter(url, allHandlers); err != nil {
		log.Fatalf("failed to add route '[%s] %s', error: %s", http.MethodDelete, url, err.Error())
	}
	var handlerPipeInfo string
	for _, h := range allHandlers {
		handlerPipeInfo += runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name() + " -> "
	}
	handlerPipeInfo = strings.TrimRight(handlerPipeInfo, " -> ")
	log.Printf("register route: [%s] '%s', handler-pipeline: %s", http.MethodDelete, url, handlerPipeInfo)
}

// # end-region http method wrap

func (c *Core) GetReqHandlers(r *http.Request) []ControllerHandler {
	method := strings.ToUpper(r.Method)
	uri := strings.ToLower(r.URL.Path)
	if trie, ok := c.router[method]; ok {
		handlerPipe := trie.FindHandler(uri)
		return handlerPipe.handlers
	}
	return nil
}

func (c *Core) FindRouteNode(r *http.Request) *node {
	if methodHandlers, ok := c.router[strings.ToUpper(r.Method)]; ok {
		return methodHandlers.root.matchNode(r.URL.Path)
	}
	return nil
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start executing summer.Core.ServeHTTP")
	ctx := NewContext(r, w)

	routeNode := c.FindRouteNode(r)
	if routeNode == nil {
		ctx.SetStatus(http.StatusNotFound).Json(fmt.Sprintf("route '[%s] %s' not found", r.Method, r.URL.Path))
		return
	}
	ctx.SetHandlers(routeNode.handlerPipe.handlers)
	ctx.SetParams(routeNode.parseParams(r.URL.Path))

	if err := ctx.Next(); err != nil {
		log.Printf("failed to exec handler pipeline for req %s, error: %s", ctx.RequestString(), err.Error())
		ctx.SetStatus(http.StatusInternalServerError).Json("inner error")
	}
}

func (c *Core) Group(prefix string) *Group {
	return NewGroup(c, prefix)
}
