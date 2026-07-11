package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
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

// LoginUser xử lý logic đăng nhập và trả về JWT Token
func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// 1. Lấy dữ liệu từ request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// 2. Tìm user trong database bằng email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
		return
	}

	// 3. Kiểm tra mật khẩu
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
		return
	}

	// 4. Tạo JWT Token
	token, err := utils.GenerateToken(user.ID, user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token đăng nhập"})
		return
	}

	// 5. Trả về token cho client
	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Logout xử lý đăng xuất bằng cách đưa token vào blacklist
func Logout(c *gin.Context) {
	// Lấy token từ context (đã được set trong middleware)
	tokenString, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy token"})
		return
	}

	// Lưu token vào bảng BlacklistedToken
	blacklistedToken := models.BlacklistedToken{
		Token: tokenString.(string),
	}

	if err := config.DB.Create(&blacklistedToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đăng xuất"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đăng xuất thành công"})
}

// ChangePassword xử lý đổi mật khẩu và vô hiệu hóa các phiên đăng nhập khác
func ChangePassword(c *gin.Context) {
	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Lấy user từ context (đã được set trong middleware)
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin người dùng"})
		return
	}
	user := userContext.(models.User)

	// Kiểm tra mật khẩu cũ
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mật khẩu cũ không đúng"})
		return
	}

	// Mã hóa mật khẩu mới
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mã hóa mật khẩu mới"})
		return
	}

	// Cập nhật mật khẩu và tăng TokenVersion
	user.Password = string(hashedPassword)
	user.TokenVersion += 1

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật mật khẩu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đổi mật khẩu thành công. Các thiết bị khác sẽ bị đăng xuất."})
}
