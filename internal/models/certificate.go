package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Certificate struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	CertificateType string             `bson:"certificate_type" json:"certificate_type"`
	Name            string             `bson:"name" json:"name"`
	Issuer          string             `bson:"issuer" json:"issuer"`
	IssueDate       time.Time          `bson:"issue_date" json:"issue_date"`
	SerialNumber    string             `bson:"serial_number" json:"serial_number"`

	RegistrationNumber string `bson:"registration_number,omitempty" json:"registration_number,omitempty"` // Số vào sổ cấp văn bằng
	Status             string `bson:"status" json:"status"`

	// Blockchain
	BlockchainTxID      string    `bson:"blockchain_tx_id,omitempty" json:"blockchain_tx_id,omitempty"`
	BlockchainHash      string    `bson:"blockchain_hash,omitempty" json:"blockchain_hash,omitempty"`
	BlockchainTimestamp time.Time `bson:"blockchain_timestamp,omitempty" json:"blockchain_timestamp,omitempty"`

	Hash              string    `bson:"hash,omitempty" json:"hash,omitempty"`
	Signature         string    `bson:"signature,omitempty" json:"signature,omitempty"`
	Signed            bool      `bson:"signed" json:"signed"`
	SignedAt          time.Time `bson:"signed_at,omitempty" json:"signed_at,omitempty"`
	SignerCommonName  string    `bson:"signer_common_name,omitempty" json:"signer_common_name,omitempty"`
	SignerCertificate string    `bson:"signer_certificate,omitempty" json:"signer_certificate,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
