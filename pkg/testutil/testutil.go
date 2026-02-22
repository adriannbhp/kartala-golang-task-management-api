package testutil

import (
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"errors"
	"os"
)

// SetupTestDB initializes a database connection for testing and runs auto-migrations
func SetupTestDB() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Override DBName with DBNameTest for integration tests
	cfg.Database.DBName = cfg.Database.DBNameTest

	db, err := database.NewDatabaseConnection(cfg.Database)
	if err != nil {
		return nil, err
	}

	// Run migrations for test environment
	err = db.AutoMigrate(
		append(users.GetDBModels(), tasks.GetDBModels()...)...,
	)
	if os.Getenv("FORCE_MIGRATE_ERROR") == "true" {
		return nil, errors.New("forced migration error")
	}

	return db, nil
}

// SetupTestRouter returns a Gin engine in test mode
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

// ClearTables removes all records from the specified tables
func ClearTables(db *gorm.DB, tables ...string) {
	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
	}
}
