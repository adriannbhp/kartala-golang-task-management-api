package users

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		username := "testuser"
		email := "test@example.com"
		password := "hashed_password"

		u := NewUser(username, email, password)

		assert.NotEqual(t, uuid.Nil, u.ID)
		assert.Equal(t, username, u.Username)
		assert.Equal(t, email, u.Email)
		assert.Equal(t, password, u.Password)
		assert.False(t, u.CreatedAt.IsZero())
		assert.False(t, u.UpdatedAt.IsZero())
		assert.Equal(t, u.CreatedAt, u.UpdatedAt)
	})
}
