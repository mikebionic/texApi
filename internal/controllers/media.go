package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Media(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/media/")
	{
		group.POST("upload/", middlewares.Guard, services.UploadFile)

		group.GET("/:uuid/:filename", services.MediaFileHandler)
		group.GET("/:uuid/:filename/:thumb", services.MediaFileHandler)
	}
}
