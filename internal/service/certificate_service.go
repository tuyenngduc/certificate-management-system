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
	// Parse universityID từ token
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	// Lấy user theo StudentCode và UniversityID (ràng buộc trường)
	user, err := s.userRepo.FindByStudentCodeAndUniversityID(ctx, req.StudentCode, universityID)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}

	if user.FacultyID.IsZero() {
		return nil, errors.New("người dùng chưa được gán khoa")
	}

	// Lấy thông tin khoa
	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	// Lấy thông tin trường
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
		CertificateType: req.CertificateType,
		Name:            req.Name,
		SerialNumber:    req.SerialNumber,
		RegNo:           req.RegNo,
		Signed:          false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Lưu vào DB
	err = s.certificateRepo.CreateCertificate(ctx, cert)
	if err != nil {
		return nil, err
	}

	// Trả kết quả
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

func (s *certificateService) SearchCertificates(ctx context.Context, params models.SearchCertificateParams) ([]*models.CertificateResponse, int64, error) {
	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		fmt.Println("DEBUG: Claims not found or invalid")
		return nil, 0, common.ErrUnauthorized
	}
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		fmt.Println("DEBUG: Invalid universityID in token:", claims.UniversityID)
		return nil, 0, common.ErrInvalidToken
	}
	fmt.Println("DEBUG: universityID from token:", universityID.Hex())

	filter := bson.M{}
	filter["university_id"] = universityID

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
		fmt.Println("DEBUG: Searching with FacultyCode:", params.FacultyCode)

		faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, params.FacultyCode, universityID)
		if err != nil {
			fmt.Println("DEBUG: Error when finding faculty:", err)
			return nil, 0, fmt.Errorf("faculty not found: %w", err)
		}
		if faculty == nil {
			fmt.Println("DEBUG: Faculty with code", params.FacultyCode, "not found in university", universityID.Hex())
			return nil, 0, fmt.Errorf("faculty not found in your university with code: %s", params.FacultyCode)
		}
		fmt.Println("DEBUG: Found faculty ID:", faculty.ID.Hex(), "with code:", faculty.FacultyCode)
		filter["faculty_id"] = faculty.ID
	}

	fmt.Println("DEBUG: Final MongoDB filter:", filter)

	certs, total, err := s.certificateRepo.FindCertificate(ctx, filter, params.Page, params.PageSize)
	if err != nil {
		fmt.Println("DEBUG: Error when finding certificates:", err)
		return nil, 0, err
	}
	fmt.Println("DEBUG: Total certs found:", total)

	var results []*models.CertificateResponse
	for _, cert := range certs {
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
		if err != nil || user == nil {
			fmt.Println("DEBUG: Could not find user for cert ID:", cert.ID.Hex())
			continue
		}

		if params.Course != "" && !strings.Contains(strings.ToLower(user.Course), strings.ToLower(params.Course)) {
			fmt.Println("DEBUG: Course mismatch. Expected contains:", params.Course, "Actual:", user.Course)
			continue
		}

		faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
		if err != nil || faculty == nil {
			fmt.Println("DEBUG: Could not find faculty for cert ID:", cert.ID.Hex())
			continue
		}
		university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if err != nil || university == nil {
			fmt.Println("DEBUG: Could not find university for cert ID:", cert.ID.Hex())
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
