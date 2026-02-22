package auth

import (
	"errors"
	"regexp"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// RegisterParameter holds user registration data
type RegisterParameter struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

func (p RegisterParameter) Validate() error {
    return validation.ValidateStruct(&p,
        validation.Field(&p.Username, 
            validation.Required.Error("username is required"), 
            validation.Length(3, 50).Error("username must be between 3 and 50 characters"),
        ),
        validation.Field(&p.Email, 
            validation.Required.Error("email is required"), 
            is.Email.Error("must be a valid email address"),
        ),
        validation.Field(&p.Password, 
            validation.Required.Error("password is required"), 
            validation.Length(8, 100).Error("password must be between 8 and 100 characters"),
            // Must contain at least one uppercase letter
            validation.Match(regexp.MustCompile(`[A-Z]`)).Error("password must contain at least one uppercase letter"),
            // Must contain at least one symbol/special character
            validation.Match(regexp.MustCompile(`[^a-zA-Z0-9]`)).Error("password must contain at least one symbol"),
            // Must contain at least one number
            validation.Match(regexp.MustCompile(`[0-9]`)).Error("password must contain at least one number"),
        ),
    )
}

// LoginParameter holds user login data
type LoginParameter struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

func (p LoginParameter) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, 
            validation.Required.Error("email is required"), 
            is.Email.Error("must be a valid email address"),
        ),
		validation.Field(&p.Password, 
            validation.Required.Error("password is required"), 
            validation.Length(8, 100).Error("password must be between 8 and 100 characters"),
        ),
	)
}

// Custom errors
var (
	ErrEmailNotFound      = errors.New(response.ErrEmailNotFound)
	ErrUserNotFound       = errors.New(response.ErrUserNotFound)
	ErrEmailAlreadyExists = errors.New(response.ErrEmailExists)
	ErrUsernameExists     = errors.New(response.ErrUsernameExists)
	ErrInvalidPassword    = errors.New(response.ErrInvalidPassword)
	ErrWeakPassword       = errors.New(response.ErrWeakPassword)
	ErrInvalidToken       = errors.New(response.ErrInvalidToken)
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)
