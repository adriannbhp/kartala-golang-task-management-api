package main


import (
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/database"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/logger"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func Run() error {
	r, _, err := SetupApp()
	if err != nil {
		return err
	}

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	logger.Logger.Infof("Server starting on port %s", port)
	return r.Run(":" + port)
}

func SetupApp() (*gin.Engine, *config.DatabaseConfig, error) {
	// Initialize logger
	logger.SetupLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	// Connect to database
	_, err = database.NewDatabaseConnection(cfg.Database)
	if err != nil {
		return nil, nil, err
	}

	// Setup repositories, usecases, and handlers

	// Setup router
	r := gin.New()

	return r, &cfg.Database, nil
}
