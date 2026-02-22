package auth

import (
	"context"
	"fmt"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/security"
	"time"

	"github.com/google/uuid"
)

type Usecase interface {
	Register(ctx context.Context, param *RegisterParameter) (*users.User, error)
	Login(ctx context.Context, param *LoginParameter) (string, string, *users.User, error)
	RefreshAccessToken(ctx context.Context, refreshToken string, jwtSecret string) (string, error)
}

type TokenProvider interface {
	GenerateAccessToken(userID uuid.UUID, secret string, duration time.Duration) (string, error)
	GenerateRefreshToken(userID uuid.UUID, secret string, duration time.Duration) (string, error)
	ValidateToken(token string, secret string) (map[string]interface{}, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) (bool, error)
}

type defaultTokenProvider struct{}

func (d *defaultTokenProvider) GenerateAccessToken(userID uuid.UUID, secret string, duration time.Duration) (string, error) {
	return security.GenerateAccessToken(userID, secret, duration)
}
func (d *defaultTokenProvider) GenerateRefreshToken(userID uuid.UUID, secret string, duration time.Duration) (string, error) {
	return security.GenerateRefreshToken(userID, secret, duration)
}
func (d *defaultTokenProvider) ValidateToken(token string, secret string) (map[string]interface{}, error) {
	return security.ValidateToken(token, secret)
}

type defaultPasswordHasher struct{}

func (d *defaultPasswordHasher) HashPassword(password string) (string, error) {
	return security.HashPassword(password)
}
func (d *defaultPasswordHasher) CheckPasswordHash(password, hash string) (bool, error) {
	return security.CheckPasswordHash(password, hash)
}

type authUsecase struct {
	userRepo   users.Repository
	jwtSecret  string
	tokens     TokenProvider
	passHasher PasswordHasher
}

func NewUsecase(userRepo users.Repository, jwtSecret string) Usecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		tokens:     &defaultTokenProvider{},
		passHasher: &defaultPasswordHasher{},
	}
}

// Internal version for testing
func NewTestUsecase(userRepo users.Repository, jwtSecret string, tokens TokenProvider, passHasher PasswordHasher) Usecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		tokens:     tokens,
		passHasher: passHasher,
	}
}

func (uc *authUsecase) Register(ctx context.Context, param *RegisterParameter) (*users.User, error) {
	if param.Password != param.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	existingEmail, err := uc.userRepo.FindByEmail(ctx, param.Email)
	if err != nil {
		return nil, fmt.Errorf("check email usecase: %w", err)
	}
	if existingEmail != nil {
		return nil, ErrEmailAlreadyExists
	}

	existingUsername, err := uc.userRepo.FindByUsername(ctx, param.Username)
	if err != nil {
		return nil, fmt.Errorf("check username usecase: %w", err)
	}
	if existingUsername != nil {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := uc.passHasher.HashPassword(param.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := users.NewUser(param.Username, param.Email, hashedPassword)
	if param.Role != "" {
		user.Role = param.Role
	}

	_, err = uc.userRepo.InsertNewUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}

func (uc *authUsecase) Login(ctx context.Context, param *LoginParameter) (string, string, *users.User, error) {
	// Try finding by email first
	user, err := uc.userRepo.FindByEmail(ctx, param.Identifier)
	if err != nil {
		return "", "", nil, err
	}

	// If not found by email, try by username
	if user == nil {
		user, err = uc.userRepo.FindByUsername(ctx, param.Identifier)
		if err != nil {
			return "", "", nil, err
		}
	}

	if user == nil {
		return "", "", nil, ErrUserNotFound
	}

	match, _ := uc.passHasher.CheckPasswordHash(param.Password, user.Password)
	if !match {
		return "", "", nil, ErrInvalidPassword
	}

	// Generate tokens
	accessToken, err := uc.tokens.GenerateAccessToken(user.ID, uc.jwtSecret, AccessTokenDuration)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh token duration based on remember me
	refreshDuration := RefreshTokenDurationShort
	if param.RememberMe {
		refreshDuration = RefreshTokenDurationLong
	}

	refreshToken, err := uc.tokens.GenerateRefreshToken(user.ID, uc.jwtSecret, refreshDuration)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

func (uc *authUsecase) RefreshAccessToken(ctx context.Context, refreshToken string, jwtSecret string) (string, error) {
	claims, err := uc.tokens.ValidateToken(refreshToken, jwtSecret)
	if err != nil {
		return "", ErrInvalidToken
	}

	if claims["type"] != TokenTypeRefresh {
		return "", ErrInvalidToken
	}

	userIDStr := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", ErrInvalidToken
	}

	// Check if user still exists
	user, err := uc.userRepo.GetUserByUserID(ctx, userID)
	if err != nil || user == nil {
		return "", ErrUserNotFound
	}

	// Generate new access token
	accessToken, err := uc.tokens.GenerateAccessToken(userID, jwtSecret, AccessTokenDuration)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}
