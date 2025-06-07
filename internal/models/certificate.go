package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Certificate struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          primitive.ObjectID `bson:"user_id"`
	StudentCode     string             `bson:"student_code"`
	CertificateType string             `bson:"certificate_type"`
	Name            string             `bson:"name"`
	Issuer          string             `bson:"issuer"`
	SerialNumber    string             `bson:"serial_number"`
	RegNo           string             `bson:"registration_number"`

	Signed    bool      `bson:"signed"`
	SignedAt  time.Time `bson:"signed_at,omitempty"`
	Signature string    `bson:"signature,omitempty"`
	SignerCN  string    `bson:"signer_common_name,omitempty"`

	BlockchainTxID string    `bson:"blockchain_tx_id,omitempty"`
	BlockchainHash string    `bson:"blockchain_hash,omitempty"`
	BCTimestamp    time.Time `bson:"blockchain_timestamp,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type CreateCertificateRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	CertificateType string `json:"certificate_type" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Issuer          string `json:"issuer" binding:"required"`
	SerialNumber    string `json:"serial_number" binding:"required"`
	RegNo           string `json:"reg_no" binding:"required"`
}

type CertificateResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	StudentCode     string    `json:"student_code"`
	CertificateType string    `json:"certificate_type"`
	Name            string    `json:"name"`
	Issuer          string    `json:"issuer"`
	SerialNumber    string    `json:"serial_number"`
	RegNo           string    `json:"reg_no"`
	Signed          bool      `json:"signed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
