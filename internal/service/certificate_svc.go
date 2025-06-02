package service

import (
	"context"
	"errors"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CertificateService struct {
	certRepo repository.CertificateRepository
	userRepo repository.UserRepository
}

func NewCertificateService(certRepo repository.CertificateRepository, userRepo repository.UserRepository) *CertificateService {
	return &CertificateService{
		certRepo: certRepo,
		userRepo: userRepo,
	}
}
func (s *CertificateService) GetAllCertificates(ctx context.Context) ([]*models.Certificate, error) {
	return s.certRepo.GetAllCertificates(ctx)
}
func (s *CertificateService) CreateCertificate(ctx context.Context, req request.CreateCertificateRequest) (*models.Certificate, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, errors.New("id không hợp lệ")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("lỗi truy vấn user")
	}
	if user == nil {
		return nil, errors.New("user không tồn tại")
	}

	exist, _ := s.certRepo.FindBySerialNumber(ctx, req.SerialNumber)
	if exist != nil {
		return nil, errors.New("số hiệu đã tồn tại")
	}
	exist, _ = s.certRepo.FindByRegistrationNumber(ctx, req.RegistrationNumber)
	if exist != nil {
		return nil, errors.New("số vào sổ gốc cấp văn bằng đã tồn tại")
	}

	cert := &models.Certificate{
		ID:                 primitive.NewObjectID(),
		UserID:             userID,
		CertificateType:    req.CertificateType,
		Name:               req.Name,
		Issuer:             "Học Viện Kỹ Thuật Mật Mã",
		IssueDate:          time.Now(),
		SerialNumber:       req.SerialNumber,
		RegistrationNumber: req.RegistrationNumber,
		Status:             "Chờ ký số",
		Signed:             false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.certRepo.CreateCertificate(ctx, cert); err != nil {
		return nil, err
	}

	return cert, nil
}

func (s *CertificateService) HashCertificateByID(id primitive.ObjectID) error {
	cert, err := s.certRepo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	hash, err := utils.HashCertificateData(cert)
	if err != nil {
		return err
	}

	cert.Hash = hash
	cert.UpdatedAt = time.Now()

	return s.certRepo.UpdateCertificate(cert)
}

func (s *CertificateService) GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error) {
	return s.certRepo.GetByID(ctx, id)
}

func (s *CertificateService) DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error {
	return s.certRepo.DeleteByID(ctx, id)
}
