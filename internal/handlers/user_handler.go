package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
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
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var resp []models.UserResponse
	for _, u := range users {
		resp = append(resp, models.UserResponse{
			ID:        u.ID,
			StudentID: u.StudentID,
			FullName:  u.FullName,
			Email:     u.Email,
			Faculty:   u.Faculty,
			Class:     u.Class,
			Course:    u.Course,
			Status:    u.Status,
		})
	}

	c.JSON(200, gin.H{"data": resp})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Không tìm thấy user"})
		return
	}
	resp := models.UserResponse{
		ID:        user.ID,
		StudentID: user.StudentID,
		FullName:  user.FullName,
		Email:     user.Email,
		Faculty:   user.Faculty,
		Class:     user.Class,
		Course:    user.Course,
		Status:    user.Status,
	}
	c.JSON(200, gin.H{"data": resp})
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

	var resp []models.UserResponse
	for _, u := range users {
		resp = append(resp, models.UserResponse{
			ID:        u.ID,
			StudentID: u.StudentID,
			FullName:  u.FullName,
			Email:     u.Email,
			Faculty:   u.Faculty,
			Class:     u.Class,
			Course:    u.Course,
			Status:    u.Status,
		})
	}

	c.JSON(200, gin.H{
		"data":       resp,
		"total":      total,
		"page":       params.Page,
		"page_size":  params.PageSize,
		"total_page": (total + int64(params.PageSize) - 1) / int64(params.PageSize),
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if userErr := common.ParseMongoDuplicateError(err); userErr != "" {
			c.JSON(http.StatusConflict, gin.H{"error": userErr})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(400, gin.H{"errors": errs})
			return
		}
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	err = h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			c.JSON(404, gin.H{"error": "Không tìm thấy user"})
		case common.ErrStudentIDExists:
			c.JSON(400, gin.H{"error": "Mã sinh viên đã được sử dụng"})
		case common.ErrEmailExists:
			c.JSON(400, gin.H{"error": "Email đã được sử dụng"})
		default:
			c.JSON(400, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Cập nhật thành công"})
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
