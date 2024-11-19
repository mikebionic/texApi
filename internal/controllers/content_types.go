package controllers

import (
	"texApi/config"
	"texApi/internal/services"

	"github.com/gin-gonic/gin"
)

func ContentType(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/content_type/")
	group.GET("/", services.GetContentTypes)
	group.GET("/:id", services.GetContentTypes)
}
