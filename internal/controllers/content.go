package controllers

import (
	"texApi/config"
	"texApi/internal/services"

	"github.com/gin-gonic/gin"
)

func Content(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/content/")

	group.GET("/", services.GetContents)
	group.GET("/:id", services.GetContent)
	group.POST("/", services.CreateContent)
	group.PUT("/:id", services.UpdateContent)
	group.POST("/update/:id", services.UpdateContent)
	group.DELETE("/:id", services.DeleteContent)

}
