package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo() {
	// Load biến môi trường từ .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Không tìm thấy file .env, đang dùng biến môi trường hệ thống")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")

	if mongoURI == "" || dbName == "" {
		log.Fatal("❌ Thiếu biến MONGODB_URI hoặc DB_NAME")
	}

	// Kết nối MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Kết nối MongoDB thất bại: %v", err)
	}

	// Kiểm tra kết nối
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Không thể ping MongoDB: %v", err)
	}

	MongoClient = client
	MongoDB = client.Database(dbName)

	fmt.Println("✅ Đã kết nối MongoDB thành công")
}
