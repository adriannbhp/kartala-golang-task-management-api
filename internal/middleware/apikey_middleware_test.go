package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApiKeyMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	apiKey := "test-api-key"

	t.Run("valid key", func(t *testing.T) {
		r := gin.New()
		r.Use(ApiKeyMiddleware(apiKey))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-API-KEY", apiKey)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "OK", resp.Body.String())
	})

	t.Run("missing key", func(t *testing.T) {
		r := gin.New()
		r.Use(ApiKeyMiddleware(apiKey))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("invalid key", func(t *testing.T) {
		r := gin.New()
		r.Use(ApiKeyMiddleware(apiKey))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-API-KEY", "wrong-key")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
	})
}
