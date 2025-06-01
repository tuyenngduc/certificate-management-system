package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
)

func HashCertificateData(cert *models.Certificate) (string, error) {
	hashData := models.BuildCertificateHashData(cert)
	jsonBytes, err := json.Marshal(hashData)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}
