package main

import (
	"context"
	"log"

	"github.com/tuyenngduc/certificate-management-system/internal/handler"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/pkg/database"
	"github.com/tuyenngduc/certificate-management-system/routes"
)

func main() {
	if err := database.ConnectMongo(); err != nil {
		log.Fatalf("Lỗi khi kết nối MongoDB: %v", err)
	}

	db := database.DB

	// Khởi tạo repository
	userRepo := repository.NewUserRepository(db)
	trainingDepartmentRepo := repository.NewTrainingDepartmentRepository(db)
	_ = userRepo.EnsureIndexes(context.Background())

	// Khởi tạo service
	userSvc := service.NewUserService(userRepo, trainingDepartmentRepo)
	trainingDepartmentSvc := service.NewTrainingDepartmentService(trainingDepartmentRepo)

	// Khởi tạo handler
	userHandler := handler.NewUserHandler(userSvc)
	trainingDepartmentHandler := handler.NewTrainingDepartmentHandler(trainingDepartmentSvc)

	// Khởi tạo router
	r := routes.SetupRouter(userHandler, trainingDepartmentHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Lỗi khi khởi động server: %v", err)
	}
}
