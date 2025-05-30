package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
)

func RegisterTrainingDepartmentRoutes(rg *gin.RouterGroup, h *handler.TrainingDepartmentHandler) {
	// Faculty
	rg.POST("/faculties", h.CreateFaculty)
	rg.GET("/faculties", h.GetAllFaculties)
	rg.GET("/faculties/:id", h.GetFacultyByID)
	rg.PUT("/faculties/:id", h.UpdateFaculty)
	rg.DELETE("/faculties/:id", h.DeleteFaculty)
	rg.GET("/faculties/classes/:faculty_id", h.GetClassesByFaculty)

	// Class
	rg.POST("/classes", h.CreateClass)
	rg.GET("/classes", h.GetAllClasses)
	rg.GET("/classes/:id", h.GetClassByID)
	rg.PUT("/classes/:id", h.UpdateClass)
	rg.DELETE("/classes/:id", h.DeleteClass)
	rg.GET("/classes/search", h.SearchClasses)

	// Lecturer
	rg.POST("/lecturers", h.CreateLecturer)
	rg.GET("/lecturers", h.GetAllLecturers)
	rg.GET("/lecturers/:id", h.GetLecturerByID)
	rg.PUT("/lecturers/:id", h.UpdateLecturer)
	rg.DELETE("/lecturers/:id", h.DeleteLecturer)
	rg.GET("/lecturers/faculty/:faculty_id", h.GetLecturersByFaculty)
	rg.GET("/lecturers/search", h.SearchLecturers)
}
