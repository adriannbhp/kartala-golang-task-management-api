package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, http.StatusOK, "Success message", gin.H{"id": 1})

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res JSONResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "success", res.Meta.Status)
	assert.Equal(t, "Success message", res.Meta.Message)
	
	data := res.Data.(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, "Error message")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var res JSONResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "error", res.Meta.Status)
	assert.Equal(t, "Error message", res.Meta.Message)
}

func TestAbortWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	AbortWithError(c, http.StatusUnauthorized, "Unauthorized")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestSuccessWithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []map[string]interface{}{{"id": 1}, {"id": 2}}
	pagination := map[string]interface{}{"total": 2, "page": 1}

	SuccessWithPagination(c, http.StatusOK, "Paginated results", data, pagination)

	assert.Equal(t, http.StatusOK, w.Code)

	var res JSONResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "success", res.Meta.Status)
	assert.Equal(t, "Paginated results", res.Meta.Message)
	assert.NotNil(t, res.Meta.Pagination)
	
	resData := res.Data.([]interface{})
	assert.Len(t, resData, 2)
}

func TestValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	validationErrors := map[string]string{
		"email": "must be a valid email address",
		"title": "is required",
	}

	ValidationError(c, http.StatusBadRequest, "Validation failed", validationErrors)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.True(t, c.IsAborted())

	var res JSONResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "error", res.Meta.Status)
	assert.Equal(t, "Validation failed", res.Meta.Message)

	resErrors := res.Errors.(map[string]interface{})
	assert.Equal(t, "must be a valid email address", resErrors["email"])
	assert.Equal(t, "is required", resErrors["title"])
}
