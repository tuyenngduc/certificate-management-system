package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/response"
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
func (h *CertificateHandler) GetAllCertificates(c *gin.Context) {
	certs, err := h.certService.GetAllCertificates(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Không thể lấy danh sách chứng chỉ"})
		return
	}

	var resp []response.CertificateResponse
	for _, cert := range certs {
		resp = append(resp, response.ToCertificateResponse(cert))
	}

	c.JSON(200, resp)
}
func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
	var req request.CreateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				switch fe.Field() {
				case "UserID":
					c.JSON(400, gin.H{"error": "Người dùng không được để trống"})
					return
				case "CertificateType":
					c.JSON(400, gin.H{"error": "Loại văn bằng/chứng chỉ không hợp lệ. Chỉ chấp nhận 'degree' hoặc 'certificate'"})
					return
				case "Name":
					c.JSON(400, gin.H{"error": "Tên văn bằng/chứng chỉ không được để trống"})
					return
				case "SerialNumber":
					c.JSON(400, gin.H{"error": "Số hiệu không được để trống"})
					return
				case "RegistrationNumber":
					c.JSON(400, gin.H{"error": "Số vào sổ không được để trống"})
					return
				}
			}
		}
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ", "detail": err.Error()})
		return
	}

	cert, err := h.certService.CreateCertificate(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "user không tồn tại":
			c.JSON(404, gin.H{"error": err.Error()})
		case "id không hợp lệ":
			c.JSON(400, gin.H{"error": err.Error()})
		case "số hiệu đã tồn tại", "số vào sổ gốc cấp văn bằng đã tồn tại":
			c.JSON(409, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "Lỗi hệ thống", "detail": err.Error()})
		}
		return
	}

	c.JSON(201, cert)
}
func (h *CertificateHandler) GetCertificateByID(c *gin.Context) {
	idStr := c.Param("id")
	certID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã chứng chỉ không hợp lệ"})
		return
	}

	cert, err := h.certService.GetCertificateByID(context.Background(), certID)
	if err != nil || cert == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chứng chỉ không tồn tại"})
		return
	}

	resp := response.ToCertificateResponse(cert)
	c.JSON(http.StatusOK, resp)
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
func (h *CertificateHandler) DeleteCertificate(c *gin.Context) {
	idStr := c.Param("id")
	certID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID chứng chỉ không hợp lệ"})
		return
	}

	err = h.certService.DeleteCertificateByID(c.Request.Context(), certID)
	if err != nil {
		if err.Error() == "chứng chỉ không tồn tại" {
			c.JSON(404, gin.H{"error": "Chứng chỉ không tồn tại"})
			return
		}
		c.JSON(500, gin.H{"error": "Lỗi hệ thống", "detail": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Xóa chứng chỉ thành công"})
}
