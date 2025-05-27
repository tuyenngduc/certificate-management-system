package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, h *handler.AuthHandler) {
	auth := rg.Group("/auth")
	auth.POST("/request-otp", h.RequestOTP)
	auth.POST("/verify-otp", h.VerifyOTP)
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

}
