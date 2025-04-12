package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Media(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/media/")
	{
		group.POST("upload/", middlewares.Guard, services.UploadFile)

		// General media routes
		group.GET("/:uuid/:filename", services.MediaFileHandler)
		group.GET("/:uuid/:filename/:thumb", services.MediaFileHandler)
	}
}
