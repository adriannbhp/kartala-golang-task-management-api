package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockUserUsecase for handler testing
type MockUserUsecase struct {
	GetUserByIDFunc func(ctx context.Context, userID uuid.UUID) (*User, error)
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	return m.GetUserByIDFunc(ctx, userID)
}

func TestHandler_GetUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := &MockUserUsecase{}
	handler := NewHandler(mockUsecase)

	router := gin.Default()
	router.GET("/me", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		handler.GetUserInfo(c)
	})

	t.Run("success", func(t *testing.T) {
		mockUsecase.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (*User, error) {
			return &User{Email: "me@example.com"}, nil
		}
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "me@example.com")
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockUsecase.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (*User, error) {
			return nil, assert.AnError
		}
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("database_error", func(t *testing.T) {
		mockUsecase.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (*User, error) {
			return nil, ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), ErrInternalDatabase.Error())
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockUsecase.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (*User, error) {
			return nil, nil
		}
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		routerNoID := gin.New()
		routerNoID.GET("/me", handler.GetUserInfo)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp := httptest.NewRecorder()
		routerNoID.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("invalid_user_id", func(t *testing.T) {
		routerInvalidID := gin.New()
		routerInvalidID.GET("/me", func(c *gin.Context) {
			c.Set("user_id", "invalid-uuid")
			handler.GetUserInfo(c)
		})

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp := httptest.NewRecorder()
		routerInvalidID.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
