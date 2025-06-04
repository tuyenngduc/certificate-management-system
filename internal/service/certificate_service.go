package service

import (
	"context"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CertificateService interface {
	CreateCertificate(ctx context.Context, req *models.CreateCertificateRequest) (*models.CertificateResponse, error)
}

type certificateService struct {
	certificateRepo repository.CertificateRepository
	userRepo        repository.UserRepository
}

func NewCertificateService(certificateRepo repository.CertificateRepository, userRepo repository.UserRepository) CertificateService {
	return &certificateService{certificateRepo: certificateRepo, userRepo: userRepo}
}

func (s *certificateService) CreateCertificate(ctx context.Context, req *models.CreateCertificateRequest) (*models.CertificateResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, common.ErrInvalidUserID
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, common.ErrUserNotExisted
	}

	cert := &models.Certificate{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		StudentID:       user.StudentID,
		CertificateType: req.CertificateType,
		Name:            req.Name,
		Issuer:          req.Issuer,
		SerialNumber:    req.SerialNumber,
		RegNo:           req.RegNo,
		Signed:          false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.certificateRepo.CreateCertificate(ctx, cert)
	if err != nil {
		return nil, err
	}

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentID:       cert.StudentID,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		Issuer:          cert.Issuer,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		Signed:          cert.Signed,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
	}, nil
}
