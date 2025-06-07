package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Faculty struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FacultyID    primitive.ObjectID `bson:"faculty_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	UniversityID primitive.ObjectID `bson:"university_id" json:"university_id"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
