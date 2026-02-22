package users

import (
	"context"

	"github.com/google/uuid"
)

// MockUserRepository is a mock implementation for testing (Mock Object Pattern)
type MockUserRepository struct {
	FindByEmailFunc     func(ctx context.Context, email string) (*User, error)
	FindByUsernameFunc  func(ctx context.Context, username string) (*User, error)
	InsertNewUserFunc   func(ctx context.Context, user *User) (uuid.UUID, error)
	GetUserByUserIDFunc func(ctx context.Context, userID uuid.UUID) (*User, error)
}

// NewMockUserRepository creates a new mock repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

// FindByEmail calls the mock function if set
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, nil
}

// FindByUsername calls the mock function if set
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	if m.FindByUsernameFunc != nil {
		return m.FindByUsernameFunc(ctx, username)
	}
	return nil, nil
}

// InsertNewUser calls the mock function if set
func (m *MockUserRepository) InsertNewUser(ctx context.Context, user *User) (uuid.UUID, error) {
	if m.InsertNewUserFunc != nil {
		return m.InsertNewUserFunc(ctx, user)
	}
	return uuid.Nil, nil
}

// GetUserByUserID calls the mock function if set
func (m *MockUserRepository) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*User, error) {
	if m.GetUserByUserIDFunc != nil {
		return m.GetUserByUserIDFunc(ctx, userID)
	}
	return nil, nil
}
