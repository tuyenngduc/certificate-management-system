package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func RegisterAccountRoutes(rg *gin.RouterGroup, h *handler.AccountHandler) {
	account := rg.Group("/accounts")
	account.GET("", h.GetAllAccounts)
}
