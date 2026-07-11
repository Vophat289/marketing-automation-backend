package main

import (
	"log"
	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load biến môi trường từ file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Không tìm thấy file .env, sẽ sử dụng biến môi trường hệ thống")
	}

	// 1. Kết nối Database
	config.ConnectDB()

	// 2. Tự động tạo bảng User và BlacklistedToken nếu chưa có
	config.DB.AutoMigrate(
		&models.User{},
		&models.BlacklistedToken{},
	)

	// 3. Khởi tạo Gin router
	router := gin.Default()

	// 4. Thiết lập các routes
	routes.SetupUserRoutes(router)

	// 5. Khởi chạy server ở port 8080
	log.Println("Server đang chạy ở port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Không thể khởi chạy server: ", err)
	}
}
