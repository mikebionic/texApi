package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/_other/services"
	"texApi/pkg/middlewares"
)

func Auth(router *gin.Engine) {
	group := router.Group("texapp/auth")

	group.POST("/admin/login", services.AdminLogin)
	// group.POST("/user/register", services.CreateUser)
	group.POST("/user/login", services.UserLogin)
	group.GET("/user/me", middlewares.Guard, services.UserGetMe)
	// group.POST("/user/verify", services.UserVerify)
	// group.POST("/user/forget-password", services.UserForgetPassword)
	// group.POST("/user/new-password", services.UserNewPassword)
	group.POST("/worker/login", services.WorkerLogin)
	group.POST("/refresh-token", services.RefreshToken)
	group.PUT(
		"/refresh-notification-token", middlewares.Guard,
		services.RefreshNotificationToken,
	)
}
