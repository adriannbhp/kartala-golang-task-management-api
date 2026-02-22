package tasks

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockTaskStore is a local mock for testing
type MockTaskStore struct {
	InsertFunc          func(ctx context.Context, task *Task) error
	FindAllByUserIDFunc func(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) ([]Task, int64, error)
	FindByIDFunc        func(ctx context.Context, id uuid.UUID) (*Task, error)
	UpdateFunc          func(ctx context.Context, task *Task) error
	DeleteFunc          func(ctx context.Context, id uuid.UUID) error
}

func (m *MockTaskStore) Insert(ctx context.Context, task *Task) error {
	return m.InsertFunc(ctx, task)
}
func (m *MockTaskStore) FindAllByUserID(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) ([]Task, int64, error) {
	return m.FindAllByUserIDFunc(ctx, userID, param)
}
func (m *MockTaskStore) FindByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *MockTaskStore) Update(ctx context.Context, task *Task) error {
	return m.UpdateFunc(ctx, task)
}
func (m *MockTaskStore) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}

func TestTaskUsecase_CreateTask(t *testing.T) {
	mockRepo := &MockTaskStore{}
	uc := NewUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		param := &CreateTaskParameter{
			Title: "Test Task",
		}

		mockRepo.InsertFunc = func(ctx context.Context, task *Task) error {
			assert.Equal(t, userID, task.UserID)
			assert.Equal(t, TaskStatusTodo, task.Status) // Default status
			return nil
		}

		task, err := uc.CreateTask(context.Background(), userID, param)
		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, "Test Task", task.Title)
	})

	t.Run("success_with_status", func(t *testing.T) {
		userID := uuid.New()
		param := &CreateTaskParameter{
			Title:  "Test Task",
			Status: "in_progress",
		}

		mockRepo.InsertFunc = func(ctx context.Context, task *Task) error {
			assert.Equal(t, "in_progress", task.Status)
			return nil
		}

		task, err := uc.CreateTask(context.Background(), userID, param)
		assert.NoError(t, err)
		assert.Equal(t, "in_progress", task.Status)
	})

	t.Run("repository_failure", func(t *testing.T) {
		mockRepo.InsertFunc = func(ctx context.Context, task *Task) error {
			return fmt.Errorf("db error")
		}
		_, err := uc.CreateTask(context.Background(), uuid.New(), &CreateTaskParameter{Title: "fail"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestTaskUsecase_GetAllTasks(t *testing.T) {
	mockRepo := &MockTaskStore{}
	uc := NewUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.FindAllByUserIDFunc = func(ctx context.Context, uid uuid.UUID, p *SearchTaskParameter) ([]Task, int64, error) {
			assert.Equal(t, 10, p.Pagination.Limit) // Default limit
			assert.Equal(t, 1, p.Pagination.Page)   // Default page
			return []Task{{Title: "Task 1", UserID: uid}}, 1, nil
		}

		param := &SearchTaskParameter{}
		res, err := uc.GetAllTasks(context.Background(), userID, param)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		tasks := res.Data.([]Task)
		assert.Len(t, tasks, 1)
	})

	t.Run("repository_failure", func(t *testing.T) {
		mockRepo.FindAllByUserIDFunc = func(ctx context.Context, uid uuid.UUID, p *SearchTaskParameter) ([]Task, int64, error) {
			return nil, 0, fmt.Errorf("db error")
		}
		_, err := uc.GetAllTasks(context.Background(), uuid.New(), &SearchTaskParameter{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestTaskUsecase_GetTaskByID(t *testing.T) {
	mockRepo := &MockTaskStore{}
	uc := NewUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		taskID := uuid.New()
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: userID, Title: "Specific Task"}, nil
		}

		task, err := uc.GetTaskByID(context.Background(), taskID, userID)
		assert.NoError(t, err)
		assert.Equal(t, "Specific Task", task.Title)
	})

	t.Run("repo_error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return nil, assert.AnError
		}
		_, err := uc.GetTaskByID(context.Background(), uuid.New(), uuid.New())
		assert.Error(t, err)
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return nil, nil
		}
		_, err := uc.GetTaskByID(context.Background(), uuid.New(), uuid.New())
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})

	t.Run("unauthorized", func(t *testing.T) {
		userID := uuid.New()
		otherUserID := uuid.New()
		taskID := uuid.New()
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: otherUserID}, nil
		}

		_, err := uc.GetTaskByID(context.Background(), taskID, userID)
		assert.ErrorIs(t, err, ErrUnauthorizedAccess)
	})
}

func TestTaskUsecase_UpdateTask(t *testing.T) {
	mockRepo := &MockTaskStore{}
	uc := NewUsecase(mockRepo)
	userID := uuid.New()
	taskID := uuid.New()

	t.Run("success_full_update", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: userID, Title: "Old"}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, task *Task) error {
			assert.Equal(t, "New Title", task.Title)
			assert.Equal(t, "done", task.Status)
			return nil
		}

		param := &UpdateTaskParameter{Title: "New Title", Status: "done"}
		task, err := uc.UpdateTask(context.Background(), taskID, userID, param)
		assert.NoError(t, err)
		assert.Equal(t, "New Title", task.Title)
	})

	t.Run("partial_update_only_status", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: userID, Title: "Keep Title"}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, task *Task) error {
			assert.Equal(t, "Keep Title", task.Title)
			assert.Equal(t, "in_progress", task.Status)
			return nil
		}

		param := &UpdateTaskParameter{Status: "in_progress"}
		task, err := uc.UpdateTask(context.Background(), taskID, userID, param)
		assert.NoError(t, err)
		assert.Equal(t, "in_progress", task.Status)
	})

	t.Run("partial_update_only_description", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: userID, Description: "Old Desc"}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, task *Task) error {
			assert.Equal(t, "New Desc", task.Description)
			return nil
		}

		param := &UpdateTaskParameter{Description: "New Desc"}
		task, err := uc.UpdateTask(context.Background(), taskID, userID, param)
		assert.NoError(t, err)
		assert.Equal(t, "New Desc", task.Description)
	})

	t.Run("repo_find_error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return nil, assert.AnError
		}
		_, err := uc.UpdateTask(context.Background(), taskID, userID, &UpdateTaskParameter{})
		assert.Error(t, err)
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return nil, nil
		}
		_, err := uc.UpdateTask(context.Background(), taskID, userID, &UpdateTaskParameter{})
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: uuid.New()}, nil
		}
		_, err := uc.UpdateTask(context.Background(), taskID, userID, &UpdateTaskParameter{})
		assert.ErrorIs(t, err, ErrUnauthorizedAccess)
	})

	t.Run("repo_update_error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: userID}, nil
		}
		mockRepo.UpdateFunc = func(ctx context.Context, task *Task) error {
			return assert.AnError
		}
		_, err := uc.UpdateTask(context.Background(), taskID, userID, &UpdateTaskParameter{})
		assert.Error(t, err)
	})
}

func TestTaskUsecase_DeleteTask(t *testing.T) {
	mockRepo := &MockTaskStore{}
	uc := NewUsecase(mockRepo)
	userID := uuid.New()
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: userID}, nil
		}
		mockRepo.DeleteFunc = func(ctx context.Context, id uuid.UUID) error {
			return nil
		}

		err := uc.DeleteTask(context.Background(), taskID, userID)
		assert.NoError(t, err)
	})

	t.Run("repo_find_error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return nil, assert.AnError
		}
		err := uc.DeleteTask(context.Background(), taskID, userID)
		assert.Error(t, err)
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return nil, nil
		}
		err := uc.DeleteTask(context.Background(), taskID, userID)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: uuid.New()}, nil // Other user
		}
		err := uc.DeleteTask(context.Background(), taskID, userID)
		assert.ErrorIs(t, err, ErrUnauthorizedAccess)
	})

	t.Run("repo_delete_error", func(t *testing.T) {
		mockRepo.FindByIDFunc = func(ctx context.Context, id uuid.UUID) (*Task, error) {
			return &Task{ID: id, UserID: userID}, nil
		}
		mockRepo.DeleteFunc = func(ctx context.Context, id uuid.UUID) error {
			return assert.AnError
		}
		err := uc.DeleteTask(context.Background(), taskID, userID)
		assert.Error(t, err)
	})
}
