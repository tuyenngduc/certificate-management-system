package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerificationHandler struct {
	verificationService service.VerificationService
}

func NewVerificationHandler(verificationService service.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
	}
}

func (h *VerificationHandler) CreateVerificationCode(c *gin.Context) {
	var req models.CreateVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	claims, ok := c.Request.Context().Value(utils.ClaimsContextKey).(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID người dùng không hợp lệ"})
		return
	}

	expiredAt := time.Now().Add(time.Duration(req.DurationMinutes) * time.Minute)

	code := &models.VerificationCode{
		UserID:       userID,
		CanViewScore: req.CanViewScore,
		CanViewData:  req.CanViewData,
		CanViewFile:  req.CanViewFile,
		ExpiredAt:    expiredAt,
	}

	err = h.verificationService.CreateVerificationCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo mã xác minh"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       code.Code,
		"expired_at": code.ExpiredAt,
		"can_view": gin.H{
			"score": code.CanViewScore,
			"data":  code.CanViewData,
			"file":  code.CanViewFile,
		},
	})
}
func (h *VerificationHandler) GetMyCodes(c *gin.Context) {
	claims, ok := c.Request.Context().Value(utils.ClaimsContextKey).(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID không hợp lệ"})
		return
	}

	codes, err := h.verificationService.GetCodesByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy mã xác minh"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": codes})
}
