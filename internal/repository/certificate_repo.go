package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error)
	CreateCertificate(ctx context.Context, cert *models.Certificate) error
	UpdateCertificate(cert *models.Certificate) error
	GetCertificateByID(id primitive.ObjectID) (*models.Certificate, error)
	FindBySerialNumber(ctx context.Context, serial string) (*models.Certificate, error)
	FindByRegistrationNumber(ctx context.Context, reg string) (*models.Certificate, error)
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
	_, err := r.db.Collection("certificates").InsertOne(ctx, cert)
	return err
}
func (r *certificateRepository) GetCertificateByID(id primitive.ObjectID) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.db.Collection("certificates").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&cert)
	return &cert, err
}

func (r *certificateRepository) UpdateCertificate(cert *models.Certificate) error {
	_, err := r.db.Collection("certificates").UpdateOne(
		context.TODO(),
		bson.M{"_id": cert.ID},
		bson.M{"$set": cert},
	)
	return err
}

func (r *certificateRepository) FindBySerialNumber(ctx context.Context, serial string) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.db.Collection("certificates").FindOne(ctx, bson.M{"serial_number": serial}).Decode(&cert)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &cert, err
}

func (r *certificateRepository) FindByRegistrationNumber(ctx context.Context, reg string) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.db.Collection("certificates").FindOne(ctx, bson.M{"registration_number": reg}).Decode(&cert)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &cert, err
}
