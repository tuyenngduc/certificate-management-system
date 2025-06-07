package service

import (
	"context"
	"errors"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UniversityService interface {
	CreateUniversity(ctx context.Context, req *models.CreateUniversityRequest) error
	ApproveOrRejectUniversity(ctx context.Context, idStr string, action string) error
	GetAllUniversities(ctx context.Context) ([]*models.University, error)
	GetApprovedUniversities(ctx context.Context) ([]*models.University, error)
}

type universityService struct {
	universityRepo repository.UniversityRepository
}

func NewUniversityService(universityRepo repository.UniversityRepository) UniversityService {
	return &universityService{universityRepo: universityRepo}
}

func (s *universityService) CreateUniversity(ctx context.Context, req *models.CreateUniversityRequest) error {
	conflictField, err := s.universityRepo.CheckUniversityConflicts(ctx, req.UniversityName, req.EmailDomain, req.UniversityCode)
	if err != nil {
		return err
	}
	switch conflictField {
	case "university_name":
		return common.ErrUniversityNameExists
	case "email_domain":
		return common.ErrUniversityEmailDomainExists
	case "university_code":
		return common.ErrUniversityCodeExists
	}

	uni := &models.University{
		ID:             primitive.NewObjectID(),
		UniversityName: req.UniversityName,
		Address:        req.Address,
		EmailDomain:    req.EmailDomain,
		UniversityCode: req.UniversityCode,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.universityRepo.CreateUniversity(ctx, uni)
}

func (s *universityService) ApproveOrRejectUniversity(ctx context.Context, idStr string, action string) error {
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return common.ErrUniversityNotFound
	}

	university, err := s.universityRepo.FindByID(ctx, objID)
	if err != nil || university == nil {
		return common.ErrUniversityNotFound
	}

	switch action {
	case "approve":
		return s.universityRepo.UpdateStatus(ctx, objID, "approved")
	case "reject":
		return s.universityRepo.DeleteByID(ctx, objID)
	default:
		return errors.New("invalid action")
	}
}
func (s *universityService) GetAllUniversities(ctx context.Context) ([]*models.University, error) {
	return s.universityRepo.GetAllUniversities(ctx)
}
func (s *universityService) GetApprovedUniversities(ctx context.Context) ([]*models.University, error) {
	return s.universityRepo.GetApprovedUniversities(ctx)
}
