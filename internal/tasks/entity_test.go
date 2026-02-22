package tasks

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		title := "Test Task"
		description := "Test Description"

		task := NewTask(userID, title, description)

		assert.NotEqual(t, uuid.Nil, task.ID)
		assert.Equal(t, userID, task.UserID)
		assert.Equal(t, title, task.Title)
		assert.Equal(t, description, task.Description)
		assert.Equal(t, TaskStatusTodo, task.Status)
		assert.False(t, task.CreatedAt.IsZero())
		assert.False(t, task.UpdatedAt.IsZero())
		assert.Equal(t, task.CreatedAt, task.UpdatedAt)
	})
}
