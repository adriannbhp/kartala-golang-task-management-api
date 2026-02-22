package tasks

import (
	"time"

	"github.com/google/uuid"
)

// Task status constants
const (
	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"
)

// Task entity represents a task in the system (Domain Entity)
type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTask is a constructor function to create a new Task with explicit defaults
func NewTask(userID uuid.UUID, title, description string) *Task {
	now := time.Now()
	return &Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      TaskStatusTodo,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
