package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Khóa bí mật dùng để ký token (nên lấy từ biến môi trường)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken tạo ra một JWT token cho user
func GenerateToken(userID uint) (string, error) {
	// Nếu chưa cấu hình JWT_SECRET trong .env, dùng tạm một chuỗi mặc định
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("my_super_secret_key")
	}

	// 1. Tạo payload (claims) cho token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token hết hạn sau 24 giờ
	}

	// 2. Tạo token với thuật toán mã hóa HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Ký token bằng khóa bí mật
	return token.SignedString(jwtSecret)
}
