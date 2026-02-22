package integration

import (
	"bytes"
	"encoding/json"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/auth"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/delivery/http"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/testutil"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskIntegration(t *testing.T) {
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

	testutil.ClearTables(db, "tasks", "users")

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

	// 2. Prepare User and Token
	var accessToken string
	{
		registerParam := auth.RegisterParameter{
			Username:        "taskuser",
			Email:           "task@example.com",
			Password:        "Password123!",
		}
		registerBody, _ := json.Marshal(registerParam)
		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		router.ServeHTTP(httptest.NewRecorder(), req)

		loginParam := auth.LoginParameter{
			Email:    "task@example.com",
			Password:   "Password123!",
		}
		loginBody, _ := json.Marshal(loginParam)
		loginReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		loginReq.Header.Set("Content-Type", "application/json")
		loginReq.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, loginReq)
		
		require.Equal(t, 200, w.Code, w.Body.String())

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		
		data, ok := result["data"].(map[string]interface{})
		require.True(t, ok, "Login response data should be a map")
		require.NotNil(t, data["access_token"], "access_token should not be nil")
		accessToken = data["access_token"].(string)
	}

	var taskID string

	// 3. Test Create Task
	t.Run("CreateTask", func(t *testing.T) {
		param := tasks.CreateTaskParameter{
			Title:       "Integration Task",
			Description: "Integration Test Description",
		}
		body, _ := json.Marshal(param)

		req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, 201, resp.Code, resp.Body.String())
		
		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)
		data := result["data"].(map[string]interface{})
		require.Equal(t, "Integration Task", data["title"])
		taskID = data["id"].(string)
	})

	// 4. Test Get Task By ID
	t.Run("GetTaskByID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/tasks/"+taskID, nil)
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, 200, resp.Code, resp.Body.String())
		
		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)
		data := result["data"].(map[string]interface{})
		require.Equal(t, taskID, data["id"])
	})

	// 5. Test Update Task
	t.Run("UpdateTask", func(t *testing.T) {
		param := tasks.UpdateTaskParameter{
			Title:  "Updated Title",
			Status: "done",
		}
		body, _ := json.Marshal(param)

		req := httptest.NewRequest("PUT", "/api/v1/tasks/"+taskID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, 200, resp.Code, resp.Body.String())
		
		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)
		data := result["data"].(map[string]interface{})
		require.Equal(t, "Updated Title", data["title"])
		require.Equal(t, "done", data["status"])
	})

	// 6. Test GetAllTasks
	t.Run("GetAllTasks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/tasks?limit=10&page=1", nil)
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, 200, resp.Code, resp.Body.String())
		
		var result map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &result)
		data, ok := result["data"].([]interface{})
		require.True(t, ok, "GetAllTasks response data should be an array")
		require.Len(t, data, 1)
	})

	// 7. Test Delete Task
	t.Run("DeleteTask", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/tasks/"+taskID, nil)
		req.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, 200, resp.Code, resp.Body.String())

		// Verify deletion
		verifyReq := httptest.NewRequest("GET", "/api/v1/tasks/"+taskID, nil)
		verifyReq.Header.Set("X-API-KEY", cfg.Secret.ApiKey)
		verifyReq.Header.Set("Authorization", "Bearer "+accessToken)
		verifyResp := httptest.NewRecorder()
		router.ServeHTTP(verifyResp, verifyReq)
		assert.Equal(t, 404, verifyResp.Code)
	})
}
