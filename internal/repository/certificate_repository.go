package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateRepository interface {
	CreateCertificate(ctx context.Context, cert *models.Certificate) error
}
type certificateRepository struct {
	col *mongo.Collection
}

func NewCertificateRepository(db *mongo.Database) CertificateRepository {
	col := db.Collection("certificates")
	return &certificateRepository{col: col}
}

func (r *certificateRepository) CreateCertificate(ctx context.Context, cert *models.Certificate) error {
	_, err := r.col.InsertOne(ctx, cert)
	return err
}
