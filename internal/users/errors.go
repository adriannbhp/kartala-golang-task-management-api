package users

import (
	"errors"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
)

var (
	ErrUserNotFound       = errors.New(response.ErrUserNotFound)
	ErrEmailAlreadyExists = errors.New(response.ErrEmailExists)
	ErrUsernameExists     = errors.New(response.ErrUsernameExists)
	ErrInternalDatabase   = errors.New(response.ErrInternalDatabase)
)
