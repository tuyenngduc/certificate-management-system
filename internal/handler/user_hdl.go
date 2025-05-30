package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/response"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.ValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(400, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dob, err := time.Parse("02/01/2006", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sai định dạng ngày sinh (dd/mm/yyyy)"})
		return
	}

	// Tìm faculty theo code
	faculty, err := h.svc.FindFacultyByCode(c.Request.Context(), req.FacultyCode)
	if err != nil || faculty == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy khoa với mã " + req.FacultyCode})
		return
	}
	// Tìm class theo code
	class, err := h.svc.FindClassByCode(c.Request.Context(), req.ClassCode)
	if err != nil || class == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy lớp với mã " + req.ClassCode})
		return
	}

	user := &models.User{
		FullName:     req.FullName,
		StudentID:    req.StudentID,
		Email:        req.Email,
		Ethnicity:    req.Ethnicity,
		Gender:       req.Gender,
		FacultyID:    faculty.ID,
		ClassID:      class.ID,
		Course:       req.Course,
		NationalID:   req.NationalID,
		Address:      req.Address,
		PlaceOfBirth: req.PlaceOfBirth,
		DateOfBirth:  dob,
		PhoneNumber:  req.PhoneNumber,
	}

	err = h.svc.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	var (
		results      []map[string]interface{}
		successCount int
	)

	for _, u := range req.Users {
		result := map[string]interface{}{"email": u.Email}

		dob, err := time.Parse("02/01/2006", u.DateOfBirth)
		if err != nil {
			result["error"] = "Sai định dạng ngày sinh (dd/mm/yyyy)"
			results = append(results, result)
			continue
		}

		faculty, err := h.svc.FindFacultyByCode(c.Request.Context(), u.FacultyCode)
		if err != nil || faculty == nil {
			result["error"] = "Không tìm thấy khoa với mã " + u.FacultyCode
			results = append(results, result)
			continue
		}

		class, err := h.svc.FindClassByCode(c.Request.Context(), u.ClassCode)
		if err != nil || class == nil {
			result["error"] = "Không tìm thấy lớp với mã " + u.ClassCode
			results = append(results, result)
			continue
		}

		user := &models.User{
			FullName:     u.FullName,
			StudentID:    u.StudentID,
			Email:        u.Email,
			Ethnicity:    u.Ethnicity,
			Gender:       u.Gender,
			FacultyID:    faculty.ID,
			ClassID:      class.ID,
			Course:       u.Course,
			NationalID:   u.NationalID,
			Address:      u.Address,
			PlaceOfBirth: u.PlaceOfBirth,
			DateOfBirth:  dob,
			PhoneNumber:  u.PhoneNumber,
		}

		if err := h.svc.CreateUser(c.Request.Context(), user); err != nil {
			result["error"] = err.Error()
		} else {
			result["status"] = "created"
			successCount++
		}

		results = append(results, result)
	}

	if successCount == len(req.Users) {
		c.JSON(http.StatusCreated, gin.H{
			"success_count": successCount,
			"results":       results,
		})
	} else {
		c.JSON(207, gin.H{
			"success_count": successCount,
			"error_count":   len(req.Users) - successCount,
			"results":       results,
		})
	}
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

	var (
		results      []map[string]interface{}
		successCount int
	)

	for i, row := range rows {
		if i == 0 {
			continue // Bỏ qua header
		}

		result := map[string]interface{}{"row": i + 1}

		if len(row) < 13 {
			result["error"] = "Thiếu dữ liệu"
			results = append(results, result)
			continue
		}

		dob, err := time.Parse("02/01/2006", row[11])
		if err != nil {
			result["error"] = "Sai định dạng ngày sinh (dd/mm/yyyy)"
			results = append(results, result)
			continue
		}

		faculty, err := h.svc.FindFacultyByCode(c.Request.Context(), row[5])
		if err != nil || faculty == nil {
			result["error"] = "Không tìm thấy khoa với mã " + row[5]
			results = append(results, result)
			continue
		}
		class, err := h.svc.FindClassByCode(c.Request.Context(), row[6])
		if err != nil || class == nil {
			result["error"] = "Không tìm thấy lớp với mã " + row[6]
			results = append(results, result)
			continue
		}

		user := &models.User{
			FullName:     row[0],
			Email:        row[1],
			StudentID:    row[2],
			Ethnicity:    row[3],
			Gender:       row[4],
			FacultyID:    faculty.ID,
			ClassID:      class.ID,
			Course:       row[7],
			NationalID:   row[8],
			Address:      row[9],
			PlaceOfBirth: row[10],
			DateOfBirth:  dob,
			PhoneNumber:  row[12],
		}

		if err := h.svc.CreateUser(c.Request.Context(), user); err != nil {
			result["error"] = err.Error()
		} else {
			result["status"] = "created"
			successCount++
		}
		results = append(results, result)
	}

	// Trả về HTTP code phù hợp
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

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	faculties, _ := h.svc.GetAllFaculties(c.Request.Context())
	classes, _ := h.svc.GetAllClasses(c.Request.Context())

	facultyMap := make(map[primitive.ObjectID]string)
	for _, f := range faculties {
		facultyMap[f.ID] = f.Code
	}
	classMap := make(map[primitive.ObjectID]string)
	for _, cl := range classes {
		classMap[cl.ID] = cl.Code
	}

	var resp []*response.UserResponse
	for _, u := range users {
		resp = append(resp, &response.UserResponse{
			ID:           u.ID.Hex(),
			StudentID:    u.StudentID,
			FullName:     u.FullName,
			Email:        u.Email,
			Ethnicity:    u.Ethnicity,
			Gender:       u.Gender,
			FacultyCode:  facultyMap[u.FacultyID],
			ClassCode:    classMap[u.ClassID],
			Course:       u.Course,
			NationalID:   u.NationalID,
			Address:      u.Address,
			PlaceOfBirth: u.PlaceOfBirth,
			DateOfBirth:  u.DateOfBirth.Format("02/01/2006"),
			PhoneNumber:  u.PhoneNumber,
		})
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	id := c.Query("id")
	fullName := c.Query("full_name")
	email := c.Query("email")
	nationalID := c.Query("national_id")
	phone := c.Query("phone_number")
	studentID := c.Query("student_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.svc.SearchUsers(c.Request.Context(), id, fullName, email, nationalID, phone, studentID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Lấy map facultyID -> code, classID -> code
	faculties, _ := h.svc.GetAllFaculties(c.Request.Context())
	classes, _ := h.svc.GetAllClasses(c.Request.Context())
	facultyMap := make(map[primitive.ObjectID]string)
	for _, f := range faculties {
		facultyMap[f.ID] = f.Code
	}
	classMap := make(map[primitive.ObjectID]string)
	for _, cl := range classes {
		classMap[cl.ID] = cl.Code
	}

	var result []response.UserResponse
	for _, u := range users {
		result = append(result, response.UserResponse{
			ID:           u.ID.Hex(),
			StudentID:    u.StudentID,
			FullName:     u.FullName,
			Email:        u.Email,
			Ethnicity:    u.Ethnicity,
			Gender:       u.Gender,
			FacultyCode:  facultyMap[u.FacultyID],
			ClassCode:    classMap[u.ClassID],
			Course:       u.Course,
			NationalID:   u.NationalID,
			Address:      u.Address,
			PlaceOfBirth: u.PlaceOfBirth,
			DateOfBirth:  u.DateOfBirth.Format("02/01/2006"),
			PhoneNumber:  u.PhoneNumber,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       result,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				if msg, ok := request.ValidateMessages[fe.Field()][fe.Tag()]; ok {
					c.JSON(400, gin.H{"error": msg})
					return
				}
			}
		}
		c.JSON(400, gin.H{"error": err.Error()})
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
		if err.Error() == "user không tồn tại" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa user thành công"})
}
