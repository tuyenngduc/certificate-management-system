package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/response"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/xuri/excelize/v2"
)

type SubjectHandler struct {
	subjectService            service.SubjectService
	trainingDepartmentService *service.TrainingDepartmentService
}

func NewSubjectHandler(subjectService service.SubjectService, trainingDepartmentService *service.TrainingDepartmentService) *SubjectHandler {
	return &SubjectHandler{subjectService: subjectService, trainingDepartmentService: trainingDepartmentService}
}

// POST /subjects
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var req request.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.subjectService.CreateSubject(c.Request.Context(), &req)
	if err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "mã môn học đã tồn tại":
			c.JSON(http.StatusConflict, gin.H{"error": errMsg})
		case "khoa không tồn tại":
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		case "id khoa không hợp lệ":
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "đã tạo môn học thành công"})
}

func (h *SubjectHandler) SearchSubjects(c *gin.Context) {
	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")
	creditStr := c.Query("credit")

	var credit *int
	if creditStr != "" {
		val, err := strconv.Atoi(creditStr)
		if err == nil {
			credit = &val
		}
	}

	subjects, err := h.subjectService.Search(c.Request.Context(), id, code, name, credit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, subjects)
}

// PUT /subjects/:id
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.subjectService.UpdateSubject(c.Request.Context(), id, &req); err != nil {
		errMsg := err.Error()
		switch {
		case errMsg == "id môn học không hợp lệ" || errMsg == "id khoa không hợp lệ":
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		case errMsg == "môn học không tồn tại" || errMsg == "khoa không tồn tại":
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		case errMsg == "mã môn học đã tồn tại":
			c.JSON(http.StatusConflict, gin.H{"error": errMsg}) // 👈 Thêm dòng này
		case errMsg == "không có dữ liệu cập nhật":
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "đã cập nhật môn học thành công"})
}

// DELETE /subjects/:id
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	id := c.Param("id")
	if err := h.subjectService.DeleteSubject(c.Request.Context(), id); err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "id môn học không hợp lệ":
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		case "không tìm thấy môn học":
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "đã xóa môn học thành công"})
}

// GET /subjects/:id
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	id := c.Param("id")
	subject, err := h.subjectService.GetSubjectByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if subject == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "không tìm thấy môn học"})
		return
	}

	// Lấy tên khoa từ FacultyID
	faculty, err := h.trainingDepartmentService.GetFacultyByID(c.Request.Context(), subject.FacultyID)
	facultyName := ""
	if err == nil && faculty != nil {
		facultyName = faculty.Name
	}

	resp := response.SubjectResponse{
		ID:          subject.ID.Hex(),
		Code:        subject.Code,
		Name:        subject.Name,
		Credit:      subject.Credit,
		FacultyID:   facultyName, // trả về tên khoa
		Description: subject.Description,
	}

	c.JSON(http.StatusOK, resp)
}

// GET /subjects
func (h *SubjectHandler) ListSubjects(c *gin.Context) {
	subjects, err := h.subjectService.ListSubjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Lấy danh sách khoa và tạo map facultyID -> facultyName
	faculties, err := h.trainingDepartmentService.GetAllFaculties(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không lấy được danh sách khoa"})
		return
	}
	facultyMap := make(map[string]string)
	for _, f := range faculties {
		facultyMap[f.ID.Hex()] = f.Name
	}

	var result []response.SubjectResponse
	for _, s := range subjects {
		result = append(result, response.SubjectResponse{
			ID:          s.ID.Hex(),
			Code:        s.Code,
			Name:        s.Name,
			Credit:      s.Credit,
			FacultyID:   facultyMap[s.FacultyID.Hex()], // Trả về tên khoa thay vì id
			Description: s.Description,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *SubjectHandler) ImportSubjects(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy file"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không mở được file"})
		return
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File không đúng định dạng Excel"})
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không đọc được sheet"})
		return
	}

	type ImportResult struct {
		Row     int    `json:"row"`
		Code    string `json:"code"`
		Status  string `json:"status"`  // "success" hoặc "failed"
		Message string `json:"message"` // Lý do thành công/thất bại
	}

	var results []ImportResult

	for i, row := range rows {
		if i == 0 {
			continue // Bỏ qua header
		}
		result := ImportResult{
			Row:    i + 1,
			Code:   "",
			Status: "failed",
		}
		if len(row) < 4 {
			result.Message = "Thiếu dữ liệu bắt buộc"
			results = append(results, result)
			continue
		}
		code := row[0]
		creditStr := row[2]
		credit, err := strconv.Atoi(creditStr)
		result.Code = code

		if err != nil {
			result.Message = "Credit không hợp lệ"
			results = append(results, result)
			continue
		}

		req := request.CreateSubjectByExcelRequest{
			Code:        code,
			Name:        row[1],
			Credit:      credit,
			FacultyCode: row[3],
		}
		if len(row) > 4 {
			req.Description = row[4]
		}

		if err := h.subjectService.CreateSubjectByFacultyCode(c.Request.Context(), &req); err == nil {
			result.Status = "success"
			result.Message = "Thành công"
		} else {
			result.Message = err.Error()
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import hoàn tất",
		"results": results,
	})
}
