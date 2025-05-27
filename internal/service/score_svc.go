package service

import (
	"context"
	"errors"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreService interface {
	CreateScore(ctx context.Context, req *request.CreateScoreRequest) error
	CreateScoreByCode(ctx context.Context, req *request.CreateScoreByExcelRequest) error
}

type scoreService struct {
	repo        *repository.ScoreRepository
	userRepo    *repository.UserRepository
	subjectRepo repository.SubjectRepository
}

func NewScoreService(repo *repository.ScoreRepository, userRepo *repository.UserRepository, subjectRepo repository.SubjectRepository) ScoreService {
	return &scoreService{repo: repo, userRepo: userRepo, subjectRepo: subjectRepo}
}

// Tính điểm quá trình = attendance*0.3 + midterm*0.7
func calculateProcessScore(attendance, midterm float64) float64 {
	return attendance*0.3 + midterm*0.7
}

// Tính điểm tổng kết = processScore*0.3 + final*0.7
func calculateTotalScore(processScore, final float64) float64 {
	return processScore*0.3 + final*0.7
}

func convertGPAChar(score float64) string {
	switch {
	case score >= 9.0:
		return "A+"
	case score >= 8.5 && score <= 8.9:
		return "A"
	case score >= 7.8 && score <= 8.4:
		return "B+"
	case score >= 6.3 && score <= 6.9:
		return "C+"
	case score >= 5.5 && score <= 6.2:
		return "C"
	case score >= 4.8 && score <= 5.4:
		return "D+"
	case score >= 4.0 && score <= 4.7:
		return "D"
	default:
		return "F"
	}
}

func (s *scoreService) CreateScore(ctx context.Context, req *request.CreateScoreRequest) error {
	// Validate ObjectID
	studentID, err := primitive.ObjectIDFromHex(req.StudentID)
	if err != nil {
		return errors.New("mã sinh viên không hợp lệ")
	}
	subjectID, err := primitive.ObjectIDFromHex(req.SubjectID)
	if err != nil {
		return errors.New("mã môn học không hợp lệ")
	}

	// Kiểm tra sinh viên có tồn tại không
	user, err := s.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return errors.New("lỗi hệ thống")
	}
	if user == nil {
		return errors.New("sinh viên không tồn tại")
	}

	// Kiểm tra môn học có tồn tại không
	subject, err := s.subjectRepo.GetByID(ctx, subjectID)
	if err != nil {
		return errors.New("lỗi hệ thống")
	}
	if subject == nil {
		return errors.New("môn học không tồn tại")
	}

	// Kiểm tra điểm đã tồn tại chưa
	exists, err := s.repo.IsScoreExists(ctx, req.StudentID, req.SubjectID, req.Semester)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("dữ liệu đã tồn tại")
	}

	processScore := calculateProcessScore(req.Attendance, req.Midterm)
	totalScore := calculateTotalScore(processScore, req.Final)

	score := &models.Score{
		StudentID:    studentID,
		SubjectID:    subjectID,
		Semester:     req.Semester,
		Attendance:   req.Attendance,
		Midterm:      req.Midterm,
		Final:        req.Final,
		ProcessScore: processScore,
		Total:        totalScore,
		GPAChar:      convertGPAChar(totalScore),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.CreateScore(ctx, score)
}

func (s *scoreService) CreateScoreByCode(ctx context.Context, req *request.CreateScoreByExcelRequest) error {
	// Tìm sinh viên theo code
	user, err := s.userRepo.GetByCode(ctx, req.StudentCode)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("sinh viên không tồn tại: " + req.StudentCode)
	}

	// Tìm môn học theo code
	subject, err := s.subjectRepo.GetByCode(ctx, req.SubjectCode)
	if err != nil {
		return err
	}
	if subject == nil {
		return errors.New("môn học không tồn tại: " + req.SubjectCode)
	}

	// Kiểm tra điểm đã tồn tại chưa
	exists, err := s.repo.IsScoreExists(ctx, user.ID.Hex(), subject.ID.Hex(), req.Semester)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("điểm cho sinh viên, môn học, học kỳ này đã tồn tại")
	}

	processScore := calculateProcessScore(req.Attendance, req.Midterm)
	totalScore := calculateTotalScore(processScore, req.Final)

	score := &models.Score{
		StudentID:    user.ID,
		SubjectID:    subject.ID,
		Semester:     req.Semester,
		Attendance:   req.Attendance,
		Midterm:      req.Midterm,
		Final:        req.Final,
		ProcessScore: processScore,
		Total:        totalScore,
		GPAChar:      convertGPAChar(totalScore),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.CreateScore(ctx, score)
}
