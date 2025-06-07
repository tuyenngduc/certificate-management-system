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
