package middleware

import (
	"log"
	"time"

	"github.com/jihonghe/go-web-framework/summer/gin"
)

func Cost(c *gin.Context) {
	startTm := time.Now()
	c.Next()
	endTm := time.Now()
	log.Printf("%s costed time: %dms", c.RequestString(), endTm.Sub(startTm).Milliseconds())
}
