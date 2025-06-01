package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
	"github.com/tuyenngduc/certificate-management-system/internal/middleware"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	trainingDepartmentHandler *handler.TrainingDepartmentHandler,
	authHandler *handler.AuthHandler,
	accountHandler *handler.AccountHandler,
	subjectHandler *handler.SubjectHandler,
	scoreHandler *handler.ScoreHandler,
	certificateHandler *handler.CertificateHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	api := r.Group("/api/v1")

	public := api.Group("/")
	{
		RegisterAuthRoutes(public, authHandler)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		RegisterUserRoutes(protected, userHandler)
		RegisterTrainingDepartmentRoutes(protected, trainingDepartmentHandler)
		RegisterAccountRoutes(protected, accountHandler)
		RegisterSubjectRoutes(protected, subjectHandler)
		RegisterScoreRoutes(protected, scoreHandler)
		RegisterCertificateRoutes(protected, certificateHandler)
	}

	return r
}
