package security

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "secret123"
	
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Verify
	match, err := CheckPasswordHash(password, hash)
	assert.NoError(t, err)
	assert.True(t, match)

	// Wrong password
	match, err = CheckPasswordHash("wrong", hash)
	assert.NoError(t, err)
	assert.False(t, match)
}
