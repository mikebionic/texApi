package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Auth(router *gin.Engine) {
	group := router.Group("texapp/auth/")

	group.GET("/login/", services.UserLogin)
	group.GET("/profile/", services.UserGetMe)
	group.GET("/logout/", services.Logout)

	group.GET("/oauth/:provider/callback/", services.GetOAuthCallbackFunction)
	group.GET("/oauth/logout/:provider/", services.OAuthLogout)
	group.GET("/oauth/:provider/", services.OAuthProvider)
	group.GET("/oauth/testfront/", services.OAuthFront)
}
