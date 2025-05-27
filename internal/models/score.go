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

	Attendance   float64 `bson:"attendance"`    // điểm chuyên cần
	Midterm      float64 `bson:"midterm"`       // điểm giữa kỳ
	Final        float64 `bson:"final"`         // điểm cuối kỳ
	ProcessScore float64 `bson:"process_score"` // điểm quá trình (tính từ attendance + midterm)
	Total        float64 `bson:"total"`         // điểm kết thúc học phần
	GPAChar      string  `bson:"gpa_char"`      // điểm chữ

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
