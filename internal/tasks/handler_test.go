package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/pagination"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockTaskUsecase struct {
	CreateTaskFunc  func(ctx context.Context, userID uuid.UUID, param *CreateTaskParameter) (*Task, error)
	GetAllTasksFunc func(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) (*pagination.PaginationResult, error)
	GetTaskByIDFunc func(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Task, error)
	UpdateTaskFunc  func(ctx context.Context, id uuid.UUID, userID uuid.UUID, param *UpdateTaskParameter) (*Task, error)
	DeleteTaskFunc  func(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

func (m *MockTaskUsecase) CreateTask(ctx context.Context, userID uuid.UUID, param *CreateTaskParameter) (*Task, error) {
	return m.CreateTaskFunc(ctx, userID, param)
}
func (m *MockTaskUsecase) GetAllTasks(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) (*pagination.PaginationResult, error) {
	return m.GetAllTasksFunc(ctx, userID, param)
}
func (m *MockTaskUsecase) GetTaskByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Task, error) {
	return m.GetTaskByIDFunc(ctx, id, userID)
}
func (m *MockTaskUsecase) UpdateTask(ctx context.Context, id uuid.UUID, userID uuid.UUID, param *UpdateTaskParameter) (*Task, error) {
	return m.UpdateTaskFunc(ctx, id, userID, param)
}
func (m *MockTaskUsecase) DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return m.DeleteTaskFunc(ctx, id, userID)
}

func setupTaskRouter(mockUsecase *MockTaskUsecase, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(mockUsecase)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	
	v1 := router.Group("/api/v1/tasks")
	{
		v1.POST("", handler.CreateTask)
		v1.GET("", handler.GetAllTasks)
		v1.GET("/:id", handler.GetTaskByID)
		v1.PUT("/:id", handler.UpdateTask)
		v1.DELETE("/:id", handler.DeleteTask)
	}
	return router
}

func TestTaskHandler_CreateTask(t *testing.T) {
	mockUsecase := &MockTaskUsecase{}
	userID := uuid.New().String()
	router := setupTaskRouter(mockUsecase, userID)

	t.Run("success", func(t *testing.T) {
		param := CreateTaskParameter{Title: "New Task"}
		body, _ := json.Marshal(param)

		mockUsecase.CreateTaskFunc = func(ctx context.Context, uid uuid.UUID, p *CreateTaskParameter) (*Task, error) {
			return &Task{ID: uuid.New(), Title: p.Title}, nil
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Contains(t, resp.Body.String(), "New Task")
	})

	t.Run("invalid_json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("validation_error", func(t *testing.T) {
		param := CreateTaskParameter{Title: ""} // Empty title should fail validation
		body, _ := json.Marshal(param)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "required")
	})

	t.Run("internal_db_error", func(t *testing.T) {
		param := CreateTaskParameter{Title: "DB Error"}
		body, _ := json.Marshal(param)
		mockUsecase.CreateTaskFunc = func(ctx context.Context, uid uuid.UUID, p *CreateTaskParameter) (*Task, error) {
			return nil, ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), ErrInternalDatabase.Error())
	})

	t.Run("generic_error", func(t *testing.T) {
		param := CreateTaskParameter{Title: "Generic Error"}
		body, _ := json.Marshal(param)
		mockUsecase.CreateTaskFunc = func(ctx context.Context, uid uuid.UUID, p *CreateTaskParameter) (*Task, error) {
			return nil, errors.New("something went wrong")
		}
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	mockUsecase := &MockTaskUsecase{}
	userID := uuid.New().String()
	router := setupTaskRouter(mockUsecase, userID)

	t.Run("success", func(t *testing.T) {
		mockUsecase.GetAllTasksFunc = func(ctx context.Context, uid uuid.UUID, p *SearchTaskParameter) (*pagination.PaginationResult, error) {
			return &pagination.PaginationResult{
				Data: []Task{{Title: "Task 1"}},
				Metadata: pagination.PaginationMetadata{Total: 1},
			}, nil
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?page=1&limit=10", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid_pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?page=abc", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_filter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?id=invalid-uuid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		t.Logf("Response: %s", resp.Body.String())
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), response.ErrInvalidRequest)
	})

	t.Run("filter_validation_error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?status=invalid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("internal_db_error", func(t *testing.T) {
		mockUsecase.GetAllTasksFunc = func(ctx context.Context, uid uuid.UUID, p *SearchTaskParameter) (*pagination.PaginationResult, error) {
			return nil, ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("generic_error", func(t *testing.T) {
		mockUsecase.GetAllTasksFunc = func(ctx context.Context, uid uuid.UUID, p *SearchTaskParameter) (*pagination.PaginationResult, error) {
			return nil, errors.New("generic error")
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestTaskHandler_GetTaskByID(t *testing.T) {
	mockUsecase := &MockTaskUsecase{}
	userID := uuid.New().String()
	router := setupTaskRouter(mockUsecase, userID)
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockUsecase.GetTaskByIDFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) (*Task, error) {
			return &Task{ID: tid, Title: "Specific"}, nil
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("internal_db_error", func(t *testing.T) {
		mockUsecase.GetTaskByIDFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) (*Task, error) {
			return nil, ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		mockUsecase.GetTaskByIDFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) (*Task, error) {
			return nil, ErrTaskNotFound
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockUsecase.GetTaskByIDFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) (*Task, error) {
			return nil, ErrUnauthorizedAccess
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("task_is_nil", func(t *testing.T) {
		mockUsecase.GetTaskByIDFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) (*Task, error) {
			return nil, nil // Return nil task without error
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("other_error_mapped_to_not_found", func(t *testing.T) {
		mockUsecase.GetTaskByIDFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) (*Task, error) {
			return nil, errors.New("unexpected")
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	mockUsecase := &MockTaskUsecase{}
	userID := uuid.New().String()
	router := setupTaskRouter(mockUsecase, userID)
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		param := UpdateTaskParameter{Title: "Updated"}
		body, _ := json.Marshal(param)
		mockUsecase.UpdateTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID, p *UpdateTaskParameter) (*Task, error) {
			return &Task{ID: tid, Title: p.Title}, nil
		}
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("validation_error", func(t *testing.T) {
		param := UpdateTaskParameter{Status: "pending"} // Invalid status
		body, _ := json.Marshal(param)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("internal_db_error", func(t *testing.T) {
		param := UpdateTaskParameter{Title: "Updated"}
		body, _ := json.Marshal(param)
		mockUsecase.UpdateTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID, p *UpdateTaskParameter) (*Task, error) {
			return nil, ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		param := UpdateTaskParameter{Title: "Updated"}
		body, _ := json.Marshal(param)
		mockUsecase.UpdateTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID, p *UpdateTaskParameter) (*Task, error) {
			return nil, ErrTaskNotFound
		}
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		param := UpdateTaskParameter{Title: "Updated"}
		body, _ := json.Marshal(param)
		mockUsecase.UpdateTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID, p *UpdateTaskParameter) (*Task, error) {
			return nil, ErrUnauthorizedAccess
		}
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("generic_error", func(t *testing.T) {
		param := UpdateTaskParameter{Title: "Updated"}
		body, _ := json.Marshal(param)
		mockUsecase.UpdateTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID, p *UpdateTaskParameter) (*Task, error) {
			return nil, errors.New("generic error")
		}
		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+taskID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	mockUsecase := &MockTaskUsecase{}
	userID := uuid.New().String()
	router := setupTaskRouter(mockUsecase, userID)
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockUsecase.DeleteTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) error {
			return nil
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("internal_db_error", func(t *testing.T) {
		mockUsecase.DeleteTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) error {
			return ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		mockUsecase.DeleteTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) error {
			return ErrTaskNotFound
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockUsecase.DeleteTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) error {
			return ErrUnauthorizedAccess
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("generic_error", func(t *testing.T) {
		mockUsecase.DeleteTaskFunc = func(ctx context.Context, tid uuid.UUID, uid uuid.UUID) error {
			return errors.New("generic error")
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+taskID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
