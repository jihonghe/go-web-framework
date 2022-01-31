package summer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func TimeoutHandler(fn ControllerHandler, d time.Duration) ControllerHandler {
	return func(c *Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.request.WithContext(durationCtx)

		go func() {
			// 业务逻辑异常捕获
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			// 执行具体业务逻辑
			fn(c)

			// 发送业务逻辑执行完毕信号
			finish <- struct{}{}
		}()

		// 执行业务逻辑后的操作
		select {
		case p := <-panicChan:
			log.Printf("panic: %+v", p)
			c.response.WriteHeader(http.StatusInternalServerError)
		case <-finish:
			fmt.Printf("[OK] finished handling req: [%s] '%s'\n", c.request.Method, c.request.URL.Path)
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.response.Write([]byte("server time out"))
		}
		return nil
	}
}
