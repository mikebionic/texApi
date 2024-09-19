package controllers

import (
	"texApi/pkg/services"

	"github.com/gin-gonic/gin"
)

func Content(router *gin.Engine) {
	group := router.Group("texapp/content")

	group.GET("", services.GetContents)
}
