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
	DeleteCertificate(ctx context.Context, id primitive.ObjectID) error
	UploadCertificateFile(ctx context.Context, certificateID primitive.ObjectID, fileData []byte, filename string) (string, error)
	GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.CertificateResponse, error)
	GetCertificateBySerialNumber(ctx context.Context, serial string) (*models.Certificate, error)
	GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Certificate, error)
	CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) (*models.CertificateResponse, error)
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
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, common.ErrInvalidUserID
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, common.ErrUserNotExisted
	}
	if user.FacultyID.IsZero() {
		return nil, errors.New("người dùng chưa được gán khoa")
	}
	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	university, err := s.universityRepo.FindByID(ctx, universityID)
	if err != nil || university == nil {
		return nil, common.ErrUniversityNotFound
	}

	cert := &models.Certificate{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		FacultyID:       user.FacultyID,
		UniversityID:    universityID,
		StudentCode:     user.StudentCode,
		CertificateType: req.CertificateType,
		Name:            req.Name,
		SerialNumber:    req.SerialNumber,
		RegNo:           req.RegNo,
		Signed:          false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.certificateRepo.CreateCertificate(ctx, cert)
	if err != nil {
		return nil, err
	}

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     cert.StudentCode,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		Signed:          cert.Signed,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
	}, nil
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

func (s *certificateService) GetCertificateBySerialNumber(ctx context.Context, serial string) (*models.Certificate, error) {
	return s.certificateRepo.FindBySerialNumber(ctx, serial)
}
func (s *certificateService) GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Certificate, error) {
	return s.certificateRepo.FindCertificatesByUserID(ctx, userID)
}

func (s *certificateService) SearchCertificates(ctx context.Context, params models.SearchCertificateParams) ([]*models.CertificateResponse, int64, error) {
	filter := bson.M{}

	if params.StudentCode != "" {
		filter["student_code"] = bson.M{"$regex": params.StudentCode, "$options": "i"}
	}
	if params.CertificateType != "" {
		filter["certificate_type"] = bson.M{"$regex": params.CertificateType, "$options": "i"}
	}
	if params.Signed != nil {
		filter["signed"] = *params.Signed
	}

	// Filter by FacultyCode nếu được truyền vào
	if params.FacultyCode != "" {
		faculty, err := s.facultyRepo.FindByFacultyCode(ctx, params.FacultyCode)
		if err != nil {
			return nil, 0, fmt.Errorf("faculty not found: %w", err)
		}
		if faculty == nil {
			return nil, 0, fmt.Errorf("faculty not found with code: %s", params.FacultyCode)
		}
		filter["faculty_id"] = faculty.ID
	}

	// Truy vấn danh sách certificate
	certs, total, err := s.certificateRepo.FindCertificate(ctx, filter, params.Page, params.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []*models.CertificateResponse
	for _, cert := range certs {
		// Filter theo Course nếu có
		if params.Course != "" {
			user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
			if err != nil || user == nil {
				continue
			}
			if !strings.Contains(strings.ToLower(user.Course), strings.ToLower(params.Course)) {
				continue
			}
		}

		// Lấy Faculty, University để hiển thị thông tin
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
		results = append(results, resp)
	}

	return results, total, nil
}
