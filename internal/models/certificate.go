package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Certificate struct {
	ID              primitive.ObjectID `bson:"_id"`
	UserID          primitive.ObjectID `bson:"user_id"`
	ScoreID         primitive.ObjectID `bson:"score_id"`
	FacultyID       primitive.ObjectID `bson:"faculty_id"`
	UniversityID    primitive.ObjectID `bson:"university_id"`
	StudentCode     string             `bson:"student_code"`
	CertificateType string             `bson:"certificate_type"` //1 xuất săc, //2 giỏi, //3 khóa //4trung bình //5 yếu
	Name            string             `bson:"name"`
	SerialNumber    string             `bson:"serial_number"`
	RegNo           string             `bson:"registration_number"`
	Path            string             `bson:"path"`
	Signed          bool               `bson:"signed"`
	SignedAt        time.Time          `bson:"signed_at,omitempty"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`

	VerificationCode string    `bson:"verification_code,omitempty"`
	CodeExpiredAt    time.Time `bson:"code_expired_at,omitempty"`
}

type CreateCertificateRequest struct {
	StudentCode     string `json:"student_code" binding:"required"`
	CertificateType string `json:"certificate_type" binding:"required"`
	Name            string `json:"name" binding:"required"`
	SerialNumber    string `json:"serial_number" binding:"required"`
	RegNo           string `json:"reg_no" binding:"required"`
}

type CertificateResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	StudentCode     string    `json:"student_code"`
	StudentName     string    `json:"student_name"`
	CertificateType string    `json:"certificate_type"`
	Name            string    `json:"name"`
	SerialNumber    string    `json:"serial_number"`
	RegNo           string    `json:"reg_no"`
	Path            string    `bson:"path"`
	FacultyCode     string    `json:"faculty_code"`
	FacultyName     string    `json:"faculty_name"`
	UniversityCode  string    `json:"university_code"`
	UniversityName  string    `json:"university_name"`
	Signed          bool      `json:"signed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
type SearchCertificateParams struct {
	StudentCode     string `form:"student_code"`
	FacultyCode     string `form:"faculty_code"`
	Course          string `form:"course"`
	Signed          *bool  `form:"signed"`
	CertificateType string `form:"certificate_type"`
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=10"`
}
