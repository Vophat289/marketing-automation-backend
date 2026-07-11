package routes

import (
	"marketing-automation-backend/controllers"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes định nghĩa các đường dẫn API liên quan đến User
func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/users")
	{
		// Định nghĩa route POST /api/users/register
		userGroup.POST("/register", controllers.Register)
		
		//Lấy danh sách người dùng
		userGroup.GET("/listusers", controllers.ListUsers)

		//Post đăng nhập user
		userGroup.POST("/login", controllers.LoginUser)
	}
}
