package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handlers"
	"github.com/tuyenngduc/certificate-management-system/internal/middleware"
)

func SetupRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	certificateHandler *handlers.CertificateHandler,
	universityHandler *handlers.UniversityHandler,
) *gin.Engine {
	r := gin.Default()

	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api/v1")

	// ===== Auth routes =====
	authPublic := api.Group("/auth")
	authPublic.POST("/login", authHandler.Login)
	authPublic.POST("/request-otp", authHandler.RequestOTP)
	authPublic.POST("/verify-otp", authHandler.VerifyOTP)
	authPublic.POST("/register", authHandler.Register)

	authPrivate := api.Group("/auth")
	authPrivate.Use(middleware.JWTAuthMiddleware())
	authPrivate.GET("/accounts", authHandler.GetAllAccounts)
	authPrivate.DELETE("/accounts", authHandler.DeleteAccount)
	authPrivate.POST("/change-password", authHandler.ChangePassword)

	// ===== User routes =====
	userGroup := api.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware()) // Bảo vệ toàn bộ route users
	userGroup.POST("/import-excel", userHandler.ImportUsersFromExcel)
	userGroup.GET("", userHandler.GetAllUsers)
	userGroup.POST("", userHandler.CreateUser)
	userGroup.GET("/:id", userHandler.GetUserByID)
	userGroup.GET("/search", userHandler.SearchUsers)
	userGroup.PUT("/:id", userHandler.UpdateUser)
	userGroup.DELETE("/:id", userHandler.DeleteUser)

	// ===== Certificate routes =====
	certificateGroup := api.Group("/certificates")
	certificateGroup.Use(middleware.JWTAuthMiddleware())
	certificateGroup.GET("", certificateHandler.GetAllCertificates)
	certificateGroup.POST("", certificateHandler.CreateCertificate)
	certificateGroup.GET("/:id", certificateHandler.GetCertificateByID)

	// ===== University routes =====
	universityGroup := api.Group("/universities")
	universityGroup.Use(middleware.JWTAuthMiddleware())
	universityGroup.POST("", universityHandler.CreateUniversity)
	universityGroup.POST("/approve-or-reject", universityHandler.ApproveOrRejectUniversity)
	universityGroup.GET("", universityHandler.GetAllUniversities)
	universityGroup.GET("/status", universityHandler.GetUniversities)

	return r
}
