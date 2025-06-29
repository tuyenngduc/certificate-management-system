package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FacultyHandler struct {
	facultyService service.FacultyService
}

func NewFacultyHandler(facultyService service.FacultyService) *FacultyHandler {
	return &FacultyHandler{facultyService: facultyService}
}

func (h *FacultyHandler) CreateFaculty(c *gin.Context) {
	var req models.CreateFacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		return
	}
	claims, ok := val.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	err := h.facultyService.CreateFaculty(c.Request.Context(), claims, &req)
	if err != nil {
		switch err {
		case common.ErrUniversityNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Trường đại học không tồn tại"})
		case common.ErrFacultyCodeExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Mã khoa đã tồn tại trong trường"})
		case common.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ hoặc không chứa thông tin trường đại học"})
		default:
			log.Printf("Internal server error when creating faculty: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Đã xảy ra lỗi hệ thống, vui lòng thử lại sau."})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tạo khoa thành công!"})
}

func (h *FacultyHandler) GetAllFaculties(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin đăng nhập"})
		return
	}

	userClaims := claims.(*utils.CustomClaims)
	universityID, err := primitive.ObjectIDFromHex(userClaims.UniversityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "university_id không hợp lệ trong token"})
		return
	}

	resp, err := h.facultyService.GetAllFaculties(c.Request.Context(), universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *FacultyHandler) GetFacultyByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	faculty, err := h.facultyService.GetFacultyByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Khoa không tồn tại"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": faculty})
}

func (h *FacultyHandler) UpdateFaculty(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateFacultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	faculty, err := h.facultyService.UpdateFaculty(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case common.ErrInvalidUserID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID khoa không hợp lệ"})
		case common.ErrFacultyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Khoa không tồn tại"})
		case common.ErrFacultyCodeExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Mã khoa đã tồn tại"})
		case common.ErrNoFieldsToUpdate:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không có trường nào để cập nhật"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": faculty})
}

func (h *FacultyHandler) DeleteFaculty(c *gin.Context) {
	id := c.Param("id")

	err := h.facultyService.DeleteFaculty(c.Request.Context(), id)
	if err != nil {
		switch err {
		case common.ErrInvalidUserID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID khoa không hợp lệ"})
		case common.ErrFacultyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Khoa không tồn tại"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa khoa thành công"})
}
