package middleware

import (
	"log"

	"go-web-framework/summer"
)

func TestMiddleware() summer.ControllerHandler {
	return func(c *summer.Context) error {
		log.Println("middleware before TestMiddleware()")
		c.Next()
		log.Println("middleware after TestMiddleware()")
		return nil
	}
}
