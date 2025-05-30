package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tuyenngduc/certificate-management-system/internal/dto/request"
	"github.com/tuyenngduc/certificate-management-system/internal/handler"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/pkg/database"
	"github.com/tuyenngduc/certificate-management-system/routes"
	"github.com/tuyenngduc/certificate-management-system/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Không tìm thấy file .env, đang dùng biến môi trường hệ thống")
	}
	if err := database.ConnectMongo(); err != nil {
		log.Fatalf("Lỗi khi kết nối MongoDB: %v", err)
	}
	db := database.DB
	seedAdminAccount(db)

	emailSender := utils.NewSMTPSender(
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
	)

	// Khởi tạo repository
	userRepo := repository.NewUserRepository(db)
	trainingDepartmentRepo := repository.NewTrainingDepartmentRepository(db)
	authRepo := repository.NewAuthRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	scoreRepo := repository.NewScoreRepository(db)
	certRepo := repository.NewCertificateRepository(db)
	_ = userRepo.EnsureIndexes(context.Background())

	// Khởi tạo service
	subjectSvc := service.NewSubjectService(subjectRepo, trainingDepartmentRepo)
	userSvc := service.NewUserService(userRepo, trainingDepartmentRepo)
	trainingDepartmentSvc := service.NewTrainingDepartmentService(trainingDepartmentRepo)
	authSvc := service.NewAuthService(authRepo, userRepo, emailSender)
	accountSvc := service.NewAccountService(accountRepo)
	scoreSvc := service.NewScoreService(scoreRepo, userRepo, subjectRepo)
	certService := service.NewCertificateService(certRepo)

	// Khởi tạo handler
	subjectHandler := handler.NewSubjectHandler(subjectSvc, trainingDepartmentSvc)
	userHandler := handler.NewUserHandler(userSvc)
	trainingDepartmentHandler := handler.NewTrainingDepartmentHandler(trainingDepartmentSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	accountHandler := handler.NewAccountHandler(accountSvc)
	scoreHandler := handler.NewScoreHandler(scoreSvc)
	certHandler := handler.NewCertificateHandler(certService)

	// Khởi tạo router
	r := routes.SetupRouter(
		userHandler,
		trainingDepartmentHandler,
		authHandler,
		accountHandler,
		subjectHandler,
		scoreHandler,
		certHandler,
	)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Lỗi khi khởi động server: %v", err)
	}
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		request.RegisterClassValidators(v)
	}
}
