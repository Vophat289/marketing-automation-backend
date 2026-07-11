package models

import "gorm.io/gorm"

// BlacklistedToken lưu trữ các token đã bị đăng xuất
type BlacklistedToken struct {
	gorm.Model
	Token string `gorm:"uniqueIndex"`
}
