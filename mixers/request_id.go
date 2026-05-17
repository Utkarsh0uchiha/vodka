package mixers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/DevanshuTripathi/vodka"
)

const defaultRequestIDKey = "request-id"
const defaultRequestIDHeader = "X-Request-ID"

// generateUUID creates a simple UUID v4 format string
func generateUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "error-" + hex.EncodeToString(make([]byte, 4))
	}

	// Set version (4) and variant bits
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// RequestID returns a middleware that generates and tracks request IDs
// Uses default header "X-Request-ID" and context key "request-id"
func RequestID() vodka.HandlerFunc {
	return RequestIDWithHeader(defaultRequestIDHeader)
}

// RequestIDWithHeader returns a middleware that generates and tracks request IDs
// with a custom header name. Context key is always "request-id"
// If the header already exists (e.g., from a load balancer), it uses the existing value.
func RequestIDWithHeader(headerName string) vodka.HandlerFunc {
	return func(c *vodka.Context) {
		// Check if request ID already exists (from upstream proxy/load balancer)
		requestID := c.Request.Header.Get(headerName)

		// If not present, generate a new one
		if requestID == "" {
			requestID = generateUUID()
		}

		// Store in context for access by handlers
		c.Set(defaultRequestIDKey, requestID)

		// Add to response header
		c.Writer.Header().Set(headerName, requestID)

		c.Next()
	}
}
