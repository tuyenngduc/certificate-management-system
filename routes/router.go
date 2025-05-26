package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	trainingDepartmentHandler *handler.TrainingDepartmentHandler,
) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")

	RegisterUserRoutes(api, userHandler)
	RegisterTrainingDepartmentRoutes(api, trainingDepartmentHandler)
	return r
}
