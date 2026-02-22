package users

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_GetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := NewUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.GetUserByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (*User, error) {
			return &User{ID: uid, Username: "founduser"}, nil
		}

		user, err := uc.GetUserByID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, userID, user.ID)
	})

	t.Run("not_found", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.GetUserByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (*User, error) {
			return nil, nil
		}
		user, err := uc.GetUserByID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("repo_error", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.GetUserByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (*User, error) {
			return nil, assert.AnError
		}
		_, err := uc.GetUserByID(context.Background(), userID)
		assert.Error(t, err)
	})

	t.Run("mock_fallbacks", func(t *testing.T) {
		emptyMock := NewMockUserRepository()
		ctx := context.Background()
		
		user, err := emptyMock.FindByEmail(ctx, "")
		assert.NoError(t, err)
		assert.Nil(t, user)

		user, err = emptyMock.FindByUsername(ctx, "")
		assert.NoError(t, err)
		assert.Nil(t, user)

		id, err := emptyMock.InsertNewUser(ctx, &User{})
		assert.NoError(t, err)
		assert.Equal(t, uuid.Nil, id)

		user, err = emptyMock.GetUserByUserID(ctx, uuid.Nil)
		assert.NoError(t, err)
		assert.Nil(t, user)

		// Hit the positive paths too
		emptyMock.FindByEmailFunc = func(ctx context.Context, email string) (*User, error) { return &User{}, nil }
		emptyMock.FindByUsernameFunc = func(ctx context.Context, username string) (*User, error) { return &User{}, nil }
		emptyMock.InsertNewUserFunc = func(ctx context.Context, user *User) (uuid.UUID, error) { return uuid.New(), nil }
		
		_, _ = emptyMock.FindByEmail(ctx, "")
		_, _ = emptyMock.FindByUsername(ctx, "")
		_, _ = emptyMock.InsertNewUser(ctx, &User{})
	})
}
