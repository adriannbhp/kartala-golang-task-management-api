package http
 
import (
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/auth"
	"github.com/adriannbhp/kartala-golang-task-management-api/internal/users"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/adriannbhp/kartala-golang-task-management-api/docs"
)
 
// SetupRoutes initializes the API routes
func SetupRoutes(r *gin.Engine, authHandler *auth.Handler, userHandler *users.Handler, jwtSecret string, apiKey string) {
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
 
	// API V1
	v1 := r.Group("/api/v1")
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
		{
			userGroup.GET("/me", userHandler.GetUserInfo)
		}
	}
}
