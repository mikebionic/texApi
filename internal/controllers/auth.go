package controllers

import (
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.Engine) {
	group := router.Group("texapp/auth/")

	group.GET("/login/", services.UserLogin)
	group.GET("/profile/", middlewares.Guard, services.UserGetMe)
	group.POST("/profile/update/", middlewares.Guard, services.ProfileUpdate)
	group.GET("/logout/", services.Logout)

	group.POST("/register-request/", services.RegisterRequest)
	group.POST("/validate-otp/", services.ValidateOTP)
	group.POST("/register/", services.Register)
	group.POST("/password/forgot/", services.ForgotPassword)
	group.POST("/password/update/", services.UpdatePasswordOTP)
	group.POST("/refresh-token/", services.RefreshToken)

	group.GET("/oauth/:provider/callback/", services.GetOAuthCallbackFunction)
	group.GET("/oauth/logout/:provider/", services.OAuthLogout)
	group.GET("/oauth/:provider/", services.OAuthProvider)
	group.GET("/oauth/testfront/", services.OAuthFront)
}
