package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jihonghe/go-web-framework/provider/demo"
	"github.com/jihonghe/go-web-framework/summer/gin"
)

func main() {
	engine := gin.New()
	// 绑定demo.DemoServiceProvider服务
	engine.BindSrvProvider(&demo.DemoServiceProvider{})
	registerRouter(engine)
	server := &http.Server{
		Handler: engine,
		Addr:    "localhost:8888",
	}
	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server shutdown: ", err)
	}
}
