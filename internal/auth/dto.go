package auth

import (
	"errors"
)

// RegisterParameter holds user registration data
type RegisterParameter struct {
	Username        string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email           string `json:"email" binding:"required,email" example:"john@example.com"`
	Password        string `json:"password" binding:"required,min=8,max=100" example:"password123"`
}

// LoginParameter holds user login data
type LoginParameter struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// Custom errors
var (
	ErrEmailNotFound      = errors.New("email not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)
