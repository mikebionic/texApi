package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Content(router *gin.Engine) {
	group := router.Group("texapp/content/")

	group.GET("", services.GetContents)
	group.GET("/:id", services.GetContent)
	group.POST("", middlewares.Guard, services.CreateContent)
	group.DELETE("/:id", middlewares.Guard, services.DeleteContent)

}
