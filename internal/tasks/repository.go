package tasks

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
	return []any{&taskDB{}}
}

type Repository interface {
	Insert(ctx context.Context, task *Task) error
	FindAllByUserID(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) ([]Task, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// taskDB represents the database model for GORM
type taskDB struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title       string    `gorm:"not null"`
	Description string    ``
	Status      string    `gorm:"default:'todo'"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (taskDB) TableName() string {
	return "tasks"
}

func (m *taskDB) toDomain() *Task {
	return &Task{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      m.Status,
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func fromDomain(t *Task) *taskDB {
	return &taskDB{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		UserID:      t.UserID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Insert(ctx context.Context, task *Task) error {
	model := fromDomain(task)
	err := r.db.WithContext(ctx).Create(model).Error
	return r.mapError(err)
}

func (r *gormRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID, param *SearchTaskParameter) ([]Task, int64, error) {
	var models []taskDB
	var total int64

	db := r.db.WithContext(ctx).Model(&taskDB{}).Where("user_id = ?", userID)

	// Filtering
	if param.Filter.Title != "" {
		db = db.Where("title ILIKE ?", "%"+param.Filter.Title+"%")
	}
	if param.Filter.Status != "" {
		db = db.Where("status = ?", param.Filter.Status)
	}

	// Count total before pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, r.mapError(err)
	}

	// Sorting
	sort := param.Pagination.Sort
	if sort == "" {
		sort = "created_at desc"
	}
	db = db.Order(sort)

	// Pagination
	offset := (param.Pagination.Page - 1) * param.Pagination.Limit
	if err := db.Limit(param.Pagination.Limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, r.mapError(err)
	}

	tasks := make([]Task, len(models))
	for i, m := range models {
		tasks[i] = *m.toDomain()
	}

	return tasks, total, nil
}

func (r *gormRepository) FindByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	var model taskDB
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, r.mapError(err)
	}
	return model.toDomain(), nil
}

func (r *gormRepository) Update(ctx context.Context, task *Task) error {
	model := fromDomain(task)
	err := r.db.WithContext(ctx).Save(model).Error
	return r.mapError(err)
}

func (r *gormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&taskDB{}, "id = ?", id).Error
	return r.mapError(err)
}

func (r *gormRepository) mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	errStr := err.Error()

	if strings.Contains(errStr, "relation") && strings.Contains(errStr, "does not exist") {
		return ErrInternalDatabase
	}

	return err
}
