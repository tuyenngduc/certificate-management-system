package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Cán bộ giáo viên
type Lecturer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string             `bson:"code" json:"code"`
	FullName  string             `bson:"fullName" json:"full_name"`
	Email     string             `bson:"email" json:"email"`
	Title     string             `bson:"title" json:"title"`
	FacultyID primitive.ObjectID `bson:"faculty_id" json:"faculty_id"` // Giảng viên thuộc khoa nào
}

// Lớp học
type Class struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Course    string             `bson:"course" json:"course"`
	FacultyID primitive.ObjectID `bson:"faculty_id" json:"faculty_id"`
}

// Khoa
type Faculty struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Code           string             `bson:"code" json:"code"`
	TrainingPeriod string             `bson:"trainingPeriod" json:"training_period"` // Thời gian đào tạo
}
