package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/xuri/excelize/v2"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dob, err := time.Parse("02/01/2006", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sai định dạng ngày sinh (dd/mm/yyyy)"})
		return
	}

	user := &models.User{
		FullName:     req.FullName,
		StudentID:    req.StudentID,
		Email:        req.Email,
		Ethnicity:    req.Ethnicity,
		Gender:       req.Gender,
		Major:        req.Major,
		Class:        req.Class,
		Course:       req.Course,
		NationalID:   req.NationalID,
		Address:      req.Address,
		PlaceOfBirth: req.PlaceOfBirth,
		DateOfBirth:  dob,
		PhoneNumber:  req.PhoneNumber,
	}

	err = h.svc.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tạo user thành công"})
}

func (h *UserHandler) BulkCreateUser(c *gin.Context) {
	var req request.BulkCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var results []map[string]interface{}
	for _, u := range req.Users {
		// Parse date, validate, map sang models.User như CreateUser
		dob, err := time.Parse("02/01/2006", u.DateOfBirth)
		if err != nil {
			results = append(results, map[string]interface{}{"email": u.Email, "error": "Sai định dạng ngày sinh"})
			continue
		}
		user := &models.User{
			FullName:     u.FullName,
			Email:        u.Email,
			Ethnicity:    u.Ethnicity,
			Gender:       u.Gender,
			Major:        u.Major,
			Class:        u.Class,
			Course:       u.Course,
			NationalID:   u.NationalID,
			Address:      u.Address,
			PlaceOfBirth: u.PlaceOfBirth,
			DateOfBirth:  dob,
			PhoneNumber:  u.PhoneNumber,
		}
		err = h.svc.CreateUser(c.Request.Context(), user)
		if err != nil {
			results = append(results, map[string]interface{}{"email": u.Email, "error": err.Error()})
		} else {
			results = append(results, map[string]interface{}{"email": u.Email, "status": "created"})
		}
	}
	c.JSON(http.StatusCreated, gin.H{"results": results})
}

func (h *UserHandler) ImportUsersFromExcel(c *gin.Context) {
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

	var results []map[string]interface{}
	for i, row := range rows {
		if i == 0 {
			continue // Bỏ qua dòng tiêu đề
		}
		if len(row) < 12 {
			results = append(results, map[string]interface{}{"row": i + 1, "error": "Thiếu dữ liệu"})
			continue
		}
		dob, err := time.Parse("02/01/2006", row[11])
		if err != nil {
			results = append(results, map[string]interface{}{"row": i + 1, "error": "Sai định dạng ngày sinh"})
			continue
		}
		user := &models.User{
			FullName:     row[0],
			Email:        row[1],
			StudentID:    row[2],
			Ethnicity:    row[3],
			Gender:       row[4],
			Major:        row[5],
			Class:        row[6],
			Course:       row[7],
			NationalID:   row[8],
			Address:      row[9],
			PlaceOfBirth: row[10],
			DateOfBirth:  dob,
			PhoneNumber:  row[11],
		}
		err = h.svc.CreateUser(c.Request.Context(), user)
		if err != nil {
			results = append(results, map[string]interface{}{"row": i + 1, "error": err.Error()})
		} else {
			results = append(results, map[string]interface{}{"row": i + 1, "status": "created"})
		}
	}
	c.JSON(http.StatusCreated, gin.H{"results": results})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	fullName := c.Query("full_name")
	email := c.Query("email")
	nationalID := c.Query("national_id")
	phone := c.Query("phone_number")
	studentID := c.Query("student_id")

	users, err := h.svc.SearchUsers(c.Request.Context(), fullName, email, nationalID, phone, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.svc.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật user thành công"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa user thành công"})
}
