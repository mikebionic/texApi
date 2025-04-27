package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func User(router *gin.Engine) {
	userGroup := router.Group(config.ENV.API_PREFIX + "/users/")
	{
		userGroup.GET("/", services.GetUserList)
		userGroup.GET("/:id", services.GetUser)
		userGroup.GET("/:id/rich", services.GetUserRichInfo)
		userGroup.POST("/", middlewares.GuardAdmin, services.CreateUser)
		userGroup.PUT("/:id", middlewares.GuardAdmin, services.UpdateUser)
		userGroup.DELETE("/:id", middlewares.GuardAdmin, services.DeleteUser)
	}
}
