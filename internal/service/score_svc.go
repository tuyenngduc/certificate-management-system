package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/response"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreService interface {
	CreateScore(ctx context.Context, req *request.CreateScoreRequest) error
	CreateScoreByCode(ctx context.Context, req *request.CreateScoreByExcelRequest) error
	GetScoresBySubjectID(ctx context.Context, subjectID string) ([]*response.ScoreWithSubjectAndStudentResponse, error)
	UpdateScore(ctx context.Context, id string, req *request.UpdateScoreRequest) error
	ImportScoresBySubjectExcel(ctx context.Context, subjectCode string, reqs []request.ImportScoresBySubjectExcelRequest) ([]string, error)
	CalculateCGPA(ctx context.Context, studentID string) (*response.CGPAResponse, error)
	GetScoresByStudentID(ctx context.Context, studentID string) ([]*response.ScoreWithSubjectAndStudentResponse, error)
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

func roundToOneDecimalPlace(val float64) float64 {
	return math.Round(val*10) / 10
}

func convertGPAChar(score float64) string {
	switch {
	case score >= 9.0:
		return "A+"
	case score >= 8.5:
		return "A"
	case score >= 7.8:
		return "B+"
	case score >= 7.0:
		return "B"
	case score >= 6.3:
		return "C+"
	case score >= 5.5:
		return "C"
	case score >= 4.8:
		return "D+"
	case score >= 4.0:
		return "D"
	default:
		return "F"
	}
}
func gpaCharTo4Scale(gpaChar string) float64 {
	switch gpaChar {
	case "A+":
		return 4.0
	case "A":
		return 3.8
	case "B+":
		return 3.5
	case "B":
		return 3.0
	case "C+":
		return 2.4
	case "D+":
		return 1.5
	case "D":
		return 1.0
	default:
		return 0.0
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

	processScore := roundToOneDecimalPlace(calculateProcessScore(req.Attendance, req.Midterm))
	totalScore := roundToOneDecimalPlace(calculateTotalScore(processScore, req.Final))

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
		Passed:       req.Final >= 2 && totalScore >= 4.0,
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
		Passed:       req.Final >= 2 && totalScore >= 4.0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.CreateScore(ctx, score)
}

func (s *scoreService) GetScoresByStudentID(ctx context.Context, studentID string) ([]*response.ScoreWithSubjectAndStudentResponse, error) {
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, errors.New("mã sinh viên không hợp lệ")
	}

	scores, err := s.repo.GetScoresByStudentID(ctx, objID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, objID)
	if err != nil || user == nil {
		return nil, errors.New("không tìm thấy sinh viên")
	}

	var results []*response.ScoreWithSubjectAndStudentResponse
	for _, score := range scores {
		subject, err := s.subjectRepo.GetByID(ctx, score.SubjectID)
		if err != nil || subject == nil {
			continue
		}
		fmt.Println("DEBUG:", score.Final, score.Total, score.Final >= 2 && score.Total >= 4.0)

		results = append(results, &response.ScoreWithSubjectAndStudentResponse{
			ID:           score.ID,
			StudentName:  user.FullName,
			SubjectName:  subject.Name,
			Credit:       subject.Credit,
			Semester:     score.Semester,
			Attendance:   score.Attendance,
			Midterm:      score.Midterm,
			Final:        score.Final,
			ProcessScore: roundToOneDecimalPlace(score.ProcessScore),
			Total:        roundToOneDecimalPlace(score.Total),
			GPAChar:      score.GPAChar,
			Passed:       score.Final >= 2 && score.Total >= 4.0,
		})
	}

	return results, nil
}

func (s *scoreService) GetScoresBySubjectID(ctx context.Context, subjectID string) ([]*response.ScoreWithSubjectAndStudentResponse, error) {
	objID, err := primitive.ObjectIDFromHex(subjectID)
	if err != nil {
		return nil, errors.New("mã môn học không hợp lệ")
	}

	scores, err := s.repo.GetScoresBySubjectID(ctx, objID)
	if err != nil {
		return nil, err
	}

	subject, err := s.subjectRepo.GetByID(ctx, objID)
	if err != nil || subject == nil {
		return nil, errors.New("không tìm thấy môn học")
	}

	var result []*response.ScoreWithSubjectAndStudentResponse
	for _, score := range scores {

		user, err := s.userRepo.GetByID(ctx, score.StudentID)
		if err != nil || user == nil {
			continue
		}

		result = append(result, &response.ScoreWithSubjectAndStudentResponse{
			ID:           score.ID,
			StudentName:  user.FullName,
			SubjectName:  subject.Name,
			Credit:       subject.Credit,
			Semester:     score.Semester,
			Attendance:   score.Attendance,
			Midterm:      score.Midterm,
			Final:        score.Final,
			ProcessScore: roundToOneDecimalPlace(score.ProcessScore),
			Total:        roundToOneDecimalPlace(score.Total),
			GPAChar:      score.GPAChar,
			Passed:       score.Final >= 2 && score.Total >= 4.0,
		})
	}

	return result, nil
}

func (s *scoreService) UpdateScore(ctx context.Context, id string, req *request.UpdateScoreRequest) error {
	scoreID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id điểm không hợp lệ")
	}

	// Lấy bản ghi điểm hiện tại
	score, err := s.repo.GetScoreByID(ctx, scoreID)
	if err != nil {
		return errors.New("lỗi hệ thống")
	}
	if score == nil {
		return errors.New("không tìm thấy điểm")
	}

	update := bson.M{}

	if req.Attendance != nil {
		update["attendance"] = *req.Attendance
	}
	if req.Midterm != nil {
		update["midterm"] = *req.Midterm
	}
	if req.Final != nil {
		update["final"] = *req.Final
	}

	// Tính lại điểm quá trình, tổng kết, GPAChar nếu có thay đổi
	attendance := score.Attendance
	midterm := score.Midterm
	final := score.Final
	if req.Attendance != nil {
		attendance = *req.Attendance
	}
	if req.Midterm != nil {
		midterm = *req.Midterm
	}
	if req.Final != nil {
		final = *req.Final
	}
	processScore := roundToOneDecimalPlace(calculateProcessScore(attendance, midterm))
	totalScore := roundToOneDecimalPlace(calculateTotalScore(processScore, final))
	update["process_score"] = processScore
	update["total"] = totalScore
	update["gpa_char"] = convertGPAChar(totalScore)

	return s.repo.UpdateScore(ctx, scoreID, update)
}

func (s *scoreService) ImportScoresBySubjectExcel(ctx context.Context, subjectID string, reqs []request.ImportScoresBySubjectExcelRequest) ([]string, error) {
	var results []string

	subjectObjID, err := primitive.ObjectIDFromHex(subjectID)
	if err != nil {
		return nil, errors.New("mã môn học không hợp lệ")
	}

	subject, err := s.subjectRepo.GetByID(ctx, subjectObjID)
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi lấy môn học")
	}
	if subject == nil {
		return nil, errors.New("không tìm thấy môn học")
	}

	for i, req := range reqs {
		// Tìm sinh viên theo code
		user, err := s.userRepo.GetByCode(ctx, req.StudentID)
		if err != nil {
			results = append(results, "Dòng "+itoa(i+2)+": lỗi hệ thống khi lấy sinh viên")
			continue
		}
		if user == nil {
			results = append(results, "Dòng "+itoa(i+2)+": sinh viên không tồn tại")
			continue
		}

		// Kiểm tra điểm đã tồn tại chưa
		exists, err := s.repo.IsScoreExists(ctx, user.ID.Hex(), subject.ID.Hex(), req.Semester)
		if err != nil {
			results = append(results, "Dòng "+itoa(i+2)+": lỗi kiểm tra trùng điểm")
			continue
		}
		if exists {
			results = append(results, "Dòng "+itoa(i+2)+": điểm đã tồn tại")
			continue
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
			Passed:       req.Final >= 2 && totalScore >= 4.0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err = s.repo.CreateScore(ctx, score)
		if err != nil {
			results = append(results, "Dòng "+itoa(i+2)+": lỗi lưu điểm")
		} else {
			results = append(results, "Dòng "+itoa(i+2)+": thành công")
		}
	}

	return results, nil
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func (s *scoreService) CalculateCGPA(ctx context.Context, studentID string) (*response.CGPAResponse, error) {
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil, errors.New("mã sinh viên không hợp lệ")
	}

	// Lấy tất cả điểm của sinh viên
	scores, err := s.repo.GetScoresByStudentID(ctx, objID)
	if err != nil {
		return nil, err
	}
	if len(scores) == 0 {
		return &response.CGPAResponse{CGPA: 0, TotalSubjects: 0, TotalCredits: 0, TotalFailedSubjects: 0}, nil
	}

	var (
		totalWeighted       float64
		totalCredits        int
		totalFailedSubjects int
		totalSubjects       int
	)

	for _, score := range scores {
		subject, err := s.subjectRepo.GetByID(ctx, score.SubjectID)
		if err != nil || subject == nil {
			continue
		}
		totalSubjects++

		credit := subject.Credit
		if score.Passed {
			totalWeighted += gpaCharTo4Scale(score.GPAChar) * float64(credit)
			totalCredits += credit
		} else {
			totalFailedSubjects++
		}
	}

	cgpa := 0.0
	if totalCredits > 0 {
		cgpa = roundToOneDecimalPlace(totalWeighted / float64(totalCredits))
	}

	return &response.CGPAResponse{
		CGPA:                cgpa,
		TotalSubjects:       totalSubjects,
		TotalCredits:        totalCredits,
		TotalFailedSubjects: totalFailedSubjects,
	}, nil
}
