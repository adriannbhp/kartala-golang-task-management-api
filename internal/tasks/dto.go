package tasks

import (
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/pagination"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

// CreateTaskParameter holds task creation data
type CreateTaskParameter struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (p CreateTaskParameter) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, 
			validation.Required.Error("title is required"), 
			validation.Length(3, 100).Error("title must be between 3 and 100 characters"),
		),
		validation.Field(&p.Status, validation.In("todo", "in_progress", "done").Error("must be one of: todo, in_progress, done")),
	)
}

// UpdateTaskParameter holds task update data
type UpdateTaskParameter struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (p UpdateTaskParameter) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Length(3, 100)),
		validation.Field(&p.Status, validation.In("todo", "in_progress", "done").Error("must be one of: todo, in_progress, done")),
	)
}

// TaskFilterParameter holds task-specific filtering data
type TaskFilterParameter struct {
	ID     uuid.UUID `form:"id"`
	Status string    `form:"status"`
	Title  string    `form:"title"`
}

func (p TaskFilterParameter) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Status, validation.In("todo", "in_progress", "done").Error("must be one of: todo, in_progress, done")),
	)
}

// SearchTaskParameter combines pagination and filtering
type SearchTaskParameter struct {
	Pagination pagination.PaginationParameter
	Filter     TaskFilterParameter
}
