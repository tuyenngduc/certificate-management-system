package repository

import (
	"context"
	"errors"

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
	GetAllCertificates(ctx context.Context) ([]*models.Certificate, error)
	FindByRegistrationNumber(ctx context.Context, reg string) (*models.Certificate, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type certificateRepository struct {
	db *mongo.Database
}

func NewCertificateRepository(db *mongo.Database) CertificateRepository {
	return &certificateRepository{db: db}
}
func (r *certificateRepository) GetAllCertificates(ctx context.Context) ([]*models.Certificate, error) {
	cursor, err := r.db.Collection("certificates").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certs []*models.Certificate
	for cursor.Next(ctx) {
		var cert models.Certificate
		if err := cursor.Decode(&cert); err != nil {
			return nil, err
		}
		certs = append(certs, &cert)
	}
	return certs, nil
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
func (r *certificateRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.db.Collection("certificates").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("chứng chỉ không tồn tại")
	}
	return nil
}
