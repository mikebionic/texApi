package controllers

import (
	"texApi/pkg/services"

	"github.com/gin-gonic/gin"
)

func WebSocket(router *gin.Engine) {
	group := router.Group("texapp/ws")

	group.GET("/:token", services.HandleConnections)
}
