package tasks

import (
	"context"
	"fmt"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/pagination"

	"github.com/google/uuid"
)

type Usecase interface {
	CreateTask(ctx context.Context, userID uuid.UUID, param *CreateTaskParameter) (*Task, error)
	GetAllTasks(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) (*pagination.PaginationResult, error)
	GetTaskByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Task, error)
	UpdateTask(ctx context.Context, id uuid.UUID, userID uuid.UUID, param *UpdateTaskParameter) (*Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type taskUsecase struct {
	taskRepo Repository
}

func NewUsecase(taskRepo Repository) Usecase {
	return &taskUsecase{
		taskRepo: taskRepo,
	}
}

func (uc *taskUsecase) CreateTask(ctx context.Context, userID uuid.UUID, param *CreateTaskParameter) (*Task, error) {
	task := NewTask(userID, param.Title, param.Description)
	if param.Status != "" {
		task.Status = param.Status
	}

	if err := uc.taskRepo.Insert(ctx, task); err != nil {
		return nil, fmt.Errorf("create task usecase: %w", err)
	}

	return task, nil
}

func (uc *taskUsecase) GetAllTasks(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) (*pagination.PaginationResult, error) {
	if param.Pagination.Limit <= 0 {
		param.Pagination.Limit = 10
	}
	if param.Pagination.Page <= 0 {
		param.Pagination.Page = 1
	}

	tasks, total, err := uc.taskRepo.FindAllByUserID(ctx, userID, param)
	if err != nil {
		return nil, fmt.Errorf("get all tasks usecase: %w", err)
	}

	totalPages := int((total + int64(param.Pagination.Limit) - 1) / int64(param.Pagination.Limit))

	return &pagination.PaginationResult{
		Metadata: pagination.PaginationMetadata{
			Total:       total,
			Page:        param.Pagination.Page,
			Limit:       param.Pagination.Limit,
			TotalPages:  totalPages,
			HasNext:     param.Pagination.Page < totalPages,
			HasPrevious: param.Pagination.Page > 1,
		},
		Data: tasks,
	}, nil
}

func (uc *taskUsecase) GetTaskByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Task, error) {
	task, err := uc.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get task by id usecase: %w", err)
	}

	if task == nil {
		return nil, ErrTaskNotFound
	}

	if task.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	return task, nil
}

func (uc *taskUsecase) UpdateTask(ctx context.Context, id uuid.UUID, userID uuid.UUID, param *UpdateTaskParameter) (*Task, error) {
	task, err := uc.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("update task usecase: %w", err)
	}

	if task == nil {
		return nil, ErrTaskNotFound
	}

	if task.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	if param.Title != "" {
		task.Title = param.Title
	}
	if param.Description != "" {
		task.Description = param.Description
	}
	if param.Status != "" {
		task.Status = param.Status
	}

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("update task usecase: %w", err)
	}

	return task, nil
}

func (uc *taskUsecase) DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	task, err := uc.taskRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("delete task usecase: %w", err)
	}

	if task == nil {
		return ErrTaskNotFound
	}

	if task.UserID != userID {
		return ErrUnauthorizedAccess
	}

	return uc.taskRepo.Delete(ctx, id)
}
