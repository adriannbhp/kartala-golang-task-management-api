package middleware

import (
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "secret"

	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})
		return r
	}

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		token, _ := security.GenerateAccessToken(userID, secret, 1*time.Hour)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		r := setupRouter()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	})

	t.Run("missing_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)

		r := setupRouter()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid_format_short", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer") // Missing token part

		r := setupRouter()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid_format_not_bearer", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

		r := setupRouter()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid_token_signature", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.signature")

		r := setupRouter()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("wrong_token_type", func(t *testing.T) {
		userID := uuid.New()
		// Use refresh token instead of access token
		token, _ := security.GenerateRefreshToken(userID, secret, 1*time.Hour)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		r := setupRouter()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
