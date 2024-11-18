package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
)

func Media(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/media/")
	group.POST("upload/", services.UploadFile)
}
