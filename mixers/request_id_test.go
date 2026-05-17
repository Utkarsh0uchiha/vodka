package mixers

import (
	"net/http/httptest"
	"testing"

	"github.com/DevanshuTripathi/vodka"
)

func TestRequestID(t *testing.T) {
	app := vodka.NewRouter()

	app.Use(RequestID())

	app.GET("/test", func(c *vodka.Context) {
		requestID, exists := c.Get("request-id")
		if !exists {
			t.Error("request-id not found in context")
		}

		if requestID == nil || requestID == "" {
			t.Error("request-id is empty")
		}

		reqHeaderID := c.Request.Header.Get("X-Request-ID")
		if reqHeaderID != requestID {
			t.Errorf("request header ID (%s) doesn't match context ID (%s)", reqHeaderID, requestID)
		}

		c.JSON(200, vodka.M{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	// Check response header
	requestIDHeader := w.Header().Get("X-Request-ID")
	if requestIDHeader == "" {
		t.Error("X-Request-ID header not found in response")
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRequestIDWithCustomHeader(t *testing.T) {
	app := vodka.NewRouter()

	customHeaderName := "X-Correlation-ID"
	app.Use(RequestIDWithHeader(customHeaderName))

	app.GET("/test", func(c *vodka.Context) {
		requestID, _ := c.Get("request-id")
		
		reqHeaderID := c.Request.Header.Get(customHeaderName)
		if reqHeaderID != requestID {
			t.Errorf("request header ID (%s) doesn't match context ID (%s)", reqHeaderID, requestID)
		}
		
		c.JSON(200, vodka.M{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	// Check custom header
	requestIDHeader := w.Header().Get(customHeaderName)
	if requestIDHeader == "" {
		t.Errorf("%s header not found in response", customHeaderName)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRequestIDUniqueness(t *testing.T) {
	app := vodka.NewRouter()

	var requestID1, requestID2 string

	app.Use(RequestID())

	app.GET("/test", func(c *vodka.Context) {
		id, _ := c.Get("request-id")
		idStr := id.(string)
		
		reqHeaderID := c.Request.Header.Get("X-Request-ID")
		if reqHeaderID != idStr {
			t.Errorf("request header ID (%s) doesn't match context ID (%s)", reqHeaderID, idStr)
		}
		
		if requestID1 == "" {
			requestID1 = idStr
		} else {
			requestID2 = idStr
		}
		c.String(200, "ok")
	})

	// First request
	req1 := httptest.NewRequest("GET", "/test", nil)
	w1 := httptest.NewRecorder()
	app.ServeHTTP(w1, req1)

	// Second request
	req2 := httptest.NewRequest("GET", "/test", nil)
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)

	// IDs should be different
	if requestID1 == requestID2 {
		t.Errorf("request IDs should be unique, but both are %s", requestID1)
	}

	// Headers should also be different
	header1 := w1.Header().Get("X-Request-ID")
	header2 := w2.Header().Get("X-Request-ID")

	if header1 == header2 {
		t.Errorf("response headers should be unique, but both are %s", header1)
	}
}

func TestRequestIDExistingHeader(t *testing.T) {
	app := vodka.NewRouter()

	app.Use(RequestID())

	app.GET("/test", func(c *vodka.Context) {
		requestID, _ := c.Get("request-id")
		
		reqHeaderID := c.Request.Header.Get("X-Request-ID")
		if reqHeaderID != requestID {
			t.Errorf("request header ID (%s) doesn't match context ID (%s)", reqHeaderID, requestID)
		}
		
		c.JSON(200, vodka.M{"request_id": requestID})
	})

	// Create request with existing X-Request-ID header (e.g., from load balancer)
	existingID := "550e8400-e29b-41d4-a716-446655440000"
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	// Should use the existing ID, not generate a new one
	responseID := w.Header().Get("X-Request-ID")
	if responseID != existingID {
		t.Errorf("expected existing ID %s, got %s", existingID, responseID)
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRequestIDFormat(t *testing.T) {
	// Test UUID format (simple check)
	id := generateUUID()

	// Should have 8-4-4-4-12 hex pattern (standard UUID format)
	if len(id) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(id))
	}

	// Check for hyphens at expected positions
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		t.Errorf("invalid UUID format: %s", id)
	}
}

func TestRequestIDAccessibleAcrossMiddleware(t *testing.T) {
	app := vodka.NewRouter()

	var requestIDFromMiddleware string

	app.Use(RequestID())

	app.Use(func(c *vodka.Context) {
		id, exists := c.Get("request-id")
		if exists {
			requestIDFromMiddleware = id.(string)
		}
		
		reqHeaderID := c.Request.Header.Get("X-Request-ID")
		if !exists || reqHeaderID != id.(string) {
			t.Errorf("request header ID (%s) doesn't match context ID (%v)", reqHeaderID, id)
		}
		
		c.Next()
	})

	app.GET("/test", func(c *vodka.Context) {
		requestID, _ := c.Get("request-id")
		c.JSON(200, vodka.M{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if requestIDFromMiddleware == "" {
		t.Error("request-id not accessible in downstream middleware")
	}

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRequestIDContextAccess(t *testing.T) {
	app := vodka.NewRouter()

	app.Use(RequestID())

	app.GET("/test", func(c *vodka.Context) {
		// Request ID should be accessible from context
		contextID, exists := c.Get("request-id")
		if !exists {
			t.Error("request-id not found in context")
		}

		responseID := c.Writer.Header().Get("X-Request-ID")

		if contextID != responseID {
			t.Errorf("context ID (%v) doesn't match response header ID (%s)", contextID, responseID)
		}

		reqHeaderID := c.Request.Header.Get("X-Request-ID")
		if contextID != reqHeaderID {
			t.Errorf("context ID (%v) doesn't match request header ID (%s)", contextID, reqHeaderID)
		}

		c.String(200, "ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
