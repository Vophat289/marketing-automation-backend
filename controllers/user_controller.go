package controllers

import (
	"marketing-automation-backend/config"
	"marketing-automation-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Get danh sách users
func ListUsers(c *gin.Context) {
	var users []models.User

	// 1. Truy vấn toàn bộ danh sách users từ database
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách người dùng"})
		return
	}

	// 2. Lọc bớt thông tin nhạy cảm như Password trước khi trả về
	var userResponses []gin.H
	for _, user := range users {
		userResponses = append(userResponses, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		})
	}

	// 3. Trả về danh sách người dùng
	c.JSON(http.StatusOK, gin.H{
		"message": "Lấy danh sách thành công",
		"data":    userResponses,
	})
}

// Register xử lý logic đăng ký người dùng mới
func Register(c *gin.Context) {
	// Định nghĩa struct để nhận dữ liệu từ request body
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// 1. Lấy dữ liệu từ request và kiểm tra tính hợp lệ
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Kiểm tra xem email đã tồn tại trong database chưa
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email đã được sử dụng"})
		return
	}

	// 3. Mã hóa mật khẩu (Hash password) để bảo mật
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu"})
		return
	}

	// 4. Tạo người dùng mới và lưu vào database
	newUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo người dùng"})
		return
	}

	// 5. Trả về phản hồi thành công (không trả về mật khẩu)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký thành công",
		"user": gin.H{
			"id":    newUser.ID,
			"name":  newUser.Name,
			"email": newUser.Email,
		},
	})
}
