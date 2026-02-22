package users

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUserMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestUserRepository_FindByEmail(t *testing.T) {
	email := "test@example.com"

	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "role", "created_at", "updated_at"}).
				AddRow(uuid.New(), "testuser", email, "secret", "user", time.Now(), time.Now()))

		user, err := repo.FindByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.FindByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnError(assert.AnError)

		user, err := repo.FindByEmail(context.Background(), email)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	username := "testuser"

	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "role", "created_at", "updated_at"}).
				AddRow(uuid.New(), username, "test@example.com", "secret", "user", time.Now(), time.Now()))

		user, err := repo.FindByUsername(context.Background(), username)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.FindByUsername(context.Background(), username)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnError(assert.AnError)

		user, err := repo.FindByUsername(context.Background(), username)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_InsertNewUser(t *testing.T) {
	user := &User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
	}

	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		id, err := repo.InsertNewUser(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, id)
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT").
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		id, err := repo.InsertNewUser(context.Background(), user)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
	})
}

func TestUserRepository_GetUserByUserID(t *testing.T) {
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "role", "created_at", "updated_at"}).
				AddRow(userID, "testuser", "test@example.com", "secret", "user", time.Now(), time.Now()))

		user, err := repo.GetUserByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetUserByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupUserMockDB(t)
		repo := NewRepository(gormDB)

		mock.ExpectQuery("SELECT").
			WillReturnError(assert.AnError)

		user, err := repo.GetUserByUserID(context.Background(), userID)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
