package auth

import (
	"errors"
	"time"
)

// RegisterParameter holds user registration data
type RegisterParameter struct {
	Username        string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email           string `json:"email" binding:"required,email" example:"john@example.com"`
	Password        string `json:"password" binding:"required,min=8,max=100" example:"password123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" example:"password123"`
	Role            string `json:"role" example:"user"` // Optional, defaults to "user"
}

// LoginParameter holds user login data
type LoginParameter struct {
	Identifier string `json:"identifier" binding:"required" example:"johndoe"` // Email or Username
	Password   string `json:"password" binding:"required,min=8" example:"password123"`
	RememberMe bool   `json:"remember_me" example:"true"` // If true, token expires in 7 days; otherwise 1 hour
}

// Custom errors
var (
	ErrEmailNotFound      = errors.New("email not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrPasswordMismatch   = errors.New("password and confirm password do not match")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// Token duration constants
const (
	AccessTokenDuration = time.Minute * 15 // 15 minutes
	RefreshTokenDurationShort = time.Hour * 24 * 7  // 7 days
	RefreshTokenDurationLong  = time.Hour * 24 * 30 // 30 days
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)
