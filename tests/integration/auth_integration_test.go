package integration

import (
	"bytes"
	"encoding/json"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/auth"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/delivery/http"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/testutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuthIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Enable .env loading for integration tests
	t.Setenv("LOAD_ENV_IN_TEST", "true")

	// 1. Setup
	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)

	testutil.ClearTables(db, "users")

	userRepo := users.NewRepository(db)
	taskRepo := tasks.NewRepository(db)
	
	// Correct parameters for NewUsecase in auth
	authUC := auth.NewUsecase(userRepo, cfg.Secret.JwtSecret, time.Hour, time.Hour) 
	taskUC := tasks.NewUsecase(taskRepo)

	authHandler := auth.NewHandler(authUC, cfg.Secret.JwtSecret)
	userHandler := users.NewHandler(users.NewUsecase(userRepo))
	taskHandler := tasks.NewHandler(taskUC)

	router := testutil.SetupTestRouter()
	http.SetupRoutes(router, authHandler, userHandler, taskHandler, cfg.Secret.JwtSecret, cfg.Secret.ApiKey)

	// 2. Test Register
	t.Run("Register", func(t *testing.T) {
		param := auth.RegisterParameter{
			Username:        "itest",
			Email:           "itest@example.com",
			Password:        "Password123!",
		}
		body, _ := json.Marshal(param)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, 201, resp.Code, resp.Body.String())
		
		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)
		require.Equal(t, "success", result["meta"].(map[string]interface{})["status"])
	})

	// 3. Test Login
	t.Run("Login", func(t *testing.T) {
		param := auth.LoginParameter{
			Email:    "itest@example.com",
			Password:   "Password123!",
		}
		body, _ := json.Marshal(param)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, 200, resp.Code, resp.Body.String())

		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)
		
		data := result["data"].(map[string]interface{})
		require.NotNil(t, data["access_token"])
	})
}
