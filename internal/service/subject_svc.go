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

type SubjectService interface {
	CreateSubject(ctx context.Context, req *request.CreateSubjectRequest) error
	UpdateSubject(ctx context.Context, id string, req *request.UpdateSubjectRequest) error
	DeleteSubject(ctx context.Context, id string) error
	GetSubjectByID(ctx context.Context, id string) (*models.Subject, error)
	ListSubjects(ctx context.Context) ([]*models.Subject, error)
	CreateSubjectByFacultyCode(ctx context.Context, req *request.CreateSubjectByExcelRequest) error
}

type subjectService struct {
	subjectRepo  repository.SubjectRepository
	trainingRepo *repository.TrainingDepartmentRepository
}

func NewSubjectService(subjectRepo repository.SubjectRepository, trainingRepo *repository.TrainingDepartmentRepository) SubjectService {
	return &subjectService{
		subjectRepo:  subjectRepo,
		trainingRepo: trainingRepo,
	}
}

func (s *subjectService) CreateSubject(ctx context.Context, req *request.CreateSubjectRequest) error {
	// Kiểm tra mã môn học trùng
	existing, err := s.subjectRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("mã môn học đã tồn tại")
	}

	// Kiểm tra FacultyID hợp lệ
	facultyID, err := primitive.ObjectIDFromHex(req.FacultyID)
	if err != nil {
		return errors.New("id khoa không hợp lệ")
	}

	// 🟡 Kiểm tra khoa có tồn tại không
	faculty, err := s.trainingRepo.GetFacultyByID(ctx, facultyID)
	if err != nil || faculty == nil {
		return errors.New("khoa không tồn tại")
	}

	// Tạo subject
	subject := &models.Subject{
		Code:        req.Code,
		Name:        req.Name,
		Credit:      req.Credit,
		FacultyID:   facultyID,
		Description: req.Description,
	}

	return s.subjectRepo.Create(ctx, subject)
}
func (s *subjectService) UpdateSubject(ctx context.Context, id string, req *request.UpdateSubjectRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id môn học không hợp lệ")
	}

	// Kiểm tra subject có tồn tại không
	existingSubject, err := s.subjectRepo.GetByID(ctx, objectID)
	if err != nil {
		return errors.New("lỗi hệ thống")
	}
	if existingSubject == nil {
		return errors.New("môn học không tồn tại")
	}

	update := bson.M{}

	// Kiểm tra code nếu truyền lên
	if req.Code != nil {
		other, err := s.subjectRepo.GetByCode(ctx, *req.Code)
		if err != nil {
			return errors.New("lỗi hệ thống")
		}
		// Nếu đã có subject khác với code này thì báo lỗi
		if other != nil && other.ID.Hex() != objectID.Hex() {
			return errors.New("mã môn học đã tồn tại")
		}
		update["code"] = *req.Code
	}

	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.Credit != nil {
		update["credit"] = *req.Credit
	}
	if req.FacultyID != nil {
		facultyID, err := primitive.ObjectIDFromHex(*req.FacultyID)
		if err != nil {
			return errors.New("id khoa không hợp lệ")
		}
		faculty, err := s.trainingRepo.GetFacultyByID(ctx, facultyID)
		if err != nil || faculty == nil {
			return errors.New("khoa không tồn tại")
		}
		update["faculty_id"] = facultyID
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if len(update) == 0 {
		return errors.New("không có dữ liệu cập nhật")
	}

	return s.subjectRepo.Update(ctx, objectID, update)
}

func (s *subjectService) DeleteSubject(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id môn học không hợp lệ")
	}
	return s.subjectRepo.Delete(ctx, objectID)
}

func (s *subjectService) GetSubjectByID(ctx context.Context, id string) (*models.Subject, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("id môn học không hợp lệ")
	}
	return s.subjectRepo.GetByID(ctx, objectID)
}

func (s *subjectService) ListSubjects(ctx context.Context) ([]*models.Subject, error) {
	return s.subjectRepo.List(ctx)
}

func (s *subjectService) CreateSubjectByFacultyCode(ctx context.Context, req *request.CreateSubjectByExcelRequest) error {
	// Kiểm tra mã môn học trùng
	existing, err := s.subjectRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("mã môn học đã tồn tại")
	}

	// Tìm faculty theo code
	faculty, err := s.trainingRepo.GetFacultyByCode(ctx, req.FacultyCode)
	if err != nil || faculty == nil {
		return errors.New("khoa không tồn tại")
	}

	subject := &models.Subject{
		Code:        req.Code,
		Name:        req.Name,
		Credit:      req.Credit,
		FacultyID:   faculty.ID,
		Description: req.Description,
	}
	return s.subjectRepo.Create(ctx, subject)
}
