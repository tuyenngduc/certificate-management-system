package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainHandler struct {
	BlockchainSvc service.BlockchainService
}

func NewBlockchainHandler(blockchainSvc service.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{BlockchainSvc: blockchainSvc}
}

func (h *BlockchainHandler) PushCertificateToChain(c *gin.Context) {
	certIDStr := c.Param("id")
	certID, err := primitive.ObjectIDFromHex(certIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	txID, err := h.BlockchainSvc.PushCertificateToChain(c.Request.Context(), certID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đưa lên blockchain", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Ghi văn bằng lên blockchain thành công",
		"transaction_id": txID,
		"certificate_id": certID.Hex(),
	})
}

func (h *BlockchainHandler) GetCertificateByID(c *gin.Context) {
	id := c.Param("id")
	result, err := h.BlockchainSvc.GetCertificateFromChain(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
