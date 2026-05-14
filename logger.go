package vodka

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		log.Printf(Blue+"%s %s %s"+Reset, c.Request.Method, c.Request.URL.Path, time.Since(start))
	}
}
