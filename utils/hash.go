package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
)

func ComputeCertificateHash(data models.CertificateHashData) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}
