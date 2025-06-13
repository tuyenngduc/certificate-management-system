package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/pkg/database"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CertificateService interface {
	GetAllCertificates(ctx context.Context) ([]*models.CertificateResponse, error)
	DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error
	DeleteCertificate(ctx context.Context, id primitive.ObjectID) error
	UploadCertificateFile(ctx context.Context, certificateID primitive.ObjectID, fileData []byte, filename string) (string, error)
	GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.CertificateResponse, error)
	GetCertificateBySerialAndUniversity(ctx context.Context, serial string, universityID primitive.ObjectID) (*models.Certificate, error)
	GetCertificateByUserID(ctx context.Context, userID primitive.ObjectID) (*models.CertificateResponse, error)
	CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) (*models.CertificateResponse, error)
	GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateResponse, error)
	SearchCertificates(ctx context.Context, params models.SearchCertificateParams) ([]*models.CertificateResponse, int64, error)
}

type certificateService struct {
	certificateRepo repository.CertificateRepository
	userRepo        repository.UserRepository
	facultyRepo     repository.FacultyRepository
	universityRepo  repository.UniversityRepository
	minioClient     *database.MinioClient
}

func NewCertificateService(
	certificateRepo repository.CertificateRepository,
	userRepo repository.UserRepository,
	facultyRepo repository.FacultyRepository,
	universityRepo repository.UniversityRepository,
	minioClient *database.MinioClient,
) CertificateService {
	return &certificateService{
		certificateRepo: certificateRepo,
		userRepo:        userRepo,
		facultyRepo:     facultyRepo,
		universityRepo:  universityRepo,
		minioClient:     minioClient,
	}
}

func (s *certificateService) CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) (*models.CertificateResponse, error) {
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	user, err := s.userRepo.FindByStudentCodeAndUniversityID(ctx, req.StudentCode, universityID)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}
	if user.FacultyID.IsZero() {
		return nil, errors.New("người dùng chưa được gán khoa")
	}

	// Kiểm tra thông tin văn bằng hoặc chứng chỉ
	if req.IsDegree {
		if err := s.validateDegreeRequest(ctx, req, universityID); err != nil {
			return nil, err
		}
	} else {
		if err := s.validateCertificateRequest(ctx, req, universityID); err != nil {
			return nil, err
		}
	}

	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	university, err := s.universityRepo.FindByID(ctx, universityID)
	if err != nil || university == nil {
		return nil, common.ErrUniversityNotFound
	}

	// Tạo certificate
	cert := &models.Certificate{
		ID:              primitive.NewObjectID(),
		UserID:          user.ID,
		FacultyID:       user.FacultyID,
		UniversityID:    universityID,
		StudentCode:     user.StudentCode,
		IsDegree:        req.IsDegree,
		Name:            req.Name,
		CertificateType: req.CertificateType,
		SerialNumber:    req.SerialNumber,
		RegNo:           req.RegNo,
		IssueDate:       req.IssueDate,
		Signed:          false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.certificateRepo.CreateCertificate(ctx, cert); err != nil {
		return nil, err
	}

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     cert.StudentCode,
		StudentName:     user.FullName,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		IssueDate:       cert.IssueDate.Format("02/01/2006"),
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		Signed:          cert.Signed,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
	}, nil
}

func (s *certificateService) validateDegreeRequest(ctx context.Context, req *models.CreateCertificateRequest, universityID primitive.ObjectID) error {
	if req.CertificateType == "" || req.SerialNumber == "" || req.RegNo == "" || req.IssueDate.IsZero() {
		return errors.New("thiếu thông tin bắt buộc cho văn bằng (CertificateType, SerialNumber, RegNo, IssueDate)")
	}

	singleDegreeTypes := map[string]bool{
		"Cử nhân": true,
		"Thạc sĩ": true,
		"Tiến sĩ": true,
		"Kỹ sư":   true,
	}

	if singleDegreeTypes[req.CertificateType] {
		alreadyIssued, err := s.certificateRepo.ExistsDegreeByStudentCodeAndType(ctx, req.StudentCode, universityID, req.CertificateType)
		if err != nil {
			return err
		}
		if alreadyIssued {
			return common.ErrCertificateAlreadyExists
		}
	}

	return nil
}

func (s *certificateService) validateCertificateRequest(ctx context.Context, req *models.CreateCertificateRequest, universityID primitive.ObjectID) error {
	if req.Name == "" || req.IssueDate.IsZero() {
		return errors.New("thiếu thông tin bắt buộc cho chứng chỉ (Name, IssueDate)")
	}
	if req.SerialNumber != "" || req.RegNo != "" || req.CertificateType != "" {
		return errors.New("không được truyền SerialNumber, RegNo hoặc CertificateType cho chứng chỉ")
	}

	alreadyIssued, err := s.certificateRepo.ExistsCertificateByStudentCodeAndName(ctx, req.StudentCode, universityID, req.Name)
	if err != nil {
		return err
	}
	if alreadyIssued {
		return common.ErrCertificateAlreadyExists
	}

	return nil
}

func (s *certificateService) GetAllCertificates(ctx context.Context) ([]*models.CertificateResponse, error) {
	// Lấy tất cả certificate từ repository
	certs, err := s.certificateRepo.GetAllCertificates(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.CertificateResponse, 0, len(certs))

	for _, cert := range certs {
		// Lấy faculty
		faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
		if err != nil || faculty == nil {
			// Nếu không tìm thấy faculty thì có thể bỏ qua hoặc gán giá trị mặc định
			faculty = &models.Faculty{
				FacultyCode: "N/A",
				FacultyName: "Không xác định",
			}
		}

		// Lấy university
		university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if err != nil || university == nil {
			university = &models.University{
				UniversityCode: "N/A",
				UniversityName: "Không xác định",
			}
		}

		responses = append(responses, &models.CertificateResponse{
			ID:              cert.ID.Hex(),
			UserID:          cert.UserID.Hex(),
			StudentCode:     cert.StudentCode,
			CertificateType: cert.CertificateType,
			Name:            cert.Name,
			SerialNumber:    cert.SerialNumber,
			RegNo:           cert.RegNo,
			Path:            cert.Path,
			FacultyCode:     faculty.FacultyCode,
			FacultyName:     faculty.FacultyName,
			UniversityCode:  university.UniversityCode,
			UniversityName:  university.UniversityName,
			Signed:          cert.Signed,
			CreatedAt:       cert.CreatedAt,
			UpdatedAt:       cert.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *certificateService) GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.CertificateResponse, error) {
	cert, err := s.certificateRepo.GetCertificateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, common.ErrCertificateNotFound
	}
	user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}

	faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
	if err != nil || faculty == nil {
		faculty = &models.Faculty{
			FacultyCode: "N/A",
			FacultyName: "Không xác định",
		}
	}

	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil || university == nil {
		university = &models.University{
			UniversityCode: "N/A",
			UniversityName: "Không xác định",
		}
	}

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     user.StudentCode,
		StudentName:     user.FullName,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		Path:            cert.Path,
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		Signed:          cert.Signed,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
	}, nil
}

func (s *certificateService) DeleteCertificate(ctx context.Context, id primitive.ObjectID) error {
	return s.certificateRepo.DeleteCertificate(ctx, id)
}

func (s *certificateService) UploadCertificateFile(
	ctx context.Context,
	certificateID primitive.ObjectID,
	fileData []byte,
	filename string,
) (string, error) {
	certificate, err := s.certificateRepo.GetCertificateByID(ctx, certificateID)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy certificate: %w", err)
	}

	university, err := s.universityRepo.FindByID(ctx, certificate.UniversityID)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy trường đại học: %w", err)
	}

	objectKey := fmt.Sprintf("certificates/%s/%s", university.UniversityCode, filename)
	contentType := http.DetectContentType(fileData)

	err = s.minioClient.UploadFile(ctx, objectKey, fileData, contentType)
	if err != nil {
		return "", fmt.Errorf("lỗi upload file lên MinIO: %w", err)
	}

	err = s.certificateRepo.UpdateCertificatePath(ctx, certificateID, objectKey)
	if err != nil {
		return "", fmt.Errorf("lỗi cập nhật path vào MongoDB: %w", err)
	}

	return objectKey, nil
}

func (s *certificateService) GetCertificateBySerialAndUniversity(ctx context.Context, serial string, universityID primitive.ObjectID) (*models.Certificate, error) {
	return s.certificateRepo.FindBySerialAndUniversity(ctx, serial, universityID)
}
func (s *certificateService) GetCertificateByUserID(ctx context.Context, userID primitive.ObjectID) (*models.CertificateResponse, error) {
	cert, err := s.certificateRepo.FindLatestCertificateByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, common.ErrCertificateNotFound
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		user = &models.User{
			FullName: "Không xác định",
		}
	}

	faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
	if err != nil || faculty == nil {
		faculty = &models.Faculty{
			FacultyCode: "N/A",
			FacultyName: "Không xác định",
		}
	}

	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil || university == nil {
		university = &models.University{
			UniversityCode: "N/A",
			UniversityName: "Không xác định",
		}
	}

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     cert.StudentCode,
		StudentName:     user.FullName,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		Path:            cert.Path,
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		Signed:          cert.Signed,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
	}, nil
}

func (s *certificateService) SearchCertificates(ctx context.Context, params models.SearchCertificateParams) ([]*models.CertificateResponse, int64, error) {
	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		return nil, 0, common.ErrUnauthorized
	}

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, 0, common.ErrInvalidToken
	}

	filter := bson.M{
		"university_id": universityID,
	}

	if params.StudentCode != "" {
		filter["student_code"] = bson.M{"$regex": params.StudentCode, "$options": "i"}
	}
	if params.CertificateType != "" {
		filter["certificate_type"] = bson.M{"$regex": params.CertificateType, "$options": "i"}
	}
	if params.Signed != nil {
		filter["signed"] = *params.Signed
	}
	if params.FacultyCode != "" {
		faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, params.FacultyCode, universityID)
		if err != nil || faculty == nil {
			return nil, 0, fmt.Errorf("faculty not found in your university with code: %s", params.FacultyCode)
		}
		filter["faculty_id"] = faculty.ID
	}

	certs, total, err := s.certificateRepo.FindCertificate(ctx, filter, params.Page, params.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []*models.CertificateResponse
	for _, cert := range certs {
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
		if err != nil || user == nil {
			continue
		}

		if params.Course != "" && !strings.Contains(strings.ToLower(user.Course), strings.ToLower(params.Course)) {
			continue
		}

		faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
		if err != nil || faculty == nil {
			continue
		}

		university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if err != nil || university == nil {
			continue
		}

		resp := &models.CertificateResponse{
			ID:              cert.ID.Hex(),
			UserID:          cert.UserID.Hex(),
			StudentCode:     cert.StudentCode,
			StudentName:     user.FullName,
			CertificateType: cert.CertificateType,
			Name:            cert.Name,
			SerialNumber:    cert.SerialNumber,
			RegNo:           cert.RegNo,
			Path:            cert.Path,
			IssueDate:       cert.IssueDate.Format("02/01/2006"),
			FacultyCode:     faculty.FacultyCode,
			FacultyName:     faculty.FacultyName,
			UniversityCode:  university.UniversityCode,
			UniversityName:  university.UniversityName,
			Signed:          cert.Signed,
			CreatedAt:       cert.CreatedAt,
			UpdatedAt:       cert.UpdatedAt,
		}
		results = append(results, resp)
	}

	return results, total, nil
}

func (s *certificateService) GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateResponse, error) {
	certs, err := s.certificateRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []*models.CertificateResponse
	for _, cert := range certs {
		// Lấy user
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
		if err != nil || user == nil {
			continue
		}

		// Lấy khoa
		faculty, _ := s.facultyRepo.FindByID(ctx, cert.FacultyID)
		if faculty == nil {
			faculty = &models.Faculty{
				FacultyCode: "N/A",
				FacultyName: "Không xác định",
			}
		}

		// Lấy trường
		university, _ := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if university == nil {
			university = &models.University{
				UniversityCode: "N/A",
				UniversityName: "Không xác định",
			}
		}

		resp := &models.CertificateResponse{
			ID:              cert.ID.Hex(),
			UserID:          cert.UserID.Hex(),
			StudentCode:     user.StudentCode,
			StudentName:     user.FullName,
			CertificateType: cert.CertificateType,
			Name:            cert.Name,
			SerialNumber:    cert.SerialNumber,
			RegNo:           cert.RegNo,
			Path:            cert.Path,
			FacultyCode:     faculty.FacultyCode,
			FacultyName:     faculty.FacultyName,
			UniversityCode:  university.UniversityCode,
			UniversityName:  university.UniversityName,
			Signed:          cert.Signed,
			CreatedAt:       cert.CreatedAt,
			UpdatedAt:       cert.UpdatedAt,
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *certificateService) DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error {
	return s.certificateRepo.DeleteCertificateByID(ctx, id)
}
