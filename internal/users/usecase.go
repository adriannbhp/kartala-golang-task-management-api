package users

import (
	"context"

	"github.com/google/uuid"
)

type Usecase interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
}

type userUsecase struct {
	userRepo Repository
}

func NewUsecase(userRepo Repository) Usecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uc *userUsecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	return uc.userRepo.GetUserByUserID(ctx, userID)
}
