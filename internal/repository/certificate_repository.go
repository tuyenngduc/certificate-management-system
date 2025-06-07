package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateRepository interface {
	GetAllCertificates(ctx context.Context) ([]*models.Certificate, error)
	GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error)
	DeleteCertificate(ctx context.Context, id primitive.ObjectID) error
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

func (r *certificateRepository) GetAllCertificates(ctx context.Context) ([]*models.Certificate, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certs []*models.Certificate
	if err := cursor.All(ctx, &certs); err != nil {
		return nil, err
	}
	return certs, nil
}
func (r *certificateRepository) GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) DeleteCertificate(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
