package database

import (
	"fmt"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabaseConnection creates a new database connection using default postgres dialector
func NewDatabaseConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	return NewDatabaseConnectionWithDialector(postgres.Open(dsn), &gorm.Config{})
}

// NewDatabaseConnectionWithDialector creates a new database connection with a custom dialector
func NewDatabaseConnectionWithDialector(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	if dialector == nil {
		return nil, fmt.Errorf("dialector is required")
	}
	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
