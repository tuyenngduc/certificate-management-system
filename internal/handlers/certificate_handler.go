package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
)

type CertificateHandler struct {
	certificateService service.CertificateService
}

func NewCertificateHandler(certService service.CertificateService) *CertificateHandler {
	return &CertificateHandler{certificateService: certService}
}

func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
	var req models.CreateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	resp, err := h.certificateService.CreateCertificate(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case common.ErrInvalidUserID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID người dùng không hợp lệ"})
			return
		case common.ErrUserNotExisted:
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy user"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, resp)
}
