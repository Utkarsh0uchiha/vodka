package mixers

import (
	"log"
	"time"

	"github.com/DevanshuTripathi/vodka"
)

// WSLogger wraps a WSHandlerFunc and logs connect, disconnect, and duration for each connection.
//
// Usage:
//
//	app.WS("/ws", mixers.WSLogger(func(c *vodka.WSContext) {
//	    // your handler
//	}))
func WSLogger(handler vodka.WSHandlerFunc) vodka.WSHandlerFunc {
	return func(c *vodka.WSContext) {
		start := time.Now()
		log.Printf(vodka.Green+"WS Connect    %s %s"+vodka.Reset, c.IP(), c.Request.URL.Path)

		handler(c)

		log.Printf(vodka.Blue+"WS Disconnect %s %s %s"+vodka.Reset, c.IP(), c.Request.URL.Path, time.Since(start))
	}
}
