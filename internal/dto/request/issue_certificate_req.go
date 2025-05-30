package request

type CreateCertificateRequest struct {
	UserID            string `json:"user_id" binding:"required"`
	CertificateType   string `json:"certificate_type" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Issuer            string `json:"issuer" binding:"required"`
	IssueDate         string `json:"issue_date" binding:"required"` // RFC3339
	CertificateNumber string `json:"certificate_number" binding:"required"`
}
