package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

type ScoreHandler struct {
	scoreService service.ScoreService
	userSvc      service.UserService
	subjectSvc   service.SubjectService
}

func NewScoreHandler(scoreSvc service.ScoreService) *ScoreHandler {
	return &ScoreHandler{
		scoreService: scoreSvc,
	}
}

// POST /scores - Thêm điểm cho sinh viên
func (h *ScoreHandler) CreateScore(c *gin.Context) {
	var req request.CreateScoreRequest

	// Parse dữ liệu từ body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dữ liệu gửi lên không hợp lệ",
		})
		return
	}

	// Gọi service xử lý logic tạo điểm
	err := h.scoreService.CreateScore(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Thành công
	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo điểm thành công",
	})
}

func (h *ScoreHandler) ImportScoresExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file"})
		return
	}
	defer f.Close()

	xlFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File không phải file Excel hợp lệ"})
		return
	}

	rows, err := xlFile.GetRows("Sheet1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc dữ liệu sheet"})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Excel không có dữ liệu"})
		return
	}

	var results []string
	for i, row := range rows[1:] {
		// Dữ liệu excel giả định theo thứ tự cột:
		// student_code | subject_code | semester | attendance | midterm | final
		if len(row) < 6 {
			results = append(results, "Dòng "+strconv.Itoa(i+2)+" bị thiếu cột")
			continue
		}

		attendance, err1 := parseFloat(row[5])
		midterm, err2 := parseFloat(row[6])
		final, err3 := parseFloat(row[7])
		if err1 != nil || err2 != nil || err3 != nil {
			results = append(results, "Dòng "+strconv.Itoa(i+2)+" điểm không hợp lệ")
			continue
		}

		req := &request.CreateScoreByExcelRequest{
			StudentCode: row[0],
			SubjectCode: row[2],
			Semester:    row[4],
			Attendance:  attendance,
			Midterm:     midterm,
			Final:       final,
		}

		err = h.scoreService.CreateScoreByCode(context.Background(), req)
		if err != nil {
			results = append(results, "Dòng "+strconv.Itoa(i+2)+": "+err.Error())
		} else {
			results = append(results, "Dòng "+strconv.Itoa(i+2)+": thành công")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import hoàn tất",
		"results": results,
	})
}

func parseFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func (h *ScoreHandler) GetScoresByStudent(c *gin.Context) {
	studentID := c.Param("student_id")
	scores, err := h.scoreService.GetScoresByStudentID(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"scores": scores})
}

// Handler
func (h *ScoreHandler) GetScoreDetailByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}

	resp, err := h.scoreService.GetScoreDetailByID(c.Request.Context(), objID)
	if err != nil {
		if err.Error() == "Không tìm thấy sinh viên" {
			c.JSON(404, gin.H{"error": "Không tìm thấy sinh viên"})
			return
		}
		if err.Error() == "Không tìm thấy môn học" {
			c.JSON(404, gin.H{"error": "Không tìm thấy môn học"})
			return
		}
		c.JSON(404, gin.H{"error": "Không tìm thấy điểm"})
		return
	}

	c.JSON(200, resp)
}

func (h *ScoreHandler) GetScoresBySubject(c *gin.Context) {
	subjectID := c.Param("subject_id")
	scores, err := h.scoreService.GetScoresBySubjectID(c.Request.Context(), subjectID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"scores": scores})
}

func (h *ScoreHandler) UpdateScore(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	err := h.scoreService.UpdateScore(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Cập nhật điểm thành công"})
}

func (h *ScoreHandler) ImportScoresBySubjectExcel(c *gin.Context) {
	subjectCode := c.Param("subject_id")
	if subjectCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu mã môn học"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file"})
		return
	}
	defer f.Close()

	xlFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File không phải Excel hợp lệ"})
		return
	}

	rows, err := xlFile.GetRows("Sheet1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc sheet"})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File Excel không có dữ liệu"})
		return
	}

	var reqs []request.ImportScoresBySubjectExcelRequest
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}
		attendance, err1 := strconv.ParseFloat(row[2], 64)
		midterm, err2 := strconv.ParseFloat(row[3], 64)
		final, err3 := strconv.ParseFloat(row[4], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}
		reqs = append(reqs, request.ImportScoresBySubjectExcelRequest{
			StudentID:   row[0],
			StudentName: row[1],
			Attendance:  attendance,
			Midterm:     midterm,
			Final:       final,
			Semester:    row[5],
		})
	}

	results, err := h.scoreService.ImportScoresBySubjectExcel(c.Request.Context(), subjectCode, reqs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import hoàn tất",
		"results": results,
	})
}

func (h *ScoreHandler) GetCGPA(c *gin.Context) {
	studentID := c.Param("student_id")
	resp, err := h.scoreService.CalculateCGPA(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}
