package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {

	return func(c *Context) {
		//start timer
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
