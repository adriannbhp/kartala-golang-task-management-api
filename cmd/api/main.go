// @title Kartala Task Management API
// @version 1.0
// @description API Server for Kartala Task Management API
// @termsOfService http://swagger.io/terms/
 
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
 
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
 
// @host localhost:8081
// @BasePath /
 
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
 
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your token.
 
package main


import (
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/config"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/database"
	"github.com/adriannbhp/kartala-golang-task-management-api/pkg/logger"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/middleware"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/delivery/http"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/auth"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var exitFunc = os.Exit

func main() {
	if err := Run(); err != nil {
		log.Printf("Application failed: %v", err)
		exitFunc(1)
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
		port = "8080"
	}

	if os.Getenv("TEST_MAIN_COVERAGE") == "true" {
		return nil
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
	db, err := database.NewDatabaseConnection(cfg.Database)
	if err != nil {
		return nil, nil, err
	}
 
	// Parse durations
	accessDuration, _ := time.ParseDuration(cfg.Auth.AccessTokenDuration)
	refreshDuration, _ := time.ParseDuration(cfg.Auth.RefreshTokenDuration)

	// Setup repositories, usecases, and handlers
	userRepo := users.NewRepository(db)
	userUsecase := users.NewUsecase(userRepo)
	userHandler := users.NewHandler(userUsecase)

	taskRepo := tasks.NewRepository(db)
	taskUsecase := tasks.NewUsecase(taskRepo)
	taskHandler := tasks.NewHandler(taskUsecase)

	authUsecase := auth.NewUsecase(userRepo, cfg.Secret.JwtSecret, accessDuration, refreshDuration)
	authHandler := auth.NewHandler(authUsecase, cfg.Secret.JwtSecret)
 

	// Setup router
	r := gin.New()
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger(logger.Logger))

	http.SetupRoutes(r, authHandler, userHandler, taskHandler, cfg.Secret.JwtSecret, cfg.Secret.ApiKey)


	return r, &cfg.Database, nil
}
