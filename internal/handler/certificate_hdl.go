package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
)

type CertificateHandler struct {
	certService *service.CertificateService
}

func NewCertificateHandler(certService *service.CertificateService) *CertificateHandler {
	return &CertificateHandler{certService: certService}
}

func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
	var req request.CreateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cert, err := h.certService.CreateCertificate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cert)
}

func (h *CertificateHandler) GetCertificateByID(c *gin.Context) {
	idStr := c.Param("id")
	certID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã chứng chỉ không hợp lệ"})
		return
	}

	cert, err := h.certService.GetCertificateByID(context.Background(), certID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chứng chỉ không tồn tại"})
		return
	}

	c.JSON(http.StatusOK, cert)
}

func (h *CertificateHandler) HashCertificate(c *gin.Context) {
	idParam := c.Param("id")
	certID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID chứng chỉ không hợp lệ"})
		return
	}

	err = h.certService.HashCertificateByID(certID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật hash chứng chỉ thành công"})
}
