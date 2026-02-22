package security

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()
	duration := 1 * time.Hour

	token, err := GenerateAccessToken(userID, secret, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
