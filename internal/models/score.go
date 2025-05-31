package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Score struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	StudentID primitive.ObjectID `bson:"student_id"`
	SubjectID primitive.ObjectID `bson:"subject_id"`
	Semester  string             `bson:"semester"`

	Attendance   float64   `bson:"attendance"`
	Midterm      float64   `bson:"midterm"`
	Final        float64   `bson:"final"`
	ProcessScore float64   `bson:"process_score"`
	Total        float64   `bson:"total"`
	GPAChar      string    `bson:"gpa_char"`
	Passed       bool      `bson:"passed"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}
