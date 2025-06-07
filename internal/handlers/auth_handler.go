package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
func (h *AuthHandler) GetAllAccounts(c *gin.Context) {
	accounts, err := h.authService.GetAllAccounts(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var resp []models.AccountResponse
	for _, acc := range accounts {
		resp = append(resp, models.AccountResponse{
			ID:            acc.ID,
			StudentID:     acc.StudentID,
			StudentEmail:  acc.StudentEmail,
			PersonalEmail: acc.PersonalEmail,
			CreatedAt:     acc.CreatedAt.Format(time.RFC3339),
			Role:          acc.Role,
		})
	}

	c.JSON(200, resp)
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var input models.RequestOTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if err := h.authService.RequestOTP(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã gửi mã OTP tới email sinh viên"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	res, err := h.authService.VerifyOTP(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": res,
	})

}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tạo tài khoản thành công"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	account, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(account.ID, account.StudentID, account.Role, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tạo được token"})
		return
	}

	resp := models.LoginResponse{
		Token: token,
		Role:  account.Role,
	}

	c.JSON(http.StatusOK, resp)
}
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email là bắt buộc"})
		return
	}

	err := h.authService.DeleteAccountByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, common.ErrAccountUniversityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Trường không tồn tại"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa tài khoản thất bại: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa tài khoản thành công"})
}
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(400, gin.H{"errors": errs})
			return
		}
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Không xác thực"})
		return
	}

	claims, ok := claimsRaw.(*utils.CustomClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Dữ liệu xác thực không hợp lệ"})
		return
	}

	accountID, err := primitive.ObjectIDFromHex(claims.AccountID)
	if err != nil {
		c.JSON(401, gin.H{"error": "ID tài khoản không hợp lệ"})
		return
	}

	err = h.authService.ChangePassword(c.Request.Context(), accountID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case common.ErrAccountNotFound:
			c.JSON(404, gin.H{"error": "Không tìm thấy tài khoản"})
		case common.ErrInvalidOldPassword:
			c.JSON(400, gin.H{"error": "Mật khẩu cũ không đúng"})
		default:
			c.JSON(500, gin.H{"error": "Lỗi hệ thống"})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Đổi mật khẩu thành công"})
}
