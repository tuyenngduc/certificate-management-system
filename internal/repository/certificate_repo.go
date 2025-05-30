package repository

import (
	"context"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error)
	CreateCertificate(ctx context.Context, cert *models.Certificate) error
}

type certificateRepository struct {
	db *mongo.Database
}

func NewCertificateRepository(db *mongo.Database) CertificateRepository {
	return &certificateRepository{db: db}
}

func (r *certificateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.db.Collection("certificates").FindOne(ctx, bson.M{"_id": id}).Decode(&cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) CreateCertificate(ctx context.Context, cert *models.Certificate) error {
	cert.ID = primitive.NewObjectID()
	cert.CreatedAt = time.Now()
	cert.UpdatedAt = time.Now()

	_, err := r.db.Collection("certificates").InsertOne(ctx, cert)
	return err
}
