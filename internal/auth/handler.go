package auth

import (
	"errors"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/logger"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authUsecase Usecase
	jwtSecret   string
}

func NewHandler(authUsecase Usecase, jwtSecret string) *Handler {
	return &Handler{
		authUsecase: authUsecase,
		jwtSecret:   jwtSecret,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with username and email
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body RegisterParameter true "Registration details"
// @Success 201 {object} response.SwaggerRegisterResponse "User created successfully"
// @Failure 400 {object} response.SwaggerBadRequestResponse "Invalid request"
// @Failure 409 {object} response.SwaggerBadRequestResponse "Conflict (Email or Username already exists)"
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var param RegisterParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	if err := param.Validate(); err != nil {
		response.ValidationError(c, http.StatusBadRequest, response.ErrInvalidRequest, err)
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), &param)
	if err != nil {
		// Log the original error for debugging
		logger.Logger.Errorf("Failed to register user: %v", err)
		
		// Map domain errors to appropriate HTTP status and friendly messages
		if errors.Is(err, ErrEmailAlreadyExists) || errors.Is(err, ErrUsernameExists) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		
		if errors.Is(err, users.ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.Success(c, http.StatusCreated, response.SuccessRegister, user)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body LoginParameter true "Login credentials"
// @Success 200 {object} response.SwaggerLoginResponse "Login successful"
// @Failure 400 {object} response.SwaggerBadRequestResponse "Invalid request"
// @Failure 401 {object} response.SwaggerUnauthorizedResponse "Invalid credentials"
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var param LoginParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	if err := param.Validate(); err != nil {
		response.ValidationError(c, http.StatusBadRequest, response.ErrInvalidRequest, err)
		return
	}

	accessToken, refreshToken, user, err := h.authUsecase.Login(c.Request.Context(), &param)
	if err != nil {
		// Log the original error for debugging
		logger.Logger.Errorf("Login failed: %v", err)
		
		if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrInvalidPassword) {
			response.Error(c, http.StatusUnauthorized, response.ErrInvalidCredentials)
			return
		}
		
		if errors.Is(err, users.ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.Success(c, http.StatusOK, response.SuccessLogin, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken handles access token refreshing
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body object{refresh_token=string} true "Refresh token"
// @Success 200 {object} response.SwaggerLoginResponse "Token refreshed successfully"
// @Failure 400 {object} response.SwaggerBadRequestResponse "Invalid request"
// @Failure 401 {object} response.SwaggerUnauthorizedResponse "Invalid or expired token"
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	accessToken, err := h.authUsecase.RefreshAccessToken(c.Request.Context(), body.RefreshToken, h.jwtSecret)
	if err != nil {
		logger.Logger.Warnf("Token refresh failed: %v", err)
		if errors.Is(err, ErrInvalidToken) {
			response.Error(c, http.StatusUnauthorized, response.ErrInvalidToken)
			return
		}
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.Success(c, http.StatusOK, response.SuccessRefresh, gin.H{
		"access_token": accessToken,
	})
}

func (h *Handler) Ping(c *gin.Context) {
	response.Success(c, http.StatusOK, "PONG", nil)
}
