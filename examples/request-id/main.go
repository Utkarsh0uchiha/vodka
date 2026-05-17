package main

import (
	"log"

	"github.com/DevanshuTripathi/vodka"
	"github.com/DevanshuTripathi/vodka/mixers"
)

func main() {
	app := vodka.DefaultRouter()

	// Add request ID middleware — generates unique ID for each request
	app.Use(mixers.RequestID())

	// Example: Access request ID in handlers
	app.GET("/api/users", func(c *vodka.Context) {
		requestID, _ := c.Get("request-id")
		log.Printf("[%v] Fetching users", requestID)

		c.JSON(200, vodka.M{
			"request_id": requestID,
			"users": []vodka.M{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
		})
	})

	// Example: Custom header name
	app2 := vodka.DefaultRouter()

	// Use X-Correlation-ID instead of default X-Request-ID
	app2.Use(mixers.RequestIDWithHeader("X-Correlation-ID"))

	app2.GET("/api/orders", func(c *vodka.Context) {
		correlationID, _ := c.Get("request-id")
		log.Printf("[%v] Processing order", correlationID)

		c.JSON(200, vodka.M{
			"correlation_id": correlationID,
			"status":         "success",
		})
	})

	// Middleware can access the request ID too
	app.Use(func(c *vodka.Context) {
		requestID, _ := c.Get("request-id")
		log.Printf("Request ID %v | Method: %s | Path: %s", requestID, c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	app.GET("/", func(c *vodka.Context) {
		requestID, _ := c.Get("request-id")
		c.JSON(200, vodka.M{
			"message":    "Welcome to Vodka!",
			"request_id": requestID,
		})
	})

	log.Println("Starting server on :8080")
	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
