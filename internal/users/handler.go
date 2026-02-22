package users

import (
	"errors"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	userUsecase Usecase
}

func NewHandler(userUsecase Usecase) *Handler {
	return &Handler{
		userUsecase: userUsecase,
	}
}

// GetUserInfo handles getting the current user's information
// @Summary Get current user info
// @Description Get profile information for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Success 200 {object} response.SwaggerUserMeResponse "User info retrieved successfully"
// @Failure 401 {object} response.SwaggerUnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response.SwaggerNotFoundResponse "User not found"
// @Router /api/v1/users/me [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, response.ErrUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	user, err := h.userUsecase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	if user == nil {
		response.Error(c, http.StatusNotFound, ErrUserNotFound.Error())
		return
	}

	response.Success(c, http.StatusOK, response.SuccessFetch, user)
}
