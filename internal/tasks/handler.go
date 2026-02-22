package tasks

import (
	"errors"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/pagination"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/logger"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

// CreateTask handles task creation
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags Tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param task body CreateTaskParameter true "Task data"
// @Success 201 {object} response.SwaggerTaskResponse
// @Failure 400 {object} response.SwaggerBadRequestResponse
// @Failure 401 {object} response.SwaggerUnauthorizedResponse
// @Failure 500 {object} response.SwaggerInternalErrorResponse
// @Router /api/v1/tasks [post]
func (h *Handler) CreateTask(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	var param CreateTaskParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	if err := param.Validate(); err != nil {
		response.ValidationError(c, http.StatusBadRequest, response.ErrInvalidRequest, err)
		return
	}

	task, err := h.usecase.CreateTask(c.Request.Context(), userID, &param)
	if err != nil {
		logger.Logger.Errorf("Failed to create task: %v", err)
		if errors.Is(err, ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.Success(c, http.StatusCreated, response.SuccessCreated, task)
}

// GetAllTasks handles getting all tasks for a user
// @Summary Get all tasks
// @Description Get a paginated list of tasks for the authenticated user with optional filtering
// @Tags Tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort query string false "Sorting field (e.g. title asc, created_at desc)"
// @Param status query string false "Filter by status (todo, in_progress, done)"
// @Param title query string false "Filter by title (ILike search)"
// @Success 200 {object} response.SwaggerTasksPaginationResponse
// @Failure 400 {object} response.SwaggerBadRequestResponse
// @Failure 401 {object} response.SwaggerUnauthorizedResponse
// @Failure 500 {object} response.SwaggerInternalErrorResponse
// @Router /api/v1/tasks [get]
func (h *Handler) GetAllTasks(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	var paginationParam pagination.PaginationParameter
	if err := c.ShouldBindQuery(&paginationParam); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	var filter TaskFilterParameter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	if err := filter.Validate(); err != nil {
		response.ValidationError(c, http.StatusBadRequest, response.ErrInvalidRequest, err)
		return
	}

	param := &SearchTaskParameter{
		Pagination: paginationParam,
		Filter:     filter,
	}

	result, err := h.usecase.GetAllTasks(c.Request.Context(), userID, param)
	if err != nil {
		logger.Logger.Errorf("Failed to get tasks: %v", err)
		if errors.Is(err, ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, response.SuccessFetch, result.Data, result.Metadata)
}

// GetTaskByID handles getting a single task
// @Summary Get task by ID
// @Description Get details of a specific task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Task ID (UUID)"
// @Success 200 {object} response.SwaggerTaskResponse
// @Failure 400 {object} response.SwaggerBadRequestResponse
// @Failure 401 {object} response.SwaggerUnauthorizedResponse
// @Failure 404 {object} response.SwaggerNotFoundResponse
// @Failure 500 {object} response.SwaggerInternalErrorResponse
// @Router /api/v1/tasks/{id} [get]
func (h *Handler) GetTaskByID(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	task, err := h.usecase.GetTaskByID(c.Request.Context(), taskID, userID)
	if err != nil {
		logger.Logger.Warnf("GetTaskByID failed: %v", err)
		if errors.Is(err, ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		if errors.Is(err, ErrTaskNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrUnauthorizedAccess) {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		response.Error(c, http.StatusNotFound, response.ErrResourceNotFound)
		return
	}

	if task == nil {
		response.Error(c, http.StatusNotFound, ErrTaskNotFound.Error())
		return
	}

	response.Success(c, http.StatusOK, response.SuccessFetch, task)
}

// UpdateTask handles task updates
// @Summary Update task
// @Description Update an existing task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Task ID (UUID)"
// @Param task body UpdateTaskParameter true "Updated task data"
// @Success 200 {object} response.SwaggerTaskResponse
// @Failure 400 {object} response.SwaggerBadRequestResponse
// @Failure 401 {object} response.SwaggerUnauthorizedResponse
// @Failure 404 {object} response.SwaggerNotFoundResponse
// @Failure 500 {object} response.SwaggerInternalErrorResponse
// @Router /api/v1/tasks/{id} [put]
func (h *Handler) UpdateTask(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	var param UpdateTaskParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	if err := param.Validate(); err != nil {
		response.ValidationError(c, http.StatusBadRequest, response.ErrInvalidRequest, err)
		return
	}

	task, err := h.usecase.UpdateTask(c.Request.Context(), taskID, userID, &param)
	if err != nil {
		logger.Logger.Errorf("Failed to update task: %v", err)
		if errors.Is(err, ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		if errors.Is(err, ErrTaskNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrUnauthorizedAccess) {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.Success(c, http.StatusOK, response.SuccessUpdate, task)
}

// DeleteTask handles task deletion
// @Summary Delete task
// @Description Delete a task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Task ID (UUID)"
// @Success 200 {object} response.SwaggerDeleteResponse
// @Failure 400 {object} response.SwaggerBadRequestResponse
// @Failure 401 {object} response.SwaggerUnauthorizedResponse
// @Failure 404 {object} response.SwaggerNotFoundResponse
// @Failure 500 {object} response.SwaggerInternalErrorResponse
// @Router /api/v1/tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidRequest)
		return
	}

	err = h.usecase.DeleteTask(c.Request.Context(), taskID, userID)
	if err != nil {
		logger.Logger.Errorf("Failed to delete task: %v", err)
		if errors.Is(err, ErrInternalDatabase) {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		if errors.Is(err, ErrTaskNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrUnauthorizedAccess) {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.ErrInternalServer)
		return
	}

	response.Success(c, http.StatusOK, response.SuccessDelete, nil)
}