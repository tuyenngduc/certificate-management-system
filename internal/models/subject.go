package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subject struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Credit      int                `bson:"credit" json:"credit"`
	FacultyID   primitive.ObjectID `bson:"faculty_id" json:"faculty_id"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   int64              `bson:"created_at" json:"created_at"`
	UpdatedAt   int64              `bson:"updated_at" json:"updated_at"`
}
