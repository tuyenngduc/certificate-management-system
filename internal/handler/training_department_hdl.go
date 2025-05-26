package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingDepartmentHandler struct {
	svc *service.TrainingDepartmentService
}

func NewTrainingDepartmentHandler(svc *service.TrainingDepartmentService) *TrainingDepartmentHandler {
	return &TrainingDepartmentHandler{svc: svc}
}

// Faculty CRUD
func (h *TrainingDepartmentHandler) CreateFaculty(c *gin.Context) {
	var req request.CreateFacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.FacultyValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.CreateFaculty(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "mã khoa đã tồn tại" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tạo khoa thành công"})
}

func (h *TrainingDepartmentHandler) GetAllFaculties(c *gin.Context) {
	faculties, err := h.svc.GetAllFaculties(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, faculties)
}

func (h *TrainingDepartmentHandler) GetFacultyByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	faculty, err := h.svc.GetFacultyByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy khoa"})
		return
	}
	c.JSON(http.StatusOK, faculty)
}

func (h *TrainingDepartmentHandler) UpdateFaculty(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateFacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.FacultyUpdateValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateFaculty(c.Request.Context(), id, &req)
	if err != nil {
		switch err.Error() {
		case "ID không hợp lệ":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case "không tìm thấy khoa":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case "mã khoa đã tồn tại":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		case "không có dữ liệu cập nhật":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật khoa thành công"})
}

func (h *TrainingDepartmentHandler) DeleteFaculty(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.svc.DeleteFaculty(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

// Class CRUD
func (h *TrainingDepartmentHandler) GetClassesByFaculty(c *gin.Context) {
	facultyID := c.Param("faculty_id")
	classes, err := h.svc.GetClassesByFacultyID(c.Request.Context(), facultyID)
	if err != nil {
		if err.Error() == "ID khoa không hợp lệ" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}
func (h *TrainingDepartmentHandler) CreateClass(c *gin.Context) {
	var req request.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.ClassValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.CreateClass(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "mã lớp đã tồn tại" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "ID khoa không hợp lệ" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "ID khoa không tồn tại" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tạo lớp thành công"})
}
func (h *TrainingDepartmentHandler) GetAllClasses(c *gin.Context) {
	classes, err := h.svc.GetAllClasses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *TrainingDepartmentHandler) GetClassByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	class, err := h.svc.GetClassByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy lớp"})
		return
	}
	c.JSON(http.StatusOK, class)
}

func (h *TrainingDepartmentHandler) UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.ClassUpdateValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateClass(c.Request.Context(), id, &req)
	if err != nil {
		switch err.Error() {
		case "ID không hợp lệ":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case "không tìm thấy lớp":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case "mã lớp đã tồn tại":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		case "không có dữ liệu cập nhật":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật lớp thành công"})
}

func (h *TrainingDepartmentHandler) DeleteClass(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.svc.DeleteClass(c.Request.Context(), id); err != nil {
		if err.Error() == "không tìm thấy lớp" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

// Lecturer CRUD
func (h *TrainingDepartmentHandler) CreateLecturer(c *gin.Context) {
	var req request.CreateLecturerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.LecturerValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.CreateLecturer(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "mã giảng viên đã tồn tại":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		case "ID khoa không hợp lệ", "ID khoa không tồn tại":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tạo giảng viên thành công"})
}

func (h *TrainingDepartmentHandler) GetAllLecturers(c *gin.Context) {
	lecturers, err := h.svc.GetAllLecturers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lecturers)
}

func (h *TrainingDepartmentHandler) GetLecturerByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	lecturer, err := h.svc.GetLecturerByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy cán bộ"})
		return
	}
	c.JSON(http.StatusOK, lecturer)
}
func (h *TrainingDepartmentHandler) GetLecturersByFaculty(c *gin.Context) {
	facultyID := c.Param("faculty_id")
	lecturers, err := h.svc.GetLecturersByFacultyID(c.Request.Context(), facultyID)
	if err != nil {
		switch err.Error() {
		case "ID khoa không hợp lệ":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case "khoa không tồn tại":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, lecturers)
}
func (h *TrainingDepartmentHandler) UpdateLecturer(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateLecturerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.LecturerUpdateValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateLecturer(c.Request.Context(), id, &req)
	if err != nil {
		switch err.Error() {
		case "ID không hợp lệ":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case "không tìm thấy cán bộ":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case "mã giảng viên đã tồn tại":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		case "không có dữ liệu cập nhật":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật cán bộ thành công"})
}

func (h *TrainingDepartmentHandler) DeleteLecturer(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := h.svc.DeleteLecturer(c.Request.Context(), id); err != nil {
		if err.Error() == "không tìm thấy giảng viên" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa giảng viên thành công"})
}
