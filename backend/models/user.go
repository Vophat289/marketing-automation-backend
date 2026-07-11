package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name         string
	Email        string `gorm:"unique"`
	Password     string
	TokenVersion int `gorm:"default:1"` // Dùng để vô hiệu hóa token cũ khi đổi mật khẩu
}
