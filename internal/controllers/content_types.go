package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func ContentType(router *gin.Engine) {
	group := router.Group("texapp/content_type/")
	group.GET("", services.GetContentTypes)
}
