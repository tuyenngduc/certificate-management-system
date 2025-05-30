package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	issueDate, err := time.Parse(time.RFC3339, req.IssueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue_date format"})
		return
	}

	cert := &models.Certificate{
		UserID:            userID,
		CertificateType:   req.CertificateType,
		Name:              req.Name,
		Issuer:            req.Issuer,
		IssueDate:         issueDate,
		CertificateNumber: req.CertificateNumber,
		Status:            "issued", // mặc định trạng thái là issued
	}

	err = h.certService.IssueCertificate(context.Background(), cert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Certificate issued successfully",
		"certificate_id": cert.ID.Hex(),
	})
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
