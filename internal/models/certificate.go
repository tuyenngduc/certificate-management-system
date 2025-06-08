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
	CertificateType string             `bson:"certificate_type"`
	Name            string             `bson:"name"`
	SerialNumber    string             `bson:"serial_number"`
	RegNo           string             `bson:"registration_number"`
	Path            string             `bson:"path"`
	Signed          bool               `bson:"signed"`
	SignedAt        time.Time          `bson:"signed_at,omitempty"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}

type CreateCertificateRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	CertificateType string `json:"certificate_type" binding:"required"`
	Name            string `json:"name" binding:"required"`
	SerialNumber    string `json:"serial_number" binding:"required"`
	RegNo           string `json:"reg_no" binding:"required"`
}

type CertificateResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	StudentCode     string    `json:"student_code"`
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

type CertificateVerification struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	CertificateID    primitive.ObjectID `bson:"certificate_id"`
	VerificationCode string             `bson:"verification_code"`
	ExpiresAt        time.Time          `bson:"expires_at"`
	CreatedAt        time.Time          `bson:"created_at"`
}
