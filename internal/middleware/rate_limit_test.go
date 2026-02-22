package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allow_requests_within_limit", func(t *testing.T) {
		r := gin.New()
		// 10 requests per second, burst of 10
		r.Use(RateLimiter(rate.Limit(10), 10))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("block_requests_exceeding_limit", func(t *testing.T) {
		r := gin.New()
		// 1 request per second, burst of 1
		r.Use(RateLimiter(rate.Limit(1), 1))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		// First request allowed
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second request immediately after should be blocked
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
		assert.Contains(t, w2.Body.String(), "Too many requests")
	})

	t.Run("allow_after_wait", func(t *testing.T) {
		r := gin.New()
		// 10 requests per second (1 every 100ms), burst of 1
		r.Use(RateLimiter(rate.Limit(10), 1))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		// First request
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Wait for 150ms to allow another request
		time.Sleep(150 * time.Millisecond)

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)
	})
}
