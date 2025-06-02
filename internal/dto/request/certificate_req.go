package request

type CreateCertificateRequest struct {
	UserID             string `json:"user_id" binding:"required"`
	CertificateType    string `json:"certificate_type" binding:"required,oneof=degree certificate"`
	Name               string `json:"name" binding:"required"`
	SerialNumber       string `json:"serial_number" binding:"required"`
	RegistrationNumber string `json:"registration_number" binding:"required"`
}
