package tasks

import (
	"context"
	"errors"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/pagination"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	return gormDB, mock
}

func TestTaskRepository_Insert(t *testing.T) {
	db, mock := setupMockDB()
	repo := NewRepository(db)
	task := &Task{Title: "Test", UserID: uuid.New()}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Insert(context.Background(), task)
	assert.NoError(t, err)
}

func TestTaskRepository_FindAllByUserID(t *testing.T) {
	db, mock := setupMockDB()
	repo := NewRepository(db)
	userID := uuid.New()

	t.Run("success_with_filters", func(t *testing.T) {
		param := &SearchTaskParameter{
			Pagination: pagination.PaginationParameter{Page: 1, Limit: 10},
			Filter:     TaskFilterParameter{Title: "Task", Status: "todo"},
		}

		mock.ExpectQuery("count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Task A"))

		tasks, total, err := repo.FindAllByUserID(context.Background(), userID, param)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, tasks, 1)
	})

	t.Run("error_on_count", func(t *testing.T) {
		mock.ExpectQuery("count").WillReturnError(assert.AnError)
		_, _, err := repo.FindAllByUserID(context.Background(), userID, &SearchTaskParameter{})
		assert.Error(t, err)
	})

	t.Run("error_on_find", func(t *testing.T) {
		mock.ExpectQuery("count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery("SELECT").WillReturnError(assert.AnError)
		_, _, err := repo.FindAllByUserID(context.Background(), userID, &SearchTaskParameter{})
		assert.Error(t, err)
	})

	t.Run("sorting_branches", func(t *testing.T) {
		testCases := []struct {
			name     string
			sort     string
			expected string
		}{
			{"empty_sort", "", "ORDER BY created_at desc"},
			{"only_asc", "asc", "ORDER BY created_at asc"},
			{"only_desc", "desc", "ORDER BY created_at desc"},
			{"custom_sort", "title asc", "ORDER BY title asc"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				param := &SearchTaskParameter{
					Pagination: pagination.PaginationParameter{Page: 1, Limit: 10, Sort: tc.sort},
				}

				mock.ExpectQuery("count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				// We match the ORDER BY clause in the expected query
				mock.ExpectQuery(tc.expected).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(uuid.New(), "Task A"))

				_, _, err := repo.FindAllByUserID(context.Background(), userID, param)
				assert.NoError(t, err)
			})
		}
	})
}

func TestTaskRepository_FindByID(t *testing.T) {
	db, mock := setupMockDB()
	repo := NewRepository(db)
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(taskID, "Found"))

		task, err := repo.FindByID(context.Background(), taskID)
		assert.NoError(t, err)
		assert.NotNil(t, task)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
		task, err := repo.FindByID(context.Background(), taskID)
		assert.NoError(t, err)
		assert.Nil(t, task)
	})

	t.Run("other_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(assert.AnError)
		_, err := repo.FindByID(context.Background(), taskID)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Update(t *testing.T) {
	db, mock := setupMockDB()
	repo := NewRepository(db)
	task := &Task{ID: uuid.New(), Title: "Updated", UpdatedAt: time.Now()}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), task)
	assert.NoError(t, err)
}

func TestTaskRepository_Delete(t *testing.T) {
	db, mock := setupMockDB()
	repo := NewRepository(db)
	taskID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), taskID)
	assert.NoError(t, err)
}
func TestGetDBModels(t *testing.T) {
	models := GetDBModels()
	assert.Len(t, models, 1)
	assert.IsType(t, &taskDB{}, models[0])
}

func TestTaskRepository_MapError(t *testing.T) {
	db, _ := setupMockDB()
	repo := NewRepository(db).(*gormRepository)

	t.Run("nil_error", func(t *testing.T) {
		assert.NoError(t, repo.mapError(nil))
	})

	t.Run("record_not_found", func(t *testing.T) {
		err := repo.mapError(gorm.ErrRecordNotFound)
		assert.NoError(t, err)
	})

	t.Run("internal_db_error", func(t *testing.T) {
		err := repo.mapError(errors.New("relation \"tasks\" does not exist"))
		assert.Equal(t, ErrInternalDatabase, err)
	})

	t.Run("generic_error", func(t *testing.T) {
		genericErr := errors.New("generic error")
		err := repo.mapError(genericErr)
		assert.Equal(t, genericErr, err)
	})
}
