package middleware

import (
	"fmt"
	"backend/config"
	"backend/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy token từ header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Không có token xác thực"})
			return
		}

		// Token thường có dạng "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Kiểm tra xem token có nằm trong Blacklist không (đã đăng xuất)
		var blacklistedToken models.BlacklistedToken
		if err := config.DB.Where("token = ?", tokenString).First(&blacklistedToken).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token đã bị đăng xuất"})
			return
		}

		// 3. Giải mã và xác thực token
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		if len(jwtSecret) == 0 {
			jwtSecret = []byte("my_super_secret_key")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ hoặc đã hết hạn"})
			return
		}

		// 4. Lấy thông tin từ payload
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Không thể đọc thông tin token"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token không chứa user_id"})
			return
		}
		userID := uint(userIDFloat)

		tokenVersionFloat, ok := claims["token_version"].(float64)
		if !ok {
			// Nếu token cũ không có token_version, mặc định là 1
			tokenVersionFloat = 1
		}
		tokenVersion := int(tokenVersionFloat)

		// 5. Kiểm tra TokenVersion trong database
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Người dùng không tồn tại"})
			return
		}

		// Nếu TokenVersion trong token khác với trong DB (do đã đổi mật khẩu), từ chối
		if tokenVersion != user.TokenVersion {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Phiên đăng nhập đã hết hạn do đổi mật khẩu. Vui lòng đăng nhập lại."})
			return
		}

		// 6. Lưu thông tin user vào context để các controller dùng
		c.Set("user", user)
		c.Set("token", tokenString) // Lưu lại token string để dùng cho chức năng Logout

		// Tiếp tục xử lý request
		c.Next()
	}
}
