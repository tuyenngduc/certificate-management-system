package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	authRepo    repository.AuthRepository
	userRepo    *repository.UserRepository
	emailSender utils.EmailSender
}

func NewAuthService(authRepo repository.AuthRepository, userRepo *repository.UserRepository, emailSender utils.EmailSender) *AuthService {
	return &AuthService{
		authRepo:    authRepo,
		userRepo:    userRepo,
		emailSender: emailSender,
	}
}

func (s *AuthService) RequestOTP(ctx context.Context, input models.RequestOTPInput) error {
	_, err := s.userRepo.FindByEmail(ctx, input.StudentEmail)
	if err != nil {
		return fmt.Errorf("email không tồn tại trong hệ thống")
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	otpData := models.OTP{
		Email:     input.StudentEmail,
		Code:      otp,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}

	if err := s.authRepo.SaveOTP(ctx, otpData); err != nil {
		return err
	}

	body := fmt.Sprintf("Mã OTP của bạn là: %s. Có hiệu lực trong 3 phút.", otp)
	return s.emailSender.SendEmail(input.StudentEmail, "Mã xác thực OTP", body)
}

func (s *AuthService) VerifyOTP(ctx context.Context, input *models.VerifyOTPRequest) (string, error) {
	otpRecord, err := s.authRepo.FindLatestOTPByEmail(ctx, input.StudentEmail)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy mã OTP")
	}

	if otpRecord.Code != input.OTP {
		return "", fmt.Errorf("mã OTP không đúng")
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		return "", fmt.Errorf("mã OTP đã hết hạn")
	}

	user, err := s.userRepo.FindByEmail(ctx, input.StudentEmail)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy người dùng")
	}

	return user.ID.Hex(), nil
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) error {
	// Check xem PersonalEmail đã tồn tại chưa
	exists, err := s.authRepo.IsPersonalEmailExist(ctx, req.PersonalEmail)
	if err != nil {
		return fmt.Errorf("lỗi kiểm tra email: %w", err)
	}
	if exists {
		return fmt.Errorf("email cá nhân đã được sử dụng")
	}

	// Parse user_id
	userObjID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("user_id không hợp lệ")
	}

	// Lấy user để lấy email sinh viên
	user, err := s.authRepo.GetUserByID(ctx, userObjID)
	if err != nil {
		return fmt.Errorf("không tìm thấy user")
	}

	// Hash mật khẩu
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("lỗi hash mật khẩu: %w", err)
	}

	account := &models.Account{
		UserID:        user.ID,
		StudentEmail:  user.Email, // Email @actvn.edu.vn
		PersonalEmail: req.PersonalEmail,
		PasswordHash:  hash,
		CreatedAt:     time.Now(),
		Role:          "student",
	}

	if err := s.authRepo.CreateAccount(ctx, account); err != nil {
		return fmt.Errorf("không tạo được tài khoản: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.Account, error) {
	account, err := s.authRepo.FindByPersonalEmail(ctx, email)
	if err != nil {
		return nil, errors.New("tài khoản không tồn tại")
	}

	if !utils.ComparePassword(account.PasswordHash, password) {
		return nil, errors.New("mật khẩu không đúng")
	}

	return account, nil
}
