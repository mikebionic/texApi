package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/_other/services"
	"texApi/pkg/middlewares"
)

func Users(router *gin.Engine) {
	group := router.Group("texapp/users", middlewares.Guard)

	group.GET("", services.GetUsers)
	group.GET("/:id", services.GetUser)
	group.PUT("", services.UpdateUser)
	group.DELETE("/:id", services.DeleteUser)
}
