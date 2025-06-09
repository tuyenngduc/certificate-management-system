package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	resp, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	userResp, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userResp})
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	var params models.SearchUserParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(400, gin.H{"error": "Tham số không hợp lệ"})
		return
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	users, total, err := h.userService.SearchUsers(c.Request.Context(), params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":       users,
		"total":      total,
		"page":       params.Page,
		"page_size":  params.PageSize,
		"total_page": (total + int64(params.PageSize) - 1) / int64(params.PageSize),
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
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

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Truyền thêm `claims` vào service
	resp, err := h.userService.CreateUser(c.Request.Context(), claims, &req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		case errors.Is(err, common.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		case errors.Is(err, common.ErrStudentIDExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Mã sinh viên đã tồn tại"})
		case errors.Is(err, common.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Email đã tồn tại"})
		case errors.Is(err, common.ErrUniversityNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy trường đại học"})
		case errors.Is(err, common.ErrFacultyNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy khoa hoặc khoa không thuộc trường"})
		default:
			fmt.Printf("CreateUser unexpected error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống, vui lòng thử lại sau"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		}
		return
	}

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không có quyền hoặc token không hợp lệ"})
		return
	}

	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không có quyền hoặc token không hợp lệ"})
		return
	}

	// Tạo context có chứa claims để truyền vào service
	ctx := context.WithValue(c.Request.Context(), "claims", claims)

	err = h.userService.UpdateUser(ctx, id, req)
	if err != nil {
		log.Printf("UpdateUser Error: %v\n", err)
		switch err {
		case common.ErrStudentIDExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Mã sinh viên đã tồn tại"})
		case common.ErrEmailExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Email đã tồn tại"})
		case common.ErrUniversityNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Trường đại học không tồn tại"})
		case common.ErrFacultyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Khoa không tồn tại"})
		case common.ErrUnauthorized, common.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Không có quyền hoặc token không hợp lệ"})
		default:
			if err.Error() == "không có trường nào để cập nhật" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi server"})
			}
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật user thành công"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "Không tìm thấy user"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Xóa user thành công"})
}

func (h *UserHandler) ImportUsersFromExcel(c *gin.Context) {
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

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng upload file Excel"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể mở file"})
		return
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File không đúng định dạng Excel"})
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil || len(rows) == 0 {
		rows, err = f.GetRows("Sheet")
		if err != nil || len(rows) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không đọc được sheet dữ liệu (Sheet1 hoặc Sheet)"})
			return
		}
	}

	var (
		results      []map[string]interface{}
		successCount int
	)

	for i, row := range rows {
		if i == 0 {
			continue // Bỏ qua dòng tiêu đề
		}

		result := map[string]interface{}{"row": i + 1}

		if len(row) < 5 {
			result["error"] = "Thiếu dữ liệu"
			results = append(results, result)
			continue
		}

		user := &models.CreateUserRequest{
			StudentCode: row[0],
			FullName:    row[1],
			Email:       row[2],
			FacultyCode: row[3],
			Course:      row[4],
		}

		_, err := h.userService.CreateUser(c.Request.Context(), claims, user)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrUnauthorized):
				result["error"] = "Bạn chưa đăng nhập hoặc token không hợp lệ"
			case errors.Is(err, common.ErrInvalidToken):
				result["error"] = "Token không hợp lệ"
			case errors.Is(err, common.ErrStudentIDExists):
				result["error"] = "Mã sinh viên đã tồn tại"
			case errors.Is(err, common.ErrEmailExists):
				result["error"] = "Email đã tồn tại"
			case errors.Is(err, common.ErrFacultyNotFound):
				result["error"] = "Không tìm thấy khoa"
			case errors.Is(err, common.ErrUniversityNotFound):
				result["error"] = "Không tìm thấy trường đại học"
			default:
				result["error"] = err.Error()
			}
		} else {
			result["status"] = "created"
			successCount++
		}
		results = append(results, result)
	}

	if successCount == len(results) {
		c.JSON(http.StatusCreated, gin.H{
			"success_count": successCount,
			"results":       results,
		})
	} else {
		c.JSON(http.StatusMultiStatus, gin.H{
			"success_count": successCount,
			"error_count":   len(results) - successCount,
			"results":       results,
		})
	}
}
