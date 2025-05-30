package service

import (
	"context"
	"errors"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingDepartmentService struct {
	repo *repository.TrainingDepartmentRepository
}

func NewTrainingDepartmentService(repo *repository.TrainingDepartmentRepository) *TrainingDepartmentService {
	return &TrainingDepartmentService{repo: repo}
}

// Faculty
func (s *TrainingDepartmentService) CreateFaculty(ctx context.Context, req *request.CreateFacultyRequest) error {
	exist, _ := s.repo.FindFacultyByCode(ctx, req.Code)
	if exist != nil {
		return errors.New("mã khoa đã tồn tại")
	}
	return s.repo.CreateFaculty(ctx, &models.Faculty{
		Name:           req.Name,
		Code:           req.Code,
		TrainingPeriod: req.TrainingPeriod,
	})
}
func (s *TrainingDepartmentService) GetAllFaculties(ctx context.Context) ([]models.Faculty, error) {
	return s.repo.GetAllFaculties(ctx)
}
func (s *TrainingDepartmentService) GetFacultyByID(ctx context.Context, id primitive.ObjectID) (*models.Faculty, error) {
	return s.repo.GetFacultyByID(ctx, id)
}

func (s *TrainingDepartmentService) UpdateFaculty(ctx context.Context, id string, req *request.UpdateFacultyRequest) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID không hợp lệ")
	}

	// Kiểm tra faculty có tồn tại không
	faculty, _ := s.repo.GetFacultyByID(ctx, objID)
	if faculty == nil {
		return errors.New("không tìm thấy khoa")
	}

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Code != "" {
		// Kiểm tra mã khoa đã tồn tại cho faculty khác
		exist, _ := s.repo.FindFacultyByCode(ctx, req.Code)
		if exist != nil && exist.ID != objID {
			return errors.New("mã khoa đã tồn tại")
		}
		update["code"] = req.Code
	}
	if len(update) == 0 {
		return errors.New("không có dữ liệu cập nhật")
	}
	if req.TrainingPeriod != "" {
		update["trainingPeriod"] = req.TrainingPeriod
	}

	return s.repo.UpdateFaculty(ctx, objID, update)
}
func (s *TrainingDepartmentService) DeleteFaculty(ctx context.Context, id primitive.ObjectID) error {
	deleted, err := s.repo.DeleteFaculty(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("không tìm thấy khoa")
	}
	return nil
}
func (s *TrainingDepartmentService) FindFacultyByCode(ctx context.Context, code string) (*models.Faculty, error) {
	return s.repo.FindFacultyByCode(ctx, code)
}

// Class
func (s *TrainingDepartmentService) GetClassesByFacultyID(ctx context.Context, facultyID string) ([]models.Class, error) {
	objID, err := primitive.ObjectIDFromHex(facultyID)
	if err != nil {
		return nil, errors.New("ID khoa không hợp lệ")
	}
	return s.repo.GetClassesByFacultyID(ctx, objID)
}
func (s *TrainingDepartmentService) CreateClass(ctx context.Context, req *request.CreateClassRequest) error {
	exist, _ := s.repo.FindClassByCode(ctx, req.Code)
	if exist != nil {
		return errors.New("mã lớp đã tồn tại")
	}
	facultyID, err := primitive.ObjectIDFromHex(req.FacultyID)
	if err != nil {
		return errors.New("ID khoa không hợp lệ")
	}
	faculty, _ := s.repo.GetFacultyByID(ctx, facultyID)
	if faculty == nil {
		return errors.New("ID khoa không tồn tại")
	}
	return s.repo.CreateClass(ctx, &models.Class{
		Code:      req.Code,
		Course:    req.Course,
		FacultyID: facultyID,
	})
}
func (s *TrainingDepartmentService) GetAllClasses(ctx context.Context) ([]models.Class, error) {
	return s.repo.GetAllClasses(ctx)
}
func (s *TrainingDepartmentService) GetClassByID(ctx context.Context, id primitive.ObjectID) (*models.Class, error) {
	return s.repo.GetClassByID(ctx, id)
}
func (s *TrainingDepartmentService) UpdateClass(ctx context.Context, id string, req *request.UpdateClassRequest) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID không hợp lệ")
	}

	// Kiểm tra class có tồn tại không
	class, _ := s.repo.GetClassByID(ctx, objID)
	if class == nil {
		return errors.New("không tìm thấy lớp")
	}

	update := bson.M{}
	if req.Code != "" {
		exist, _ := s.repo.FindClassByCode(ctx, req.Code)
		if exist != nil && exist.ID != objID {
			return errors.New("mã lớp đã tồn tại")
		}
		update["code"] = req.Code
	}
	if req.Course != "" {
		update["course"] = req.Course
	}
	if req.FacultyID != "" {
		facultyID, err := primitive.ObjectIDFromHex(req.FacultyID)
		if err != nil {
			return errors.New("ID khoa không hợp lệ")
		}
		faculty, _ := s.repo.GetFacultyByID(ctx, facultyID)
		if faculty == nil {
			return errors.New("ID khoa không tồn tại")
		}
		update["faculty_id"] = facultyID
	}
	if len(update) == 0 {
		return errors.New("không có dữ liệu cập nhật")
	}

	return s.repo.UpdateClass(ctx, objID, update)
}
func (s *TrainingDepartmentService) DeleteClass(ctx context.Context, id primitive.ObjectID) error {
	deleted, err := s.repo.DeleteClass(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("không tìm thấy lớp")
	}
	return nil
}
func (s *TrainingDepartmentService) FindClassByCode(ctx context.Context, code string) (*models.Class, error) {
	return s.repo.FindClassByCode(ctx, code)
}

// Lecturer
func (s *TrainingDepartmentService) CreateLecturer(ctx context.Context, req *request.CreateLecturerRequest) error {
	exist, _ := s.repo.FindLecturerByCode(ctx, req.Code)
	if exist != nil {
		return errors.New("mã giảng viên đã tồn tại")
	}
	existEmail, _ := s.repo.FindLecturerByEmail(ctx, req.Email)
	if existEmail != nil {
		return errors.New("email giảng viên đã tồn tại")
	}

	facultyID, err := primitive.ObjectIDFromHex(req.FacultyID)
	if err != nil {
		return errors.New("ID khoa không hợp lệ")
	}
	faculty, _ := s.repo.GetFacultyByID(ctx, facultyID)
	if faculty == nil {
		return errors.New("ID khoa không tồn tại")
	}

	return s.repo.CreateLecturer(ctx, &models.Lecturer{
		ID:        primitive.NewObjectID(),
		Code:      req.Code,
		FullName:  req.FullName,
		Email:     req.Email,
		Title:     req.Title,
		FacultyID: facultyID,
	})
}
func (s *TrainingDepartmentService) GetLecturersByFacultyID(ctx context.Context, facultyID string) ([]models.Lecturer, error) {
	objID, err := primitive.ObjectIDFromHex(facultyID)
	if err != nil {
		return nil, errors.New("ID khoa không hợp lệ")
	}
	faculty, _ := s.repo.GetFacultyByID(ctx, objID)
	if faculty == nil {
		return nil, errors.New("khoa không tồn tại")
	}
	return s.repo.GetLecturersByFacultyID(ctx, objID)
}
func (s *TrainingDepartmentService) GetAllLecturers(ctx context.Context) ([]models.Lecturer, error) {
	return s.repo.GetAllLecturers(ctx)
}
func (s *TrainingDepartmentService) SearchLecturers(ctx context.Context, id, code, fullName string) ([]*models.Lecturer, error) {
	return s.repo.SearchLecturers(ctx, id, code, fullName)
}
func (s *TrainingDepartmentService) GetLecturerByID(ctx context.Context, id primitive.ObjectID) (*models.Lecturer, error) {
	return s.repo.GetLecturerByID(ctx, id)
}

func (s *TrainingDepartmentService) UpdateLecturer(ctx context.Context, id string, req *request.UpdateLecturerRequest) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID không hợp lệ")
	}
	lecturer, _ := s.repo.GetLecturerByID(ctx, objID)
	if lecturer == nil {
		return errors.New("không tìm thấy giảng viên")
	}

	update := bson.M{}
	if req.Code != "" {
		// Kiểm tra mã giảng viên đã tồn tại cho lecturer khác
		exist, _ := s.repo.FindLecturerByCode(ctx, req.Code)
		if exist != nil && exist.ID != objID {
			return errors.New("mã giảng viên đã tồn tại")
		}
		update["code"] = req.Code
	}
	if req.FullName != "" {
		update["fullName"] = req.FullName
	}
	if req.Email != "" {
		existEmail, _ := s.repo.FindLecturerByEmail(ctx, req.Email)
		if existEmail != nil && existEmail.ID != objID {
			return errors.New("email giảng viên đã tồn tại")
		}
		update["email"] = req.Email
	}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if len(update) == 0 {
		return errors.New("không có dữ liệu cập nhật")
	}

	return s.repo.UpdateLecturer(ctx, objID, update)
}

func (s *UserService) GetUsersByClassID(ctx context.Context, classID string) ([]*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(classID)
	if err != nil {
		return nil, errors.New("ID lớp không hợp lệ")
	}
	return s.repo.GetUsersByClassID(ctx, objID)
}
func (s *TrainingDepartmentService) SearchClasses(ctx context.Context, id, code, course string, page, pageSize int) ([]*models.Class, int, error) {
	return s.repo.SearchClasses(ctx, id, code, course, page, pageSize)
}
func (s *TrainingDepartmentService) DeleteLecturer(ctx context.Context, id primitive.ObjectID) error {
	deleted, err := s.repo.DeleteLecturer(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("không tìm thấy giảng viên")
	}
	return nil
}

func (s *TrainingDepartmentService) FindLecturerByCode(ctx context.Context, code string) (*models.Lecturer, error) {
	return s.repo.FindLecturerByCode(ctx, code)
}
