package middleware

import (
	"net/http"

	"go-web-framework/summer"
)

func Recovery(ctx *summer.Context) error {
	defer func() {
		if err := recover(); err != nil {
			ctx.Json(http.StatusInternalServerError, err)
		}
	}()

	ctx.Next()
	return nil
}
