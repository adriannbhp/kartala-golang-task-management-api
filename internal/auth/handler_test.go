package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockAuthUsecase for handler testing
type MockAuthUsecase struct {
	RegisterFunc           func(ctx context.Context, param *RegisterParameter) (*users.User, error)
	LoginFunc              func(ctx context.Context, param *LoginParameter) (string, string, *users.User, error)
	RefreshAccessTokenFunc func(ctx context.Context, refreshToken string, secret string) (string, error)
}

func (m *MockAuthUsecase) Register(ctx context.Context, param *RegisterParameter) (*users.User, error) {
	return m.RegisterFunc(ctx, param)
}
func (m *MockAuthUsecase) Login(ctx context.Context, param *LoginParameter) (string, string, *users.User, error) {
	return m.LoginFunc(ctx, param)
}
func (m *MockAuthUsecase) RefreshAccessToken(ctx context.Context, refreshToken string, secret string) (string, error) {
	return m.RefreshAccessTokenFunc(ctx, refreshToken, secret)
}

func TestHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := &MockAuthUsecase{}
	handler := NewHandler(mockUsecase, "secret")

	router := gin.Default()
	router.POST("/register", handler.Register)

	t.Run("success", func(t *testing.T) {
		param := RegisterParameter{Username: "newuser", Email: "new@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)

		mockUsecase.RegisterFunc = func(ctx context.Context, p *RegisterParameter) (*users.User, error) {
			return &users.User{}, nil
		}

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("validation_error", func(t *testing.T) {
		param := RegisterParameter{Username: "u", Email: "invalid", Password: "p"}
		body, _ := json.Marshal(param)
		
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "errors")
	})

	t.Run("conflict_email", func(t *testing.T) {
		param := RegisterParameter{Username: "exists", Email: "exists@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)
		mockUsecase.RegisterFunc = func(ctx context.Context, p *RegisterParameter) (*users.User, error) {
			return nil, ErrEmailAlreadyExists
		}
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("conflict_username", func(t *testing.T) {
		param := RegisterParameter{Username: "exists", Email: "new@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)
		mockUsecase.RegisterFunc = func(ctx context.Context, p *RegisterParameter) (*users.User, error) {
			return nil, ErrUsernameExists
		}
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("database_error", func(t *testing.T) {
		param := RegisterParameter{Username: "dbuser", Email: "db@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)
		mockUsecase.RegisterFunc = func(ctx context.Context, p *RegisterParameter) (*users.User, error) {
			return nil, users.ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("binding_error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("internal_error", func(t *testing.T) {
		param := RegisterParameter{Username: "erruser", Email: "err@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)
		mockUsecase.RegisterFunc = func(ctx context.Context, p *RegisterParameter) (*users.User, error) {
			return nil, assert.AnError
		}
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := &MockAuthUsecase{}
	handler := NewHandler(mockUsecase, "secret")

	router := gin.Default()
	router.POST("/login", handler.Login)

	t.Run("success", func(t *testing.T) {
		param := LoginParameter{Email: "test@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)
		mockUsecase.LoginFunc = func(ctx context.Context, p *LoginParameter) (string, string, *users.User, error) {
			return "at", "rt", &users.User{}, nil
		}
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("validation_error", func(t *testing.T) {
		param := LoginParameter{Email: "invalid", Password: ""}
		body, _ := json.Marshal(param)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("unauthorized_password", func(t *testing.T) {
		param := LoginParameter{Email: "test@example.com", Password: "wrongpassword"}
		body, _ := json.Marshal(param)
		mockUsecase.LoginFunc = func(ctx context.Context, p *LoginParameter) (string, string, *users.User, error) {
			return "", "", nil, ErrInvalidPassword
		}
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("user_not_found", func(t *testing.T) {
		param := LoginParameter{Email: "notfound@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)
		mockUsecase.LoginFunc = func(ctx context.Context, p *LoginParameter) (string, string, *users.User, error) {
			return "", "", nil, ErrUserNotFound
		}
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("database_error", func(t *testing.T) {
		param := LoginParameter{Email: "test@example.com", Password: "Password123!"}
		body, _ := json.Marshal(param)
		mockUsecase.LoginFunc = func(ctx context.Context, p *LoginParameter) (string, string, *users.User, error) {
			return "", "", nil, users.ErrInternalDatabase
		}
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("binding_error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("internal_error", func(t *testing.T) {
		param := LoginParameter{Email: "test@example.com", Password: "password123"}
		body, _ := json.Marshal(param)
		mockUsecase.LoginFunc = func(ctx context.Context, p *LoginParameter) (string, string, *users.User, error) {
			return "", "", nil, assert.AnError
		}
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := &MockAuthUsecase{}
	handler := NewHandler(mockUsecase, "secret")

	router := gin.Default()
	router.POST("/refresh", handler.RefreshToken)

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": "valid"}
		body, _ := json.Marshal(reqBody)
		mockUsecase.RefreshAccessTokenFunc = func(ctx context.Context, token string, secret string) (string, error) {
			return "new-at", nil
		}
		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": "expired"}
		body, _ := json.Marshal(reqBody)
		mockUsecase.RefreshAccessTokenFunc = func(ctx context.Context, token string, secret string) (string, error) {
			return "", ErrInvalidToken
		}
		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("binding_error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("internal_error", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": "valid"}
		body, _ := json.Marshal(reqBody)
		mockUsecase.RefreshAccessTokenFunc = func(ctx context.Context, token string, secret string) (string, error) {
			return "", assert.AnError
		}
		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestHandler_Ping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(nil, "")
	router := gin.Default()
	router.GET("/ping", handler.Ping)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "PONG")
}
