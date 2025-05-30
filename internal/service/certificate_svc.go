package service

import (
	"context"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CertificateService struct {
	certRepo repository.CertificateRepository
}

func NewCertificateService(certRepo repository.CertificateRepository) *CertificateService {
	return &CertificateService{certRepo: certRepo}
}

func (s *CertificateService) IssueCertificate(ctx context.Context, cert *models.Certificate) error {
	hashData := models.BuildCertificateHashData(cert)

	hash, err := utils.ComputeCertificateHash(hashData)
	if err != nil {
		return err
	}

	cert.BlockchainHash = hash
	cert.BlockchainTimestamp = time.Now()

	return s.certRepo.CreateCertificate(ctx, cert)
}

func (s *CertificateService) GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error) {
	return s.certRepo.GetByID(ctx, id)
}
