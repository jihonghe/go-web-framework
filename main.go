package main

import (
	"net/http"

	"go-web-framework/summer"
)

func main() {
	core := summer.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    "localhost:8888",
	}
	server.ListenAndServe()
}
