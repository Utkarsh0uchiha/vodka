package vodka

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name            string
		targetURL       string
		expectedCode    int
		expectedMessage M
	}{
		{
			name:            "Error",
			targetURL:       "/error",
			expectedCode:    404,
			expectedMessage: M{"message": "Basic Error", "success": false},
		},
		{
			name:            "No Error",
			targetURL:       "/no-error",
			expectedCode:    200,
			expectedMessage: nil,
		},
		{
			name:            "Multiple Error",
			targetURL:       "/multiple-error",
			expectedCode:    404,
			expectedMessage: M{"message": "Second", "success": false},
		},
		{
			name:            "Middleware Error",
			targetURL:       "/group/middleware-error",
			expectedCode:    500,
			expectedMessage: M{"message": "Middleware Error", "success": false},
		},
	}

	app := DefaultRouter()

	app.GET("/error", func(c *Context) {
		c.Error(404, errors.New("Basic Error"))
	})

	app.GET("/no-error", func(c *Context) {})

	app.GET("/multiple-error", func(c *Context) {
		c.Error(500, errors.New("First"))
		c.Error(404, errors.New("Second"))
	})

	gr := app.Group("/group", func(c *Context) {
		c.Error(500, errors.New("Middleware Error"))
		c.Abort()
	})

	gr.GET("/middleware-error", func(c *Context) {
		c.JSON(200, M{"success": true})
	})

	for _, tt := range tests {
		var got M

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.targetURL, nil)
			w := httptest.NewRecorder()

			app.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code=%d, got=%d", tt.expectedCode, w.Code)
			}

			json.Unmarshal(w.Body.Bytes(), &got)

			if !reflect.DeepEqual(got, tt.expectedMessage) {
				t.Errorf("Expected message=%v, got=%v", tt.expectedMessage, got)
			}
		})
	}
}
