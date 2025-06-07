package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin xác thực"})
		return
	}

	myClaims, ok := claims.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Thông tin xác thực không hợp lệ"})
		return
	}

	universityID, err := primitive.ObjectIDFromHex(myClaims.UniversityID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "university_id trong token không hợp lệ"})
		return
	}

	faculty, err := h.facultyService.CreateFaculty(c.Request.Context(), &req, universityID)
	if err != nil {
		switch err {
		case common.ErrUniversityNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Trường đại học không tồn tại"})
		case common.ErrFacultyCodeExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Mã khoa đã tồn tại trong trường"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": faculty})
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

	faculties, err := h.facultyService.GetAllFaculties(c.Request.Context(), universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []models.FacultyResponse
	for _, f := range faculties {
		resp = append(resp, models.FacultyResponse{
			ID:           f.ID,
			FacultyCode:  f.FacultyCode,
			FacultyName:  f.FacultyName,
			UniversityID: f.UniversityID,
			CreatedAt:    f.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
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
