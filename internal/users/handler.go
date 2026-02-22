package users

import (
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
func (h *Handler) GetUserInfo(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to parse user ID: "+err.Error())
		return
	}

	user, err := h.userUsecase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get user info: "+err.Error())
		return
	}

	if user == nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, "User info retrieved successfully", user)
}
