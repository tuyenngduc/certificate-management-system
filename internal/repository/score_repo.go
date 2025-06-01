package repository

import (
	"context"
	"time"

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

func (r *ScoreRepository) GetScoresByStudentID(ctx context.Context, studentID primitive.ObjectID) ([]*models.Score, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []*models.Score
	for cursor.Next(ctx) {
		var score models.Score
		if err := cursor.Decode(&score); err != nil {
			return nil, err
		}
		scores = append(scores, &score)
	}
	return scores, nil
}
func (r *ScoreRepository) GetScoresBySubjectID(ctx context.Context, subjectID primitive.ObjectID) ([]*models.Score, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"subject_id": subjectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []*models.Score
	for cursor.Next(ctx) {
		var score models.Score
		if err := cursor.Decode(&score); err != nil {
			return nil, err
		}
		scores = append(scores, &score)
	}
	return scores, nil
}

func (r *ScoreRepository) UpdateScore(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *ScoreRepository) GetScoreByID(ctx context.Context, id primitive.ObjectID) (*models.Score, error) {
	var score models.Score
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&score)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &score, nil
}
