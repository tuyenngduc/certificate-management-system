package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func RegisterSubjectRoutes(r *gin.RouterGroup, subjectHandler *handler.SubjectHandler) {
	r.POST("/subjects", subjectHandler.CreateSubject)
	r.PUT("/subjects/:id", subjectHandler.UpdateSubject)
	r.DELETE("/subjects/:id", subjectHandler.DeleteSubject)
	r.GET("/subjects/:id", subjectHandler.GetSubject)
	r.GET("/subjects", subjectHandler.ListSubjects)
	r.POST("/subjects/import", subjectHandler.ImportSubjects)
}
