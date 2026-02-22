package auth

import (
	"context"
	"fmt"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockTokenProvider struct {
	GenerateAccessTokenFunc  func(userID uuid.UUID, secret string, duration time.Duration) (string, error)
	GenerateRefreshTokenFunc func(userID uuid.UUID, secret string, duration time.Duration) (string, error)
	ValidateTokenFunc        func(token string, secret string) (map[string]interface{}, error)
}

func (m *MockTokenProvider) GenerateAccessToken(userID uuid.UUID, secret string, duration time.Duration) (string, error) {
	return m.GenerateAccessTokenFunc(userID, secret, duration)
}
func (m *MockTokenProvider) GenerateRefreshToken(userID uuid.UUID, secret string, duration time.Duration) (string, error) {
	return m.GenerateRefreshTokenFunc(userID, secret, duration)
}
func (m *MockTokenProvider) ValidateToken(token string, secret string) (map[string]interface{}, error) {
	return m.ValidateTokenFunc(token, secret)
}

type MockPasswordHasher struct {
	HashPasswordFunc      func(password string) (string, error)
	CheckPasswordHashFunc func(password, hash string) (bool, error)
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	return m.HashPasswordFunc(password)
}
func (m *MockPasswordHasher) CheckPasswordHash(password, hash string) (bool, error) {
	return m.CheckPasswordHashFunc(password, hash)
}

func TestAuthUsecase_Register(t *testing.T) {
	mockRepo := users.NewMockUserRepository()
	mockTokens := &MockTokenProvider{}
	mockHasher := &MockPasswordHasher{}
	uc := NewTestUsecase(mockRepo, "secret", mockTokens, mockHasher)

	t.Run("success", func(t *testing.T) {
		param := &RegisterParameter{
			Username:        "testuser",
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) { return nil, nil }
		mockRepo.InsertNewUserFunc = func(ctx context.Context, user *users.User) (uuid.UUID, error) { return uuid.New(), nil }
		mockHasher.HashPasswordFunc = func(p string) (string, error) { return "hashed", nil }

		user, err := uc.Register(context.Background(), param)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "hashed", user.Password)
	})

	t.Run("email_exists", func(t *testing.T) {
		param := &RegisterParameter{Email: "exists@example.com"}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) {
			return &users.User{Email: email}, nil
		}
		_, err := uc.Register(context.Background(), param)
		assert.ErrorIs(t, err, ErrEmailAlreadyExists)
	})

	t.Run("username_exists", func(t *testing.T) {
		param := &RegisterParameter{Username: "exists", Password: "p", ConfirmPassword: "p"}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) {
			return &users.User{Username: username}, nil
		}
		_, err := uc.Register(context.Background(), param)
		assert.ErrorIs(t, err, ErrUsernameExists)
	})

	t.Run("password_mismatch", func(t *testing.T) {
		param := &RegisterParameter{Password: "p1", ConfirmPassword: "p2"}
		_, err := uc.Register(context.Background(), param)
		assert.ErrorIs(t, err, ErrPasswordMismatch)
	})

	t.Run("repository_error", func(t *testing.T) {
		param := &RegisterParameter{Email: "error@e.com", Password: "p", ConfirmPassword: "p"}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) {
			return nil, fmt.Errorf("db error")
		}
		_, err := uc.Register(context.Background(), param)
		assert.Error(t, err)
	})

	t.Run("check_username_error", func(t *testing.T) {
		param := &RegisterParameter{Username: "u", Email: "e", Password: "p", ConfirmPassword: "p"}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) {
			return nil, assert.AnError
		}
		_, err := uc.Register(context.Background(), param)
		assert.Error(t, err)
	})

	t.Run("hash_error", func(t *testing.T) {
		param := &RegisterParameter{Username: "u", Email: "e", Password: "p", ConfirmPassword: "p"}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) { return nil, nil }
		mockHasher.HashPasswordFunc = func(p string) (string, error) { return "", assert.AnError }
		_, err := uc.Register(context.Background(), param)
		assert.Error(t, err)
	})

	t.Run("insert_error", func(t *testing.T) {
		param := &RegisterParameter{Username: "u", Email: "e", Password: "p", ConfirmPassword: "p"}
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) { return nil, nil }
		mockHasher.HashPasswordFunc = func(p string) (string, error) { return "h", nil }
		mockRepo.InsertNewUserFunc = func(ctx context.Context, user *users.User) (uuid.UUID, error) {
			return uuid.Nil, assert.AnError
		}
		_, err := uc.Register(context.Background(), param)
		assert.Error(t, err)
	})
}

func TestAuthUsecase_Login(t *testing.T) {
	mockRepo := users.NewMockUserRepository()
	mockTokens := &MockTokenProvider{}
	mockHasher := &MockPasswordHasher{}
	uc := NewTestUsecase(mockRepo, "secret", mockTokens, mockHasher)

	t.Run("success_by_email", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) {
			return &users.User{ID: uuid.New(), Email: email, Password: "h"}, nil
		}
		mockHasher.CheckPasswordHashFunc = func(p, h string) (bool, error) { return true, nil }
		mockTokens.GenerateAccessTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "at", nil }
		mockTokens.GenerateRefreshTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "rt", nil }

		accessToken, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "test@e.com", Password: "p"})
		assert.NoError(t, err)
		assert.Equal(t, "at", accessToken)
	})

	t.Run("success_by_username", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) {
			return &users.User{ID: uuid.New(), Username: username, Password: "h"}, nil
		}
		mockHasher.CheckPasswordHashFunc = func(p, h string) (bool, error) { return true, nil }
		mockTokens.GenerateAccessTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "at", nil }
		mockTokens.GenerateRefreshTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "rt", nil }

		_, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "user", Password: "p"})
		assert.NoError(t, err)
	})

	t.Run("remember_me", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) {
			return &users.User{ID: uuid.New()}, nil
		}
		mockHasher.CheckPasswordHashFunc = func(p, h string) (bool, error) { return true, nil }
		mockTokens.GenerateAccessTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "at", nil }
		mockTokens.GenerateRefreshTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) {
			if d == RefreshTokenDurationLong {
				return "long-rt", nil
			}
			return "short-rt", nil
		}

		_, rt, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "e", Password: "p", RememberMe: true})
		assert.NoError(t, err)
		assert.Equal(t, "long-rt", rt)
	})

	t.Run("user_not_found", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) { return nil, nil }
		_, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "n", Password: "p"})
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("invalid_password", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return &users.User{}, nil }
		mockHasher.CheckPasswordHashFunc = func(p, h string) (bool, error) { return false, nil }
		_, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "e", Password: "w"})
		assert.ErrorIs(t, err, ErrInvalidPassword)
	})

	t.Run("db_error_email", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, assert.AnError }
		_, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "e", Password: "p"})
		assert.Error(t, err)
	})

	t.Run("db_error_username", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return nil, nil }
		mockRepo.FindByUsernameFunc = func(ctx context.Context, username string) (*users.User, error) { return nil, assert.AnError }
		_, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "u", Password: "p"})
		assert.Error(t, err)
	})

	t.Run("access_token_error", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return &users.User{}, nil }
		mockHasher.CheckPasswordHashFunc = func(p, h string) (bool, error) { return true, nil }
		mockTokens.GenerateAccessTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "", assert.AnError }
		_, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "e", Password: "p"})
		assert.Error(t, err)
	})

	t.Run("refresh_token_error", func(t *testing.T) {
		mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*users.User, error) { return &users.User{}, nil }
		mockHasher.CheckPasswordHashFunc = func(p, h string) (bool, error) { return true, nil }
		mockTokens.GenerateAccessTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "at", nil }
		mockTokens.GenerateRefreshTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "", assert.AnError }
		_, _, _, err := uc.Login(context.Background(), &LoginParameter{Identifier: "e", Password: "p"})
		assert.Error(t, err)
	})
}

func TestAuthUsecase_RefreshAccessToken(t *testing.T) {
	mockRepo := users.NewMockUserRepository()
	mockTokens := &MockTokenProvider{}
	mockHasher := &MockPasswordHasher{}
	uc := NewTestUsecase(mockRepo, "secret", mockTokens, mockHasher)
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockTokens.ValidateTokenFunc = func(t string, s string) (map[string]interface{}, error) {
			return map[string]interface{}{"user_id": userID.String(), "type": TokenTypeRefresh}, nil
		}
		mockRepo.GetUserByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (*users.User, error) { return &users.User{ID: uid}, nil }
		mockTokens.GenerateAccessTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "new-at", nil }

		at, err := uc.RefreshAccessToken(context.Background(), "rt", "secret")
		assert.NoError(t, err)
		assert.Equal(t, "new-at", at)
	})

	t.Run("invalid_token", func(t *testing.T) {
		mockTokens.ValidateTokenFunc = func(t string, s string) (map[string]interface{}, error) { return nil, assert.AnError }
		_, err := uc.RefreshAccessToken(context.Background(), "rt", "secret")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("wrong_token_type", func(t *testing.T) {
		mockTokens.ValidateTokenFunc = func(t string, s string) (map[string]interface{}, error) {
			return map[string]interface{}{"type": "access"}, nil
		}
		_, err := uc.RefreshAccessToken(context.Background(), "rt", "secret")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("invalid_uuid_claim", func(t *testing.T) {
		mockTokens.ValidateTokenFunc = func(t string, s string) (map[string]interface{}, error) {
			return map[string]interface{}{"user_id": "not-uuid", "type": TokenTypeRefresh}, nil
		}
		_, err := uc.RefreshAccessToken(context.Background(), "rt", "secret")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("user_not_found_in_db", func(t *testing.T) {
		mockTokens.ValidateTokenFunc = func(t string, s string) (map[string]interface{}, error) {
			return map[string]interface{}{"user_id": userID.String(), "type": TokenTypeRefresh}, nil
		}
		mockRepo.GetUserByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (*users.User, error) { return nil, nil }
		_, err := uc.RefreshAccessToken(context.Background(), "rt", "secret")
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("db_error", func(t *testing.T) {
		mockTokens.ValidateTokenFunc = func(t string, s string) (map[string]interface{}, error) {
			return map[string]interface{}{"user_id": userID.String(), "type": TokenTypeRefresh}, nil
		}
		mockRepo.GetUserByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (*users.User, error) { return nil, assert.AnError }
		_, err := uc.RefreshAccessToken(context.Background(), "rt", "secret")
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("new_access_token_error", func(t *testing.T) {
		mockTokens.ValidateTokenFunc = func(t string, s string) (map[string]interface{}, error) {
			return map[string]interface{}{"user_id": userID.String(), "type": TokenTypeRefresh}, nil
		}
		mockRepo.GetUserByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (*users.User, error) { return &users.User{}, nil }
		mockTokens.GenerateAccessTokenFunc = func(u uuid.UUID, s string, d time.Duration) (string, error) { return "", assert.AnError }
		_, err := uc.RefreshAccessToken(context.Background(), "rt", "secret")
		assert.Error(t, err)
	})
}

func TestUsecase_Defaults(t *testing.T) {
	// This test is just to cover the default provider wrappers
	repo := users.NewMockUserRepository()
	uc := NewUsecase(repo, "secret")
	assert.NotNil(t, uc)

	tp := &defaultTokenProvider{}
	_, err := tp.GenerateAccessToken(uuid.New(), "s", time.Hour)
	assert.NoError(t, err)
	_, err = tp.GenerateRefreshToken(uuid.New(), "s", time.Hour)
	assert.NoError(t, err)
	_, err = tp.ValidateToken("invalid", "s")
	assert.Error(t, err)

	ph := &defaultPasswordHasher{}
	h, err := ph.HashPassword("p")
	assert.NoError(t, err)
	match, err := ph.CheckPasswordHash("p", h)
	assert.NoError(t, err)
	assert.True(t, match)
}
