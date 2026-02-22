package main

import (
	"context"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMainFunc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping seeder main test in short mode")
	}

	t.Run("success", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")

		assert.NotPanics(t, func() {
			main()
		})
	})

	t.Run("failure", func(t *testing.T) {
		// Mock exitFunc
		oldExit := exitFunc
		defer func() { exitFunc = oldExit }()
		
		exitCalled := false
		exitFunc = func(code int) {
			exitCalled = true
		}

		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("DB_HOST", "invalid_host")
		
		main()
		assert.True(t, exitCalled)
	})
}

type MockUserStore struct {
	FindByEmailFunc     func(ctx context.Context, email string) (*users.User, error)
	InsertNewUserFunc   func(ctx context.Context, user *users.User) (uuid.UUID, error)
	FindByUsernameFunc  func(ctx context.Context, username string) (*users.User, error)
	GetUserByUserIDFunc func(ctx context.Context, userID uuid.UUID) (*users.User, error)
}

func (m *MockUserStore) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	return m.FindByEmailFunc(ctx, email)
}
func (m *MockUserStore) InsertNewUser(ctx context.Context, user *users.User) (uuid.UUID, error) {
	return m.InsertNewUserFunc(ctx, user)
}
func (m *MockUserStore) FindByUsername(ctx context.Context, username string) (*users.User, error) {
	return m.FindByUsernameFunc(ctx, username)
}
func (m *MockUserStore) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*users.User, error) {
	return m.GetUserByUserIDFunc(ctx, userID)
}

type MockTaskStore struct {
	InsertFunc          func(ctx context.Context, task *tasks.Task) error
	FindByIDFunc        func(ctx context.Context, id uuid.UUID) (*tasks.Task, error)
	FindAllByUserIDFunc func(ctx context.Context, userID uuid.UUID, param *tasks.SearchTaskParameter) ([]tasks.Task, int64, error)
	UpdateFunc          func(ctx context.Context, task *tasks.Task) error
	DeleteFunc          func(ctx context.Context, id uuid.UUID) error
}

func (m *MockTaskStore) Insert(ctx context.Context, task *tasks.Task) error {
	return m.InsertFunc(ctx, task)
}
func (m *MockTaskStore) FindByID(ctx context.Context, id uuid.UUID) (*tasks.Task, error) {
	return m.FindByIDFunc(ctx, id)
}
func (m *MockTaskStore) FindAllByUserID(ctx context.Context, userID uuid.UUID, param *tasks.SearchTaskParameter) ([]tasks.Task, int64, error) {
	return m.FindAllByUserIDFunc(ctx, userID, param)
}
func (m *MockTaskStore) Update(ctx context.Context, task *tasks.Task) error {
	return m.UpdateFunc(ctx, task)
}
func (m *MockTaskStore) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}

func TestRunSeed(t *testing.T) {
	mockUserRepo := &MockUserStore{}
	mockTaskRepo := &MockTaskStore{}

	t.Run("success_seed", func(t *testing.T) {
		mockUserRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) {
			return nil, nil // Not found
		}
		mockUserRepo.InsertNewUserFunc = func(ctx context.Context, user *users.User) (uuid.UUID, error) {
			return uuid.New(), nil
		}
		mockTaskRepo.InsertFunc = func(ctx context.Context, task *tasks.Task) error {
			return nil
		}

		err := RunSeed(context.Background(), mockUserRepo, mockTaskRepo)
		assert.NoError(t, err)
	})

	t.Run("mixed_scenarios", func(t *testing.T) {
		count := 0
		mockUserRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) {
			count++
			if count == 1 {
				return &users.User{Email: email, ID: uuid.New()}, nil // First user exists
			}
			if count == 2 {
				return nil, assert.AnError // Second user find fails
			}
			return nil, nil // Not found
		}
		
		mockUserRepo.InsertNewUserFunc = func(ctx context.Context, user *users.User) (uuid.UUID, error) {
			if user.Email == "user@example.com" {
				return uuid.Nil, assert.AnError // Error inserting
			}
			return uuid.New(), nil
		}

		mockTaskRepo.InsertFunc = func(ctx context.Context, task *tasks.Task) error {
			if task.Title == "Setup Project" {
				return nil
			}
			return assert.AnError // Error inserting task
		}

		err := RunSeed(context.Background(), mockUserRepo, mockTaskRepo)
		assert.NoError(t, err)
	})
}

func TestBootstrapSeeder(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping bootstrap seeder in short mode")
	}

	t.Run("success", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		err := BootstrapSeeder()
		assert.NoError(t, err)
	})

	t.Run("error_config", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("JWT_SECRET", "") 
		err := BootstrapSeeder()
		assert.Error(t, err)
	})

	t.Run("error_db", func(t *testing.T) {
		t.Setenv("GO_ENV", "test")
		t.Setenv("LOAD_ENV_IN_TEST", "true")
		t.Setenv("DB_HOST", "invalid_host")
		err := BootstrapSeeder()
		assert.Error(t, err)
	})
}
