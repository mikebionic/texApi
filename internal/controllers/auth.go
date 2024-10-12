package controllers

import (
	"texApi/internal/services"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.Engine) {
	group := router.Group("texapp/auth/")

	group.GET("/login/", services.UserLogin)
	group.GET("/profile/", services.UserGetMe)
	group.GET("/logout/", services.Logout)
	group.GET("/register-request/", services.RegisterRequest)

	group.GET("/oauth/:provider/callback/", services.GetOAuthCallbackFunction)
	group.GET("/oauth/logout/:provider/", services.OAuthLogout)
	group.GET("/oauth/:provider/", services.OAuthProvider)
	group.GET("/oauth/testfront/", services.OAuthFront)
}
