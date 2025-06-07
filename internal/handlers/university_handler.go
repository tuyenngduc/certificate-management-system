package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
)

type UniversityHandler struct {
	universityService service.UniversityService
}

func NewUniversityHandler(s service.UniversityService) *UniversityHandler {
	return &UniversityHandler{universityService: s}
}

func (h *UniversityHandler) CreateUniversity(c *gin.Context) {
	var req models.CreateUniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("BindJSON error: %v", err)
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(400, gin.H{"errors": errs})
			return
		}
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	err := h.universityService.CreateUniversity(c.Request.Context(), &req)
	switch err {
	case common.ErrUniversityNameExists:
		c.JSON(400, gin.H{"error": "Tên trường đã tồn tại"})
		return
	case common.ErrUniversityEmailDomainExists:
		c.JSON(400, gin.H{"error": "Tên miền email đã tồn tại"})
		return
	case common.ErrUniversityCodeExists:
		c.JSON(400, gin.H{"error": "Mã trường đã tồn tại"})
		return
	}
	c.JSON(200, gin.H{"message": "Đã gửi yêu cầu sử dụng hệ thống, chờ admin phê duyệt"})
}

func (h *UniversityHandler) ApproveOrRejectUniversity(c *gin.Context) {
	var req models.ApproveOrRejectUniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("BindJSON error: %v", err)
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(400, gin.H{"errors": errs})
			return
		}
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	err := h.universityService.ApproveOrRejectUniversity(c.Request.Context(), req.UniversityID, req.Action)
	if err != nil {
		switch err {
		case common.ErrUniversityNotFound:
			c.JSON(404, gin.H{"error": "Không tìm thấy trường"})
			return
		case common.ErrUniversityCodeExists:
			c.JSON(409, gin.H{"error": "Mã trường đã tồn tại"})
			return
		default:
			c.JSON(500, gin.H{"error": "Lỗi hệ thống: " + err.Error()})
			return
		}
	}

	switch req.Action {
	case "approve":
		c.JSON(200, gin.H{"message": "Trường đã được phê duyệt và đã gửi tài khoản qua email"})
	case "reject":
		c.JSON(200, gin.H{"message": "Đã từ chối và xóa trường khỏi hệ thống"})
	default:
		c.JSON(400, gin.H{"error": "Hành động không hợp lệ"})
	}
}

func (h *UniversityHandler) GetAllUniversities(c *gin.Context) {
	universities, err := h.universityService.GetAllUniversities(c.Request.Context())
	if err != nil {
		log.Printf("Error getting universities: %v", err)
		c.JSON(500, gin.H{"error": "Lỗi hệ thống"})
		return
	}
	var resp []models.UniversityResponse
	for _, u := range universities {
		resp = append(resp, models.UniversityResponse{
			ID:             u.ID.Hex(),
			UniversityName: u.UniversityName,
			UniversityCode: u.UniversityCode,
			EmailDomain:    u.EmailDomain,
			Address:        u.Address,
			Status:         u.Status,
		})
	}

	c.JSON(200, gin.H{"data": resp})
}

func (h *UniversityHandler) GetApprovedUniversities(c *gin.Context) {
	universities, err := h.universityService.GetApprovedUniversities(c.Request.Context())
	if err != nil {
		log.Printf("Error getting approved universities: %v", err)
		c.JSON(500, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	var resp []models.UniversityResponse
	for _, u := range universities {
		resp = append(resp, models.UniversityResponse{
			ID:             u.ID.Hex(),
			UniversityName: u.UniversityName,
			UniversityCode: u.UniversityCode,
			EmailDomain:    u.EmailDomain,
			Address:        u.Address,
			Status:         u.Status,
		})
	}

	c.JSON(200, gin.H{"data": resp})
}
