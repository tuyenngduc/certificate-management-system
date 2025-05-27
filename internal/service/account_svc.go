package service

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
)

type AccountService interface {
	GetAllAccounts(ctx context.Context) ([]*models.Account, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) GetAllAccounts(ctx context.Context) ([]*models.Account, error) {
	return s.repo.GetAllAccounts(ctx)
}
