package tasks

import (
	"errors"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
)

var (
	ErrTaskNotFound       = errors.New(response.ErrTaskNotFound)
	ErrUnauthorizedAccess = errors.New(response.ErrUnauthorizedAccess)
	ErrInternalDatabase   = errors.New(response.ErrInternalDatabase)
)
