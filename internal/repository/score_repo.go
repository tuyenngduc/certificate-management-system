package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScoreRepository struct {
	collection *mongo.Collection
}

func NewScoreRepository(db *mongo.Database) *ScoreRepository {
	return &ScoreRepository{
		collection: db.Collection("scores"),
	}
}

func (r *ScoreRepository) CreateScore(ctx context.Context, score *models.Score) error {
	_, err := r.collection.InsertOne(ctx, score)
	return err
}

func (r *ScoreRepository) IsScoreExists(ctx context.Context, studentID, subjectID, semester string) (bool, error) {
	studentObjID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return false, err
	}
	subjectObjID, err := primitive.ObjectIDFromHex(subjectID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"student_id": studentObjID,
		"subject_id": subjectObjID,
		"semester":   semester,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
