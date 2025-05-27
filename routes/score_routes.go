package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func RegisterScoreRoutes(rg *gin.RouterGroup, handler *handler.ScoreHandler) {
	score := rg.Group("/scores")
	{
		score.POST("/", handler.CreateScore)
		score.POST("/import-excel", handler.ImportScoresExcel)
		// Có thể thêm các route khác như:
		// score.GET("/", handler.GetScores)
		// score.GET("/:id", handler.GetScoreByID)
		// score.PUT("/:id", handler.UpdateScore)
		// score.DELETE("/:id", handler.DeleteScore)
	}
}
