package repository

import (
	"context"
	"errors"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubjectRepository interface {
	Create(ctx context.Context, subject *models.Subject) error
	GetByCode(ctx context.Context, code string) (*models.Subject, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Subject, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]*models.Subject, error)
}

type subjectRepository struct {
	collection *mongo.Collection
}

func NewSubjectRepository(db *mongo.Database) SubjectRepository {
	return &subjectRepository{
		collection: db.Collection("subjects"),
	}
}

func (r *subjectRepository) Create(ctx context.Context, subject *models.Subject) error {
	subject.CreatedAt = time.Now().Unix()
	subject.UpdatedAt = subject.CreatedAt
	_, err := r.collection.InsertOne(ctx, subject)
	return err
}

func (r *subjectRepository) GetByCode(ctx context.Context, code string) (*models.Subject, error) {
	var subject models.Subject
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&subject)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &subject, err
}

func (r *subjectRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Subject, error) {
	var subject models.Subject
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&subject)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &subject, err
}

func (r *subjectRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now().Unix()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *subjectRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("không tìm thấy môn học")
	}
	return nil
}

func (r *subjectRepository) List(ctx context.Context) ([]*models.Subject, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subjects []*models.Subject
	for cursor.Next(ctx) {
		var subject models.Subject
		if err := cursor.Decode(&subject); err != nil {
			return nil, err
		}
		subjects = append(subjects, &subject)
	}
	return subjects, nil
}
