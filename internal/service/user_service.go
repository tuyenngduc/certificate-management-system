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
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	SearchUsers(ctx context.Context, params models.SearchUserParams) ([]*models.User, int64, error)
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	UpdateUser(ctx context.Context, id primitive.ObjectID, req models.UpdateUserRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}
func (s *userService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}
func (s *userService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}
func (s *userService) SearchUsers(ctx context.Context, params models.SearchUserParams) ([]*models.User, int64, error) {
	return s.userRepo.SearchUsers(ctx, params)
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Kiểm tra StudentID tồn tại
	exists, err := s.userRepo.ExistsByStudentID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrStudentIDExists
	}

	// Kiểm tra Email tồn tại
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrEmailExists
	}

	user := &models.User{
		ID:        primitive.NewObjectID(),
		StudentID: req.StudentID,
		FullName:  req.FullName,
		Email:     req.Email,
		Faculty:   req.Faculty,
		Class:     req.Class,
		Course:    req.Course,
		Status:    "Đang học",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := &models.UserResponse{
		ID:        user.ID,
		StudentID: user.StudentID,
		FullName:  user.FullName,
		Email:     user.Email,
		Faculty:   user.Faculty,
		Class:     user.Class,
		Course:    user.Course,
		Status:    user.Status,
	}
	return resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, id primitive.ObjectID, req models.UpdateUserRequest) error {
	update := bson.M{}
	if req.StudentID != "" {
		exist, err := s.userRepo.FindByStudentID(ctx, req.StudentID)
		if err == nil && exist != nil && exist.ID != id {
			return common.ErrStudentIDExists
		}
		update["studentId"] = req.StudentID
	}
	if req.Email != "" {
		exist, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && exist != nil && exist.ID != id {
			return common.ErrEmailExists
		}
		update["email"] = req.Email
	}

	if req.StudentID != "" {
		update["studentId"] = req.StudentID
	}
	if req.FullName != "" {
		update["fullName"] = req.FullName
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.Faculty != "" {
		update["facultyId"] = req.Faculty
	}
	if req.Class != "" {
		update["classId"] = req.Class
	}
	if req.Course != "" {
		update["course"] = req.Course
	}
	if req.Status != "" {
		update["status"] = req.Status
	}
	update["updatedAt"] = time.Now()
	if len(update) == 1 {
		return errors.New("không có trường nào để cập nhật")
	}
	return s.userRepo.UpdateUser(ctx, id, update)
}

func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.userRepo.DeleteUser(ctx, id)
}
