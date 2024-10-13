package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Content(router *gin.Engine) {
	group := router.Group("texapp/content/")

	group.GET("", services.GetContents)
	group.GET("/:id", services.GetContent)
	group.POST("", services.CreateContent)
	group.PUT("/:id", services.UpdateContent)
	group.DELETE("/:id", services.DeleteContent)

}
