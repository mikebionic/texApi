package controllers

import (
	"texApi/pkg/middlewares"
	"texApi/pkg/services"

	"github.com/gin-gonic/gin"
)

func Users(router *gin.Engine) {
	group := router.Group("texapp/users", middlewares.Guard)

	group.GET("", services.GetUsers)
	group.GET("/:id", services.GetUser)
	group.PUT("", services.UpdateUser)
	group.DELETE("/:id", services.DeleteUser)
}
