package http
 
import (
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/auth"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/tasks"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/adriannbhp/kartala-golang-task-management-api/docs"
	"golang.org/x/time/rate"
)
 
// SetupRoutes initializes the API routes
func SetupRoutes(r *gin.Engine, authHandler *auth.Handler, userHandler *users.Handler, taskHandler *tasks.Handler, jwtSecret string, apiKey string) {
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
 
	// API V1
	v1 := r.Group("/api/v1")
	// Global Rate Limit: 5 request per second, burst of 10
	v1.Use(middleware.RateLimiter(rate.Limit(5), 10))
	v1.Use(middleware.ApiKeyMiddleware(apiKey))
	{
		// Auth routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
		}
 
		// User routes
		userGroup := v1.Group("/users")
		userGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			userGroup.GET("/me", userHandler.GetUserInfo)
		}

		// Task routes
		taskGroup := v1.Group("/tasks")
		taskGroup.Use(middleware.AuthMiddleware(jwtSecret))
		{
			taskGroup.POST("", taskHandler.CreateTask)
			taskGroup.GET("", taskHandler.GetAllTasks)
			taskGroup.GET("/:id", taskHandler.GetTaskByID)
			taskGroup.PUT("/:id", taskHandler.UpdateTask)
			taskGroup.DELETE("/:id", taskHandler.DeleteTask)
		}
	}
}
