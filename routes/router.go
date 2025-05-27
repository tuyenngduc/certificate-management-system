package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
	"github.com/tuyenngduc/certificate-management-system/internal/middleware"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	trainingDepartmentHandler *handler.TrainingDepartmentHandler,
	authHandler *handler.AuthHandler,
	accountHandler *handler.AccountHandler,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// Public routes (không cần middleware)
	public := api.Group("/")
	{
		RegisterAuthRoutes(public, authHandler)
	}

	// Protected routes (cần middleware xác thực)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		RegisterUserRoutes(protected, userHandler)
		RegisterTrainingDepartmentRoutes(protected, trainingDepartmentHandler)
		RegisterAccountRoutes(protected, accountHandler)
	}

	return r
}
