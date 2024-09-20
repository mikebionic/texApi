package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/_other/services"
)

func Statuses(router *gin.Engine) {
	group := router.Group("texapp/statuses")

	group.GET("", services.GetStatuses)
}
