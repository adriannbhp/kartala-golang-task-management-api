package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("success_status", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestLogger(logrus.New()))
		router.GET("/success", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/success", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error_status", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestLogger(logrus.New()))
		router.GET("/error", func(c *gin.Context) {
			c.Status(http.StatusInternalServerError)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
