package request

import "time"

type CreateCertificateRequest struct {
	UserID             string    `json:"user_id" binding:"required"`
	CertificateType    string    `json:"certificate_type" binding:"required"`
	Name               string    `json:"name" binding:"required"`
	Issuer             string    `json:"issuer" binding:"required"`
	IssueDate          time.Time `json:"issue_date" binding:"required"`
	SerialNumber       string    `json:"serial_number" binding:"required"`
	RegistrationNumber string    `json:"registration_number" binding:"required"`
}
