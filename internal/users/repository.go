package users

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetDBModels returns the database models for GORM migrations
func GetDBModels() []any {
	return []any{&userDB{}}
}

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	InsertNewUser(ctx context.Context, user *User) (uuid.UUID, error)
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*User, error)
}

type gormRepository struct {
	db *gorm.DB
}

// userDB represents the database model for GORM
type userDB struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex"`
	Password  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (userDB) TableName() string {
	return "users"
}

func (m *userDB) toDomain() *User {
	return &User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func fromDomain(u *User) *userDB {
	return &userDB{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var model userDB
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, r.mapError(err)
	}
	return model.toDomain(), nil
}

func (r *gormRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var model userDB
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error
	if err != nil {
		return nil, r.mapError(err)
	}
	return model.toDomain(), nil
}

func (r *gormRepository) InsertNewUser(ctx context.Context, user *User) (uuid.UUID, error) {
	model := fromDomain(user)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return uuid.Nil, r.mapError(err)
	}
	return model.ID, nil
}

func (r *gormRepository) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var model userDB
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&model).Error
	if err != nil {
		return nil, r.mapError(err)
	}
	return model.toDomain(), nil
}

func (r *gormRepository) mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	errStr := err.Error()

	// Map generic SQL errors or postgres specific errors to domain errors
	if strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "duplicate key") {
		if strings.Contains(errStr, "email") {
			return ErrEmailAlreadyExists
		}
		if strings.Contains(errStr, "username") {
			return ErrUsernameExists
		}
	}

	if strings.Contains(errStr, "relation") && strings.Contains(errStr, "does not exist") {
		return ErrInternalDatabase
	}

	return err
}
