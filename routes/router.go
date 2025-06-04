package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handlers"
)

func SetupRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	api := r.Group("/api/v1")
	// User routes
	api.GET("/users", userHandler.GetAllUsers)
	api.POST("/users", userHandler.CreateUser)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.GET("/users/search", userHandler.SearchUsers)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	// Auth routes
	api.POST("/auth/login", authHandler.Login)
	api.GET("/auth/accounts", authHandler.GetAllAccounts)
	api.POST("/auth/request-otp", authHandler.RequestOTP)
	api.POST("/auth/verify-otp", authHandler.VerifyOTP)
	api.POST("/auth/register", authHandler.Register)

	return r
}
