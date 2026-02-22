package response

import (
	"github.com/gin-gonic/gin"
)

// Meta defines the standardized metadata for API responses
type Meta struct {
	Message    string      `json:"message"`
	Code       int         `json:"code"`
	Status     string      `json:"status"`
	Pagination interface{} `json:"pagination,omitempty"`
}

// JSONResponse defines the standardized structure for all API responses
type JSONResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// Success sends a standardized success response
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, JSONResponse{
		Meta: Meta{
			Message: message,
			Code:    code,
			Status:  "success",
		},
		Data: data,
	})
}

// SuccessWithPagination sends a standardized success response with pagination info
func SuccessWithPagination(c *gin.Context, code int, message string, data interface{}, pagination interface{}) {
	c.JSON(code, JSONResponse{
		Meta: Meta{
			Message:    message,
			Code:       code,
			Status:     "success",
			Pagination: pagination,
		},
		Data: data,
	})
}

// Error sends a standardized error response
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, JSONResponse{
		Meta: Meta{
			Message: message,
			Code:    code,
			Status:  "error",
		},
	})
}

// AbortWithError sends a standardized error response and aborts the request
func AbortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, JSONResponse{
		Meta: Meta{
			Message: message,
			Code:    code,
			Status:  "error",
		},
	})
}

// --- Swagger Documentation Models ---
// These models are used ONLY for Swagger documentation to show accurate examples

// --- Swagger Documentation Models ---
// These models are used ONLY for Swagger documentation to show accurate, concrete examples

// SwaggerLoginResponse represents a concrete login response
type SwaggerLoginResponse struct {
	Meta struct {
		Message string `json:"message" example:"Login berhasil, selamat datang kembali"`
		Code    int    `json:"code" example:"200"`
		Status  string `json:"status" example:"success"`
	} `json:"meta"`
	Data SwaggerLoginResponseData `json:"data"`
}

// SwaggerRegisterResponse represents a concrete registration response
type SwaggerRegisterResponse struct {
	Meta struct {
		Message string `json:"message" example:"Registrasi akun berhasil"`
		Code    int    `json:"code" example:"201"`
		Status  string `json:"status" example:"success"`
	} `json:"meta"`
}

// SwaggerUserMeResponse represents a concrete user info response
type SwaggerUserMeResponse struct {
	Meta struct {
		Message string `json:"message" example:"Data berhasil dimuat"`
		Code    int    `json:"code" example:"200"`
		Status  string `json:"status" example:"success"`
	} `json:"meta"`
	Data SwaggerUserMeResponseData `json:"data"`
}

// SwaggerTaskResponse represents a concrete single task response
type SwaggerTaskResponse struct {
	Meta struct {
		Message string `json:"message" example:"Data berhasil dimuat"`
		Code    int    `json:"code" example:"200"`
		Status  string `json:"status" example:"success"`
	} `json:"meta"`
	Data SwaggerTask `json:"data"`
}

// SwaggerTasksPaginationResponse represents a concrete paginated tasks response
type SwaggerTasksPaginationResponse struct {
	Meta struct {
		Message    string `json:"message" example:"Data berhasil dimuat"`
		Code       int    `json:"code" example:"200"`
		Status     string `json:"status" example:"success"`
		Pagination struct {
			Total       int64 `json:"total" example:"100"`
			Page        int   `json:"page" example:"1"`
			Limit       int   `json:"limit" example:"10"`
			TotalPages  int   `json:"total_pages" example:"10"`
			HasNext     bool  `json:"has_next" example:"true"`
			HasPrevious bool  `json:"has_previous" example:"false"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []SwaggerTask `json:"data"`
}

// SwaggerTask represents a task structure for Swagger documentation
type SwaggerTask struct {
	ID          string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string `json:"title" example:"Belajar Swagger"`
	Description string `json:"description" example:"Implementasi Swagger di Gin"`
	Status      string `json:"status" example:"todo"`
	UserID      string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	CreatedAt   string `json:"created_at" example:"2026-02-21T15:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2026-02-21T15:00:00Z"`
}

// SwaggerDeleteResponse represents a concrete deletion response
type SwaggerDeleteResponse struct {
	Meta struct {
		Message string `json:"message" example:"Data berhasil dihapus"`
		Code    int    `json:"code" example:"200"`
		Status  string `json:"status" example:"success"`
	} `json:"meta"`
}

// SwaggerBadRequestResponse represents a 400 Bad Request response documentation
type SwaggerBadRequestResponse struct {
	Meta struct {
		Message string `json:"message" example:"Parameter permintaan tidak valid"`
		Code    int    `json:"code" example:"400"`
		Status  string `json:"status" example:"error"`
	} `json:"meta"`
}

// SwaggerUnauthorizedResponse represents a 401 Unauthorized response documentation
type SwaggerUnauthorizedResponse struct {
	Meta struct {
		Message string `json:"message" example:"Sesi tidak valid atau telah berakhir"`
		Code    int    `json:"code" example:"401"`
		Status  string `json:"status" example:"error"`
	} `json:"meta"`
}

// SwaggerNotFoundResponse represents a 404 Not Found response documentation
type SwaggerNotFoundResponse struct {
	Meta struct {
		Message string `json:"message" example:"Data yang Anda cari tidak ditemukan"`
		Code    int    `json:"code" example:"404"`
		Status  string `json:"status" example:"error"`
	} `json:"meta"`
}

// SwaggerInternalErrorResponse represents a 500 Internal Server Error response documentation
type SwaggerInternalErrorResponse struct {
	Meta struct {
		Message string `json:"message" example:"Terjadi gangguan pada server, silakan coba beberapa saat lagi"`
		Code    int    `json:"code" example:"500"`
		Status  string `json:"status" example:"error"`
	} `json:"meta"`
}

// Specific Data Models for Auth Documentation
type SwaggerLoginResponseData struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"def456..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" example:"900"`
}

type SwaggerUserMeResponseData struct {
	Name     string `json:"name" example:"John Doe"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
}
