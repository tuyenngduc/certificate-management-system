package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
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
	CreateAccount(ctx context.Context, acc *models.Account) error
	FindByPersonalEmail(ctx context.Context, email string) (*models.Account, error)
	GetAllAccounts(ctx context.Context) ([]*models.Account, error)
	DeleteAccountByEmail(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, accountID primitive.ObjectID, newHash string) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Account, error)
	FindPersonalAccountByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Account, error)
}

type authRepository struct {
	col *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	col := db.Collection("accounts")
	return &authRepository{col: col}
}

func (r *authRepository) SaveOTP(ctx context.Context, otp models.OTP) error {
	_, err := r.col.InsertOne(ctx, otp)
	return err
}

func (r *authRepository) FindLatestOTPByEmail(ctx context.Context, email string) (*models.OTP, error) {
	var otp models.OTP
	opts := options.FindOne().SetSort(bson.D{{Key: "expires_at", Value: -1}}) // lấy bản ghi mới nhất
	err := r.col.FindOne(ctx, bson.M{"email": email}, opts).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *authRepository) IsPersonalEmailExist(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"personal_email": email}
	count, err := r.col.CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *authRepository) CreateAccount(ctx context.Context, acc *models.Account) error {
	_, err := r.col.InsertOne(ctx, acc)
	return err
}

func (r *authRepository) FindByPersonalEmail(ctx context.Context, email string) (*models.Account, error) {
	var account models.Account

	err := r.col.FindOne(ctx, bson.M{"personal_email": email}).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *authRepository) GetAllAccounts(ctx context.Context) ([]*models.Account, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*models.Account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}
func (r *authRepository) DeleteAccountByEmail(ctx context.Context, email string) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"personal_email": email})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return common.ErrAccountUniversityNotFound
	}

	return nil
}
func (r *authRepository) UpdatePassword(ctx context.Context, accountID primitive.ObjectID, newHash string) error {
	filter := bson.M{"_id": accountID}
	update := bson.M{"$set": bson.M{"password_hash": newHash}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *authRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Account, error) {
	filter := bson.M{"_id": id}
	var acc models.Account
	err := r.col.FindOne(ctx, filter).Decode(&acc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}
func (r *authRepository) FindPersonalAccountByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Account, error) {
	filter := bson.M{
		"student_id":     userID,
		"personal_email": bson.M{"$ne": ""},
	}

	var account models.Account
	err := r.col.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}
