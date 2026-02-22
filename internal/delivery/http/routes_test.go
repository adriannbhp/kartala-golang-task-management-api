package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	// Simply verify it doesn't panic and sets up the engine
	assert.NotPanics(t, func() {
		SetupRoutes(r, nil, nil, nil, "secret", "apikey")
	})

	// Verify some routes are registered
	routes := r.Routes()
	assert.NotEmpty(t, routes)
	
	// Check for expected routes
	hasRegister := false
	for _, route := range routes {
		if route.Path == "/api/v1/auth/register" {
			hasRegister = true
			break
		}
	}
	assert.True(t, hasRegister)
}
