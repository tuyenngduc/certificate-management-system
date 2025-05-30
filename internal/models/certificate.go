package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Certificate struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
	CertificateType   string             `bson:"certificate_type" json:"certificate_type"`
	Name              string             `bson:"name" json:"name"`
	Issuer            string             `bson:"issuer" json:"issuer"`
	IssueDate         time.Time          `bson:"issue_date" json:"issue_date"`
	CertificateNumber string             `bson:"certificate_number" json:"certificate_number"`
	Status            string             `bson:"status" json:"status"` // pending, issued, revoked

	BlockchainTxID      string    `bson:"blockchain_tx_id,omitempty" json:"blockchain_tx_id,omitempty"`         // Tx ID trên Hyperledger
	BlockchainHash      string    `bson:"blockchain_hash,omitempty" json:"blockchain_hash,omitempty"`           // Hash của dữ liệu chứng chỉ
	BlockchainTimestamp time.Time `bson:"blockchain_timestamp,omitempty" json:"blockchain_timestamp,omitempty"` // Thời gian ghi blockchain

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CertificateHashData struct {
	UserID            string `json:"user_id"`
	CertificateNumber string `json:"certificate_number"`
	CertificateType   string `json:"certificate_type"`
	Name              string `json:"name"`
	Issuer            string `json:"issuer"`
	IssueDate         string `json:"issue_date"`
	Status            string `json:"status"`
}

func BuildCertificateHashData(cert *Certificate) CertificateHashData {
	return CertificateHashData{
		UserID:            cert.UserID.Hex(),
		CertificateNumber: cert.CertificateNumber,
		CertificateType:   cert.CertificateType,
		Name:              cert.Name,
		Issuer:            cert.Issuer,
		IssueDate:         cert.IssueDate.Format(time.RFC3339), // Định dạng chuẩn
		Status:            cert.Status,
	}
}
