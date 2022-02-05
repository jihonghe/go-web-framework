package middleware

import (
	"log"

	"github.com/jihonghe/go-web-framework/summer/gin"
)

func TestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("middleware before TestMiddleware()")
		c.Next()
		log.Println("middleware after TestMiddleware()")
	}
}
