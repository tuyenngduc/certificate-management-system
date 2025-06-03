package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepository interface {
	SaveOTP(ctx context.Context, otp models.OTP) error
	FindLatestOTPByEmail(ctx context.Context, email string) (*models.OTP, error)
	IsPersonalEmailExist(ctx context.Context, email string) (bool, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	CreateAccount(ctx context.Context, acc *models.Account) error
	FindByPersonalEmail(ctx context.Context, email string) (*models.Account, error)
}

type authRepository struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) SaveOTP(ctx context.Context, otp models.OTP) error {
	_, err := r.db.Collection("otp").InsertOne(ctx, otp)
	return err
}

func (r *authRepository) FindLatestOTPByEmail(ctx context.Context, email string) (*models.OTP, error) {
	var otp models.OTP
	opts := options.FindOne().SetSort(bson.D{{Key: "expires_at", Value: -1}}) // lấy bản ghi mới nhất
	err := r.db.Collection("otp").FindOne(ctx, bson.M{"email": email}, opts).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *authRepository) IsPersonalEmailExist(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"personal_email": email}
	count, err := r.db.Collection("accounts").CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *authRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateAccount(ctx context.Context, acc *models.Account) error {
	_, err := r.db.Collection("accounts").InsertOne(ctx, acc)
	return err
}

func (r *authRepository) FindByPersonalEmail(ctx context.Context, email string) (*models.Account, error) {
	var account models.Account
	collection := r.db.Collection("accounts")

	err := collection.FindOne(ctx, bson.M{"personal_email": email}).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
