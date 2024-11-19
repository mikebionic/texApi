package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/auth/")

	group.GET("/login/", services.UserLogin)
	group.GET("/profile/", middlewares.Guard, services.UserGetMe)
	group.POST("/profile/update/", middlewares.Guard, services.ProfileUpdate)
	group.GET("/logout/", middlewares.Guard, services.Logout)

	group.POST("/register-request/", services.RegisterRequest)
	group.POST("/validate-otp/", services.ValidateOTP)
	group.POST("/register/", services.Register)
	group.POST("/password/forgot/", services.ForgotPassword)
	group.POST("/password/update/", services.UpdatePasswordOTP)
	group.POST("/refresh-token/", services.RefreshToken)

	group.GET("/:provider", services.BeginOAuth)
	group.GET("/:provider/callback", services.CompleteOAuth)
}
