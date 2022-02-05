package middleware

import (
	"net/http"

	"github.com/jihonghe/go-web-framework/summer/gin"
)

func Recovery(ctx *gin.Context) error {
	defer func() {
		if err := recover(); err != nil {
			ctx.ISetStatus(http.StatusInternalServerError).IJson(err)
		}
	}()

	ctx.Next()
	return nil
}
