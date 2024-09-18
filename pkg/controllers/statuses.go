package controllers

import (
	"texApi/pkg/services"

	"github.com/gin-gonic/gin"
)

func Statuses(router *gin.Engine) {
	group := router.Group("texapp/statuses")

	group.GET("", services.GetStatuses)
}
