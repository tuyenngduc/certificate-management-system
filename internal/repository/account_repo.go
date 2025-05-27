package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountRepository interface {
	GetAllAccounts(ctx context.Context) ([]*models.Account, error)
}

type accountRepository struct {
	collection *mongo.Collection
}

func NewAccountRepository(db *mongo.Database) AccountRepository {
	return &accountRepository{
		collection: db.Collection("accounts"),
	}
}

func (r *accountRepository) GetAllAccounts(ctx context.Context) ([]*models.Account, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*models.Account
	for cursor.Next(ctx) {
		var acc models.Account
		if err := cursor.Decode(&acc); err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}

	return accounts, nil
}
