package users

import (
	"time"

	"github.com/google/uuid"
)

// User roles constants
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleModerator = "moderator"
)

// User entity represents a user in the system (Domain Entity)
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't expose password in JSON
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser is a constructor function to create a new User with explicit defaults
func NewUser(username, email, password string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      RoleUser, // Default role
		CreatedAt: now,
		UpdatedAt: now,
	}
}
