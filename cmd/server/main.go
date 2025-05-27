package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	_ = userRepo.EnsureIndexes(context.Background())

	// Khởi tạo service
	userSvc := service.NewUserService(userRepo, trainingDepartmentRepo)
	trainingDepartmentSvc := service.NewTrainingDepartmentService(trainingDepartmentRepo)
	authSvc := service.NewAuthService(authRepo, userRepo, emailSender)
	accountSvc := service.NewAccountService(accountRepo)

	// Khởi tạo handler
	userHandler := handler.NewUserHandler(userSvc)
	trainingDepartmentHandler := handler.NewTrainingDepartmentHandler(trainingDepartmentSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	accountHandler := handler.NewAccountHandler(accountSvc)

	// Khởi tạo router
	r := routes.SetupRouter(userHandler, trainingDepartmentHandler, authHandler, accountHandler)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Lỗi khi khởi động server: %v", err)
	}
}
