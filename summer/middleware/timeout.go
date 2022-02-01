package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"go-web-framework/summer"
)

func Timeout(d time.Duration) summer.ControllerHandler {
	return func(c *summer.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

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
			c.Json(http.StatusInternalServerError, "inner error")
			log.Printf("failed to exec handler, error: %v", p)
		case <-finish:
			log.Printf("finished task in duration: %s", d)
		case <-durationCtx.Done():
			c.Json(http.StatusInternalServerError, "time out")
			c.SetHasTimeout()
			log.Printf("response for %s has already timeout", c.RequestString())
		}
		return nil
	}
}
