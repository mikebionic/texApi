package chat

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/database"
	"texApi/pkg/middlewares"
)

func Chat(router *gin.Engine) {
	chatRepository := NewRepository(database.DB)
	jwtSecret := []byte(config.ENV.ACCESS_KEY)

	ChatHub := NewHub()
	go ChatHub.Run()
	apiHandler := NewAPIHandler(chatRepository, ChatHub, jwtSecret)

	notificationGroup := router.Group(config.ENV.API_PREFIX+"/ws-notification/", middlewares.SysGuard)
	{
		notificationGroup.POST("/", apiHandler.SendDirectNotification)
	}

	group := router.Group(config.ENV.API_PREFIX+"/chat/", middlewares.Guard)
	convGroup := group.Group("/conversations/:id", apiHandler.ConversationAccessMiddleware())
	{
		convGroup.GET("/", apiHandler.GetConversation)
		convGroup.PUT("/", apiHandler.UpdateConversation)
		convGroup.DELETE("/", apiHandler.DeleteConversation)

		convGroup.GET("/pinned/", apiHandler.GetPinnedMessages)
		convGroup.POST("/pin/", apiHandler.PinMessage)
		convGroup.GET("/message/", apiHandler.GetMessages)
		convGroup.POST("/message/", apiHandler.SendMessage)
		convGroup.PUT("/message/", apiHandler.EditMessage)
		convGroup.DELETE("/message/owner/", apiHandler.DeleteMessageOfOwner)
		convGroup.POST("/message/react/", apiHandler.ReactToMessage)

		convGroup.POST("/member/manage/", apiHandler.AddRemoveConversationMembers)
		convGroup.PUT("/member/", apiHandler.UpdateConversationMember)
	}
	group.GET("/conversations/", apiHandler.GetConversations)
	group.POST("/conversations/", apiHandler.CreateConversation)
	group.GET("/search/", apiHandler.SearchMessages)

	wsHandler := NewWebSocketHandler(ChatHub, chatRepository, jwtSecret)
	wsRouteGroup := router.Group(config.ENV.API_PREFIX + "/ws/")
	wsRouteGroup.Use(middlewares.GuardURLParam)
	{
		wsRouteGroup.GET("/connect/", wsHandler.HandleWebSocket)
		wsRouteGroup.GET("/join/", wsHandler.HandleJoinConversation)
		wsRouteGroup.GET("/leave/", wsHandler.HandleLeaveConversation)
	}
}
