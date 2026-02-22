package middleware

import (
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiKeyMiddleware validates the X-API-KEY header
func ApiKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-KEY")
		if key == "" || key != apiKey {
			response.Error(c, http.StatusForbidden, "Invalid or missing API Key")
			c.Abort()
			return
		}
		c.Next()
	}
}
