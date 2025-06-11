package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VerificationRepository interface {
	Save(ctx context.Context, code *models.VerificationCode) error
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.VerificationCode, error)
}

type verificationRepository struct {
	collection *mongo.Collection
}

func NewVerificationRepository(db *mongo.Database) VerificationRepository {
	return &verificationRepository{
		collection: db.Collection("verification_codes"),
	}
}

func (r *verificationRepository) Save(ctx context.Context, code *models.VerificationCode) error {
	_, err := r.collection.InsertOne(ctx, code)
	return err
}

func (r *verificationRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.VerificationCode, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var codes []models.VerificationCode
	if err := cursor.All(ctx, &codes); err != nil {
		return nil, err
	}
	return codes, nil
}
