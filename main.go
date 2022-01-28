package main

import (
	"net/http"

	"go-web-framework/summer"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: summer.NewCore(),
	}

	server.ListenAndServe()
}
