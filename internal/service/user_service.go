package service

import (
	"context"
	"errors"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.UserResponse, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.UserResponse, error)
	SearchUsers(ctx context.Context, params models.SearchUserParams) ([]models.UserResponse, int64, error)
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	UpdateUser(ctx context.Context, id primitive.ObjectID, req models.UpdateUserRequest) error
}

type userService struct {
	userRepo       repository.UserRepository
	universityRepo repository.UniversityRepository
	facultyRepo    repository.FacultyRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	universityRepo repository.UniversityRepository,
	facultyRepo repository.FacultyRepository,
) UserService {
	return &userService{
		userRepo:       userRepo,
		universityRepo: universityRepo,
		facultyRepo:    facultyRepo,
	}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.UserResponse
	for _, u := range users {
		faculty, err := s.facultyRepo.FindByID(ctx, u.FacultyID)
		if err != nil || faculty == nil {
			continue // hoặc xử lý lỗi
		}

		university, err := s.universityRepo.FindByID(ctx, faculty.UniversityID)
		if err != nil || university == nil {
			continue // hoặc xử lý lỗi
		}

		result = append(result, models.UserResponse{
			ID:             u.ID,
			StudentCode:    u.StudentCode,
			FullName:       u.FullName,
			Email:          u.Email,
			Course:         u.Course,
			Status:         u.Status,
			FacultyCode:    faculty.FacultyCode,
			FacultyName:    faculty.FacultyName,
			UniversityCode: university.UniversityCode,
			UniversityName: university.UniversityName,
		})
	}

	return result, nil
}

func (s *userService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}

	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	university, err := s.universityRepo.FindByID(ctx, user.UniversityID)
	if err != nil || university == nil {
		return nil, common.ErrUniversityNotFound
	}

	return &models.UserResponse{
		ID:             user.ID,
		StudentCode:    user.StudentCode,
		FullName:       user.FullName,
		Email:          user.Email,
		Course:         user.Course,
		Status:         user.Status,
		FacultyCode:    faculty.FacultyCode,
		FacultyName:    faculty.FacultyName,
		UniversityCode: university.UniversityCode,
		UniversityName: university.UniversityName,
	}, nil
}

func (s *userService) SearchUsers(ctx context.Context, params models.SearchUserParams) ([]models.UserResponse, int64, error) {
	users, total, err := s.userRepo.SearchUsers(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.UserResponse
	for _, u := range users {
		faculty, _ := s.facultyRepo.FindByID(ctx, u.FacultyID)
		university, _ := s.universityRepo.FindByID(ctx, u.UniversityID)

		resp := models.UserResponse{
			ID:             u.ID,
			StudentCode:    u.StudentCode,
			FullName:       u.FullName,
			Email:          u.Email,
			Course:         u.Course,
			Status:         u.Status,
			FacultyCode:    "",
			FacultyName:    "",
			UniversityCode: "",
			UniversityName: "",
		}

		if faculty != nil {
			resp.FacultyCode = faculty.FacultyCode
			resp.FacultyName = faculty.FacultyName
		}
		if university != nil {
			resp.UniversityCode = university.UniversityCode
			resp.UniversityName = university.UniversityName
		}

		responses = append(responses, resp)
	}

	return responses, total, nil
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Kiểm tra mã sinh viên trùng
	exists, err := s.userRepo.ExistsByStudentCode(ctx, req.StudentCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrStudentIDExists
	}

	// Kiểm tra Email trùng
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrEmailExists
	}

	// Tìm University theo code
	university, err := s.universityRepo.FindByCode(ctx, req.UniversityCode)
	if err != nil || university == nil {
		return nil, common.ErrUniversityNotFound
	}

	// Tìm Faculty theo code và university_id
	faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, req.FacultyCode, university.ID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	user := &models.User{
		ID:           primitive.NewObjectID(),
		StudentCode:  req.StudentCode,
		FullName:     req.FullName,
		Email:        req.Email,
		FacultyID:    faculty.ID,
		UniversityID: university.ID,
		Course:       req.Course,
		Status:       "Đang học",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := &models.UserResponse{
		ID:             user.ID,
		StudentCode:    user.StudentCode,
		FullName:       user.FullName,
		Email:          user.Email,
		FacultyCode:    faculty.FacultyCode,
		FacultyName:    faculty.FacultyName,
		UniversityCode: university.UniversityCode,
		UniversityName: university.UniversityName,
		Course:         user.Course,
		Status:         user.Status,
	}

	return resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, id primitive.ObjectID, req models.UpdateUserRequest) error {
	update := bson.M{}

	// Check trùng mã sinh viên
	if req.StudentCode != "" {
		exist, err := s.userRepo.FindByStudentCode(ctx, req.StudentCode)
		if err == nil && exist != nil && exist.ID != id {
			return common.ErrStudentIDExists
		}
		update["student_code"] = req.StudentCode
	}

	// Check trùng email
	if req.Email != "" {
		exist, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && exist != nil && exist.ID != id {
			return common.ErrEmailExists
		}
		update["email"] = req.Email
	}

	// Cập nhật các trường thông thường
	if req.FullName != "" {
		update["full_name"] = req.FullName
	}
	if req.Course != "" {
		update["course"] = req.Course
	}
	if req.Status != "" {
		update["status"] = req.Status
	}

	// Tìm Faculty theo code nếu có yêu cầu
	if req.FacultyCode != "" && req.UniversityCode != "" {
		university, err := s.universityRepo.FindByCode(ctx, req.UniversityCode)
		if err != nil || university == nil {
			return common.ErrUniversityNotFound
		}

		faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, req.FacultyCode, university.ID)
		if err != nil || faculty == nil {
			return common.ErrFacultyNotFound
		}

		update["faculty_id"] = faculty.ID
		update["university_id"] = university.ID
	} else if req.FacultyCode != "" || req.UniversityCode != "" {
		return errors.New("phải cung cấp cả faculty_code và university_code")
	}

	// Thêm trường updatedAt
	update["updated_at"] = time.Now()

	if len(update) == 1 {
		return errors.New("không có trường nào để cập nhật")
	}

	return s.userRepo.UpdateUser(ctx, id, update)
}

func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.userRepo.DeleteUser(ctx, id)
}
