package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func RegisterUserRoutes(rg *gin.RouterGroup, handler *handler.UserHandler) {
	rg.POST("/users", handler.CreateUser)
	rg.POST("/users/bulk", handler.BulkCreateUser)
	rg.POST("/users/import-excel", handler.ImportUsersFromExcel)
	rg.GET("/users", handler.GetAllUsers)
	rg.GET("/users/search", handler.SearchUsers)
	rg.GET("/users/:id", handler.GetUserByID)
	rg.PUT("/users/:id", handler.UpdateUser)
	rg.DELETE("/users/:id", handler.DeleteUser)
	rg.GET("/users/class/:class_id", handler.GetUsersByClassID)

}
