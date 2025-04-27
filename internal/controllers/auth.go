package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/auth/")
	{
		group.POST("/login/", services.UserLogin)

		group.POST("/register-request/", services.RegisterRequest)
		group.POST("/validate-otp/", services.ValidateOTP)
		group.POST("/register/", services.Register)
		group.POST("/password/forgot/", services.ForgotPassword)
		group.POST("/password/update/", services.UpdatePasswordOTP)
		group.POST("/refresh-token/", services.RefreshToken)

		// OTP login routes
		group.POST("/otp-login/request/", services.OTPLoginRequest)
		group.POST("/otp-login/validate/", services.OTPLogin)

		group.GET("/google/", services.BeginOAuth)
		group.GET("/google/callback/", services.CompleteOAuth)
		group.POST("/google/mobile/", services.BeginOAuthMobile)

		user := group.Use(middlewares.Guard)
		{
			user.GET("/user/", services.UserGetMe)
			user.POST("/user/update/", services.UserUpdate)
			user.POST("/logout/", services.Logout)
			user.POST("/logout-all/", services.LogoutAllSessions)

			// Session management
			user.GET("/sessions/", services.ListUserSessions)
			user.DELETE("/sessions/:id/", services.RevokeSession)
		}

		// Admin routes (with admin middlewares)
		admin := group.Group("/admin", middlewares.GuardAdmin)
		{
			admin.GET("/sessions", services.ListAllSessions)
			admin.DELETE("/sessions/:id", services.AdminRevokeSession)
			admin.DELETE("/user-sessions/:user_id", services.AdminRevokeUserSessions)
		}
	}
}
