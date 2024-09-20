package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/_other/services"
)

func WebSocket(router *gin.Engine) {
	group := router.Group("texapp/ws")

	group.GET("/:token", services.HandleConnections)
}
