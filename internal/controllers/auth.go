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
}
