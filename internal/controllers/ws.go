package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/database"
	"texApi/internal/chat"
	"texApi/pkg/middlewares"
)

func WS(router *gin.Engine) {
	chatRepository := chat.NewRepository(database.DB)
	jwtSecret := []byte(config.ENV.ACCESS_KEY)
	chatHub := chat.GetHub()

	wsHandler := chat.NewWebSocketHandler(chatHub, chatRepository, jwtSecret)
	wsRouteGroup := router.Group(config.ENV.API_PREFIX + "/ws/")
	wsRouteGroup.Use(middlewares.GuardURLParam)
	{
		wsRouteGroup.GET("/connect/", wsHandler.HandleWebSocket)
		wsRouteGroup.GET("/join/", wsHandler.HandleJoinConversation)
		wsRouteGroup.GET("/leave/", wsHandler.HandleLeaveConversation)
	}

}
