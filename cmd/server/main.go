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
	userRepo := repository.NewUserRepository(db)
	_ = userRepo.EnsureIndexes(context.Background())
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	r := routes.SetupRouter((userHandler))
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Lỗi khi khởi động server: %v", err)
	}
}
