package service

import (
	"context"
	"errors"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo                   *repository.UserRepository
	trainingDepartmentRepo *repository.TrainingDepartmentRepository
}

func NewUserService(repo *repository.UserRepository, trainingDepartmentRepo *repository.TrainingDepartmentRepository) *UserService {
	return &UserService{repo: repo, trainingDepartmentRepo: trainingDepartmentRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {

	if exist, _ := s.repo.FindByEmail(ctx, user.Email); exist != nil {
		return errors.New("email đã tồn tại")
	}
	if exist, _ := s.repo.FindByNationalID(ctx, user.NationalID); exist != nil {
		return errors.New("CCCD/CMND đã tồn tại")
	}

	if user.PhoneNumber != "" {
		if exist, _ := s.repo.FindByPhoneNumber(ctx, user.PhoneNumber); exist != nil {
			return errors.New("số điện thoại đã tồn tại")
		}
	}
	if user.StudentID != "" {
		if exist, _ := s.repo.FindByStudentID(ctx, user.StudentID); exist != nil {
			return errors.New("mã sinh viên đã được đăng ký tài khoản")
		}
	}
	return s.repo.Insert(ctx, user)
}
func (s *UserService) UpdateUser(ctx context.Context, id string, req *request.CreateUserRequest) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id không hợp lệ")
	}

	update := bson.M{}
	if req.FullName != "" {
		update["fullName"] = req.FullName
	}
	if req.Ethnicity != "" {
		update["ethnicity"] = req.Ethnicity
	}
	if req.Gender != "" {
		update["gender"] = req.Gender
	}
	if req.FacultyCode != "" {
		faculty, err := s.FindFacultyByCode(ctx, req.FacultyCode)
		if err != nil || faculty == nil {
			return errors.New("Không tìm thấy khoa với mã " + req.FacultyCode)
		}
		update["facultyId"] = faculty.ID
	}
	if req.ClassCode != "" {
		class, err := s.FindClassByCode(ctx, req.ClassCode)
		if err != nil || class == nil {
			return errors.New("Không tìm thấy lớp với mã " + req.ClassCode)
		}
		update["classId"] = class.ID
	}
	if req.Course != "" {
		update["course"] = req.Course
	}
	if req.NationalID != "" {
		// Kiểm tra unique
		if exist, _ := s.repo.FindByNationalID(ctx, req.NationalID); exist != nil && exist.ID != objID {
			return errors.New("cccd/cmnd đã tồn tại")
		}
		update["nationalId"] = req.NationalID
	}
	if req.StudentID != "" {
		// Kiểm tra unique
		if exist, _ := s.repo.FindByStudentID(ctx, req.StudentID); exist != nil && exist.ID != objID {
			return errors.New("mã sinh viên đã tồn tại")
		}
		update["studentId"] = req.StudentID
	}
	if req.Email != "" {
		// Kiểm tra unique
		if exist, _ := s.repo.FindByEmail(ctx, req.Email); exist != nil && exist.ID != objID {
			return errors.New("email đã tồn tại")
		}
		update["email"] = req.Email
	}
	if req.Address != "" {
		update["address"] = req.Address
	}
	if req.PlaceOfBirth != "" {
		update["placeOfBirth"] = req.PlaceOfBirth
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("02/01/2006", req.DateOfBirth)
		if err != nil {
			return errors.New("sai định dạng ngày sinh")
		}
		update["dateOfBirth"] = dob
	}
	if req.PhoneNumber != "" {
		// Kiểm tra unique
		if exist, _ := s.repo.FindByPhoneNumber(ctx, req.PhoneNumber); exist != nil && exist.ID != objID {
			return errors.New("số điện thoại đã tồn tại")
		}
		update["phoneNumber"] = req.PhoneNumber
	}

	if len(update) == 0 {
		return errors.New("không có dữ liệu cập nhật")
	}

	return s.repo.Update(ctx, objID, update)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) SearchUsers(ctx context.Context, fullName, email, nationalID, phone, studentID string) ([]*models.User, error) {
	return s.repo.Search(ctx, fullName, email, nationalID, phone, studentID)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id không hợp lệ")
	}
	deleted, err := s.repo.Delete(ctx, objID)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("user không tồn tại")
	}
	return nil
}

func (s *UserService) FindFacultyByCode(ctx context.Context, code string) (*models.Faculty, error) {
	return s.trainingDepartmentRepo.FindFacultyByCode(ctx, code)
}

func (s *UserService) FindClassByCode(ctx context.Context, code string) (*models.Class, error) {
	return s.trainingDepartmentRepo.FindClassByCode(ctx, code)
}

func (s *UserService) GetAllFaculties(ctx context.Context) ([]models.Faculty, error) {
	return s.trainingDepartmentRepo.GetAllFaculties(ctx)
}

func (s *UserService) GetAllClasses(ctx context.Context) ([]models.Class, error) {
	return s.trainingDepartmentRepo.GetAllClasses(ctx)
}
