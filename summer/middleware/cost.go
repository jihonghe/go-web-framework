package middleware

import (
	"log"
	"time"

	"go-web-framework/summer"
)

func Cost(c *summer.Context) error {
	startTm := time.Now()
	c.Next()
	endTm := time.Now()
	log.Printf("%s costed time: %dms", c.RequestString(), endTm.Sub(startTm).Milliseconds())
	return nil
}
