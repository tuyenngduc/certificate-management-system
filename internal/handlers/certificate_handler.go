package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/pkg/database"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateHandler struct {
	certificateService service.CertificateService
	universityService  service.UniversityService
	facultyService     service.FacultyService
	userService        service.UserService
	minioClient        *database.MinioClient
}

func NewCertificateHandler(
	certSvc service.CertificateService,
	uniSvc service.UniversityService,
	facultySvc service.FacultyService,
	userSvc service.UserService,
	minioClient *database.MinioClient,
) *CertificateHandler {
	return &CertificateHandler{
		certificateService: certSvc,
		universityService:  uniSvc,
		facultyService:     facultySvc,
		userService:        userSvc,
		minioClient:        minioClient,
	}
}

func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*utils.CustomClaims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	var req models.CreateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đầu vào không hợp lệ", "chi_tiet": err.Error()})
		return
	}

	res, err := h.certificateService.CreateCertificate(c.Request.Context(), claims, &req)
	if err != nil {
		switch err {
		case common.ErrInvalidUserID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID sinh viên không hợp lệ"})
		case common.ErrUserNotExisted:
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sinh viên"})
		case common.ErrFacultyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy khoa của sinh viên"})
		case common.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		case common.ErrUniversityNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy trường tương ứng"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống", "chi_tiet": err.Error()})
		}
		return
	}

	// Thành công
	c.JSON(http.StatusCreated, gin.H{
		"data": res})
}

func (h *CertificateHandler) GetAllCertificates(c *gin.Context) {
	certs, err := h.certificateService.GetAllCertificates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Lỗi hệ thống",
			"chi_tiet": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": certs,
	})

}

func (h *CertificateHandler) GetCertificateByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID certificate không hợp lệ"})
		return
	}

	cert, err := h.certificateService.GetCertificateByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, common.ErrCertificateNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy certificate"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống", "chi_tiet": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cert,
	})

}

func (h *CertificateHandler) UploadCertificateFile(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*utils.CustomClaims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn file để tải lên"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chỉ hỗ trợ file PDF, JPG, JPEG, PNG"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc file"})
		return
	}

	// Chỉ lấy serial number từ tên file (không còn kiểm tra mã trường)
	serialNumber := strings.TrimSuffix(file.Filename, ext)

	// Lấy university từ token
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ (UniversityID không đúng định dạng)"})
		return
	}
	university, err := h.universityService.GetUniversityByID(c.Request.Context(), universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không lấy được thông tin trường đại học"})
		return
	}

	// Tìm certificate theo serial number
	certificate, err := h.certificateService.GetCertificateBySerialAndUniversity(c.Request.Context(), serialNumber, university.ID)

	if errors.Is(err, mongo.ErrNoDocuments) || certificate == nil || certificate.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng với số serial đã cung cấp"})
		return
	}

	// Kiểm tra xem văn bằng có thuộc trường đang đăng nhập không
	if certificate.UniversityID != university.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không được phép cập nhật văn bằng này"})
		return
	}

	if certificate.Path != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Văn bằng này đã có file, không thể ghi đè"})
		return
	}

	filePath, err := h.certificateService.UploadCertificateFile(c.Request.Context(), certificate.ID, fileData, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tải lên thất bại: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tải file thành công",
		"path":    filePath,
	})
}

func (h *CertificateHandler) GetCertificateFile(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	certificateID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	certificate, err := h.certificateService.GetCertificateByID(ctx, certificateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng"})
		return
	}

	if certificate.Path == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Văn bằng chưa có file"})
		return
	}

	object, err := h.minioClient.Client.GetObject(ctx, h.minioClient.Bucket, certificate.Path, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không đọc được file từ MinIO"})
		return
	}
	defer object.Close()

	fileData, err := io.ReadAll(object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc nội dung file"})
		return
	}

	contentType := http.DetectContentType(fileData)

	c.DataFromReader(http.StatusOK, int64(len(fileData)), contentType, bytes.NewReader(fileData), nil)
}

func (h *CertificateHandler) GetCertificatesByStudentID(c *gin.Context) {
	ctx := c.Request.Context()
	studentIDParam := c.Param("id")

	studentID, err := primitive.ObjectIDFromHex(studentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sinh viên không hợp lệ"})
		return
	}

	// Lấy thông tin sinh viên
	user, err := h.userService.GetUserByID(ctx, studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sinh viên"})
		return
	}

	// Lấy thông tin faculty
	faculty, err := h.facultyService.GetFacultyByCode(ctx, user.FacultyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tìm thấy khoa"})
		return
	}

	// Lấy thông tin university
	university, err := h.universityService.GetUniversityByCode(ctx, user.UniversityCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tìm thấy trường đại học"})
		return
	}

	// Lấy danh sách văn bằng
	certificate, err := h.certificateService.GetCertificateByUserID(ctx, studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng của người dùng"})
		return
	}

	result := models.CertificateResponse{
		ID:              certificate.ID,
		UserID:          certificate.UserID,
		StudentCode:     certificate.StudentCode,
		CertificateType: certificate.CertificateType,
		Name:            certificate.Name,
		SerialNumber:    certificate.SerialNumber,
		RegNo:           certificate.RegNo,
		Path:            certificate.Path,
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		Signed:          certificate.Signed,
		CreatedAt:       certificate.CreatedAt,
		UpdatedAt:       certificate.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *CertificateHandler) SearchCertificates(c *gin.Context) {
	var params models.SearchCertificateParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	certs, total, err := h.certificateService.SearchCertificates(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       certs,
		"total":      total,
		"page":       params.Page,
		"page_size":  params.PageSize,
		"total_page": (total + int64(params.PageSize) - 1) / int64(params.PageSize),
	})
}

func (h *CertificateHandler) GetMyCertificates(c *gin.Context) {
	val, exists := c.Get(string(utils.ClaimsContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		return
	}
	claims, ok := val.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token không hợp lệ"})
		return
	}

	certificates, err := h.certificateService.GetCertificatesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": certificates})
}

func (h *CertificateHandler) GetMyCertificateFile(c *gin.Context) {
	ctx := c.Request.Context()

	// Lấy claims từ context
	claimsRaw, exists := c.Get(string(utils.ClaimsContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin người dùng"})
		return
	}

	claims, ok := claimsRaw.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Thông tin token không hợp lệ"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID trong token không hợp lệ"})
		return
	}

	// Lấy certificate theo user_id
	certificate, err := h.certificateService.GetCertificateByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng cho người dùng"})
		return
	}

	if certificate.Path == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Văn bằng chưa có file"})
		return
	}

	object, err := h.minioClient.Client.GetObject(ctx, h.minioClient.Bucket, certificate.Path, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không đọc được file từ MinIO"})
		return
	}
	defer object.Close()

	fileData, err := io.ReadAll(object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc nội dung file"})
		return
	}

	contentType := http.DetectContentType(fileData)

	c.DataFromReader(http.StatusOK, int64(len(fileData)), contentType, bytes.NewReader(fileData), nil)
}
