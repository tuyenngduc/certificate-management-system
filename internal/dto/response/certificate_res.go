package response

import (
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
)

type CertificateResponse struct {
	ID                 string `json:"id"`
	UserID             string `json:"user_id"`
	CertificateType    string `json:"certificate_type"`
	Name               string `json:"name"`
	Issuer             string `json:"issuer"`
	IssueDate          string `json:"issue_date"`
	SerialNumber       string `json:"serial_number"`
	RegistrationNumber string `json:"registration_number"`
	Status             string `json:"status"`
}

func ToCertificateResponse(cert *models.Certificate) CertificateResponse {
	const layout = "02/01/2006 15:04:05"

	formatTime := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Local().Format(layout)
	}

	return CertificateResponse{
		ID:                 cert.ID.Hex(),
		UserID:             cert.UserID.Hex(),
		CertificateType:    cert.CertificateType,
		Name:               cert.Name,
		Issuer:             cert.Issuer,
		IssueDate:          formatTime(cert.IssueDate),
		SerialNumber:       cert.SerialNumber,
		RegistrationNumber: cert.RegistrationNumber,
		Status:             cert.Status,
	}
}
