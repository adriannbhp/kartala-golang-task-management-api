package middleware

import (
	"net/http"
	"strings"

	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/security"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token in the Authorization header
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, response.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, response.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := security.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.ErrInvalidToken)
			c.Abort()
			return
		}

		if claims["type"] != "access" {
			response.Error(c, http.StatusUnauthorized, response.ErrInvalidToken)
			c.Abort()
			return
		}

		// Set user identification to context
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
