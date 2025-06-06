package repository

import (
	"context"
	"log"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, update bson.M) error
	SearchUsers(ctx context.Context, params models.SearchUserParams) ([]*models.User, int64, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	FindByStudentID(ctx context.Context, studentID string) (*models.User, error)
	ExistsByStudentID(ctx context.Context, studentID string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
type userRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	col := db.Collection("users")
	repo := &userRepository{col: col}
	if err := repo.initIndexes(context.Background()); err != nil {
		log.Fatal("Cannot create indexes:", err)
	}
	return repo
}
func (r *userRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
func (r *userRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) SearchUsers(ctx context.Context, params models.SearchUserParams) ([]*models.User, int64, error) {
	filter := bson.M{}
	if params.StudentID != "" {
		filter["studentId"] = bson.M{"$regex": params.StudentID, "$options": "i"}
	}
	if params.FullName != "" {
		filter["fullName"] = bson.M{"$regex": params.FullName, "$options": "i"}
	}
	if params.Email != "" {
		filter["email"] = bson.M{"$regex": params.Email, "$options": "i"}
	}
	if params.Class != "" {
		filter["classId"] = bson.M{"$regex": params.Class, "$options": "i"}
	}
	if params.Faculty != "" {
		filter["facultyId"] = bson.M{"$regex": params.Faculty, "$options": "i"}
	}

	skip := int64((params.Page - 1) * params.PageSize)
	limit := int64(params.PageSize)

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.col.InsertOne(ctx, user)
	return err
}
func (r *userRepository) UpdateUser(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	result, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *userRepository) initIndexes(ctx context.Context) error {
	studentIDIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "studentId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{studentIDIndex, emailIndex})
	return err
}
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) FindByStudentID(ctx context.Context, studentID string) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"studentId": studentID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *userRepository) ExistsByStudentID(ctx context.Context, studentID string) (bool, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"studentId": studentID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
