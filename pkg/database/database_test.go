package database

import (
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewDatabaseConnection_Error(t *testing.T) {
	cfg := config.DatabaseConfig{Host: ""}
	db, err := NewDatabaseConnection(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewDatabaseConnectionWithDialector(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dbMock, mock, _ := sqlmock.New()
		defer dbMock.Close()
		dialector := postgres.New(postgres.Config{
			Conn: dbMock,
		})
		db, err := NewDatabaseConnectionWithDialector(dialector, &gorm.Config{})
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error_nil_dialector", func(t *testing.T) {
		db, err := NewDatabaseConnectionWithDialector(nil, &gorm.Config{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dialector is required")
		assert.Nil(t, db)
	})

	t.Run("error_invalid_dsn", func(t *testing.T) {
		db, err := NewDatabaseConnectionWithDialector(postgres.Open("invalid dsn"), &gorm.Config{})
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
