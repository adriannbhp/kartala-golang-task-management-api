package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()

	t.Run("GenerateAccessToken", func(t *testing.T) {
		token, err := GenerateAccessToken(userID, secret, time.Hour)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := ValidateToken(token, secret)
		assert.NoError(t, err)
		assert.Equal(t, userID.String(), claims["user_id"])
		assert.Equal(t, "access", claims["type"])
	})

	t.Run("GenerateRefreshToken", func(t *testing.T) {
		token, err := GenerateRefreshToken(userID, secret, time.Hour)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := ValidateToken(token, secret)
		assert.NoError(t, err)
		assert.Equal(t, userID.String(), claims["user_id"])
		assert.Equal(t, "refresh", claims["type"])
	})

	t.Run("ValidateToken_Errors", func(t *testing.T) {
		// Invalid string
		_, err := ValidateToken("invalid", secret)
		assert.Error(t, err)

		// Empty string
		_, err = ValidateToken("", secret)
		assert.Error(t, err)

		// Unexpected signing method (using None)
		tokenNone := jwt.New(jwt.SigningMethodNone)
		tokenNoneString, _ := tokenNone.SignedString(jwt.UnsafeAllowNoneSignatureType)
		_, err = ValidateToken(tokenNoneString, secret)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected signing method")

		// Invalid token (expired)
		tokenExpired, _ := GenerateAccessToken(userID, secret, -time.Hour)
		_, err = ValidateToken(tokenExpired, secret)
		assert.Error(t, err)

		// Invalid token (wrong secret)
		tokenWrongSecret, _ := GenerateAccessToken(userID, "wrong", time.Hour)
		_, err = ValidateToken(tokenWrongSecret, secret)
		assert.Error(t, err)
	})
}

func TestPassword(t *testing.T) {
	password := "mypassword"

	t.Run("HashAndCheck", func(t *testing.T) {
		hash, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		match, err := CheckPasswordHash(password, hash)
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = CheckPasswordHash("wrong", hash)
		assert.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("HashError", func(t *testing.T) {
		originalCost := BcryptCost
		BcryptCost = 32 // Invalid cost (> 31)
		defer func() { BcryptCost = originalCost }()

		_, err := HashPassword("test")
		assert.Error(t, err)
	})
}
