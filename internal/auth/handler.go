package auth

import (
	"errors"
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
func (h *Handler) Register(c *gin.Context) {
	var param RegisterParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	user, err := h.authUsecase.Register(c.Request.Context(), &param)
	if err != nil {
		logger.Logger.Errorf("Failed to register user: %v", err)
		if errors.Is(err, ErrEmailAlreadyExists) || errors.Is(err, ErrUsernameExists) || errors.Is(err, ErrPasswordMismatch) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.Success(c, http.StatusCreated, response.SuccessRegister, user)
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var param LoginParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	accessToken, refreshToken, user, err := h.authUsecase.Login(c.Request.Context(), &param)
	if err != nil {
		logger.Logger.Errorf("Login failed: %v", err)
		if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrInvalidPassword) {
			response.Error(c, http.StatusUnauthorized, response.ErrInvalidCredentials)
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
