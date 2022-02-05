package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jihonghe/go-web-framework/summer/gin"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			c.ISetStatus(http.StatusInternalServerError).IJson("inner error")
			log.Printf("failed to exec handler, error: %v", p)
		case <-finish:
			log.Printf("finished task in duration: %s", d)
		case <-durationCtx.Done():
			c.ISetStatus(http.StatusInternalServerError).IJson("timeout")
			cancel()
		}
	}
}
