package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerificationService interface {
	CreateVerificationCode(ctx context.Context, code *models.VerificationCode) error
	GetCodesByUser(ctx context.Context, userID primitive.ObjectID) ([]models.VerificationCodeResponse, error)
}

type verificationService struct {
	repo repository.VerificationRepository
}

func NewVerificationService(repo repository.VerificationRepository) VerificationService {
	return &verificationService{repo: repo}
}

func (s *verificationService) CreateVerificationCode(ctx context.Context, code *models.VerificationCode) error {
	code.ID = primitive.NewObjectID()
	code.Code = generateRandomCode(8)
	code.CreatedAt = time.Now()
	code.ViewedScore = false
	code.ViewedData = false
	code.ViewedFile = false
	return s.repo.Save(ctx, code)
}

func generateRandomCode(length int) string {
	return uuid.New().String()[:length]
}

func (s *verificationService) GetCodesByUser(ctx context.Context, userID primitive.ObjectID) ([]models.VerificationCodeResponse, error) {
	codes, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var res []models.VerificationCodeResponse
	for _, code := range codes {
		minutesRemaining := int64(code.ExpiredAt.Sub(now).Minutes())
		if minutesRemaining < 0 {
			minutesRemaining = 0
		}

		res = append(res, models.VerificationCodeResponse{
			ID:               code.ID,
			Code:             code.Code,
			CanViewScore:     code.CanViewScore,
			CanViewData:      code.CanViewData,
			CanViewFile:      code.CanViewFile,
			ViewedScore:      code.ViewedScore,
			ViewedData:       code.ViewedData,
			ViewedFile:       code.ViewedFile,
			ExpiredInMinutes: minutesRemaining,
			CreatedAt:        code.CreatedAt,
		})
	}

	return res, nil
}
