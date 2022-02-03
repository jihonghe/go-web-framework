package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-web-framework/summer"
)

func main() {
	core := summer.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
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
