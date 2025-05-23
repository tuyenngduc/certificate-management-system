package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Cán bộ giáo viên
type Lecturer struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code     string             `bson:"code" json:"code"` // Mã cán bộ giáo viên, ví dụ: HT3831
	FullName string             `bson:"fullName" json:"full_name"`
	Email    string             `bson:"email" json:"email"`
	Title    string             `bson:"title" json:"title"` // Chức danh: ThS, TS, PGS, GS
}

// Lớp học
type Class struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Code     string               `bson:"code" json:"code"`
	Course   string               `bson:"course" json:"course"`
	Students []primitive.ObjectID `bson:"students,omitempty" json:"students,omitempty"` // Danh sách sinh viên thuộc lớp
}

// Khoa
type Faculty struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Code           string             `bson:"code" json:"code"`
	TrainingPeriod string             `bson:"trainingPeriod" json:"training_period"`          // Thời gian đào tạo
	Classes        []Class            `bson:"classes,omitempty" json:"classes,omitempty"`     // Danh sách lớp thuộc khoa
	Lecturers      []Lecturer         `bson:"lecturers,omitempty" json:"lecturers,omitempty"` // Danh sách cán bộ giáo viên thuộc khoa
}

// Phòng đào tạo
type TrainingDepartment struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name" json:"name"`
	Faculties []Faculty            `bson:"faculties,omitempty" json:"faculties,omitempty"` // Danh sách khoa
	Students  []primitive.ObjectID `bson:"students,omitempty" json:"students,omitempty"`   // Danh sách sinh viên toàn trường
}
