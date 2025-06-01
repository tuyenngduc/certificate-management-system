package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func RegisterCertificateRoutes(rg *gin.RouterGroup, handler *handler.CertificateHandler) {
	cert := rg.Group("/certificates")
	{
		cert.POST("/", handler.CreateCertificate)
		cert.POST("/:id/hash", handler.HashCertificate)
		cert.GET("/:id", handler.GetCertificateByID)
	}
}
