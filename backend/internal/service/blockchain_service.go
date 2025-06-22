package service

import (
	"context"
	"fmt"

	"github.com/tuyenngduc/certificate-management-system/backend/internal/common"
	"github.com/tuyenngduc/certificate-management-system/backend/internal/models"
	"github.com/tuyenngduc/certificate-management-system/backend/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/backend/pkg/blockchain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainService interface {
	PushCertificateToChain(ctx context.Context, certificateID primitive.ObjectID) (string, error)
	GetCertificateFromChain(ctx context.Context, certificateID string) (*models.CertificateOnChain, error)
}

type blockchainService struct {
	certRepo       repository.CertificateRepository
	userRepo       repository.UserRepository
	facultyRepo    repository.FacultyRepository
	universityRepo repository.UniversityRepository
	fabricClient   *blockchain.FabricClient
}

func NewBlockchainService(
	certRepo repository.CertificateRepository,
	userRepo repository.UserRepository,
	facultyRepo repository.FacultyRepository,
	universityRepo repository.UniversityRepository,
	fabricClient *blockchain.FabricClient,
) BlockchainService {
	return &blockchainService{
		certRepo:       certRepo,
		userRepo:       userRepo,
		facultyRepo:    facultyRepo,
		universityRepo: universityRepo,
		fabricClient:   fabricClient,
	}
}

func (s *blockchainService) PushCertificateToChain(ctx context.Context, certificateID primitive.ObjectID) (string, error) {
	cert, err := s.certRepo.GetCertificateByID(ctx, certificateID)
	if err != nil || cert == nil {
		return "", common.ErrCertificateNotFound
	}
	if cert.CertHash == "" {
		return "", fmt.Errorf("certificate chưa có cert_hash")
	}

	chainData := models.CertificateOnChain{
		CertID:              cert.ID.Hex(),
		CertHash:            cert.CertHash,
		UniversitySignature: "",
		DateOfIssuing:       cert.IssueDate.Format("2006-01-02"),
		SerialNumber:        cert.SerialNumber,
		RegNo:               cert.RegNo,
		Version:             1,
		UpdatedDate:         cert.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	txID, err := s.fabricClient.IssueCertificate(chainData)
	if err != nil {
		return "", err
	}
	return txID, nil
}

func (s *blockchainService) GetCertificateFromChain(ctx context.Context, certificateID string) (*models.CertificateOnChain, error) {
	cert, err := s.fabricClient.GetCertificateByID(certificateID)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
