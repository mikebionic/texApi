package chat

import (
	"github.com/gin-gonic/gin"
	"log"
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

	group := router.Group(config.ENV.API_PREFIX + "/chat/")
	group.Use(middlewares.Guard)
	convGroup := group.Group("/conversations/:id", apiHandler.ConversationAccessMiddleware())
	{
		convGroup.GET("/", apiHandler.GetConversation)
		convGroup.GET("/pinned", apiHandler.GetPinnedMessages)
		convGroup.POST("/pin", apiHandler.PinMessage)
		convGroup.GET("/message/", apiHandler.GetMessages)
		convGroup.POST("/message/", apiHandler.SendMessage)
		convGroup.PUT("/message/", apiHandler.EditMessage)
		convGroup.DELETE("/message/", apiHandler.DeleteMessage)
		convGroup.POST("/message/react", apiHandler.ReactToMessage)
	}
	group.GET("/conversations", apiHandler.GetConversations)
	group.POST("/conversations", apiHandler.CreateConversation)
	group.GET("/search", apiHandler.SearchMessages)

	wsHandler := NewWebSocketHandler(ChatHub, chatRepository, jwtSecret)
	wsRouteGroup := router.Group(config.ENV.API_PREFIX + "/ws/")
	wsRouteGroup.Use(middlewares.Guard)

	wsRouteGroup.GET("/connect",
		func(c *gin.Context) { log.Printf("WebSocket connection attempt from %s", c.ClientIP()) },
		wsHandler.HandleWebSocket)
	wsRouteGroup.GET("/join",
		func(c *gin.Context) { log.Printf("Conversation join attempt from %s", c.ClientIP()) },
		wsHandler.HandleJoinConversation)
	wsRouteGroup.GET("/leave",
		func(c *gin.Context) { log.Printf("Conversation leave attempt from %s", c.ClientIP()) },
		wsHandler.HandleLeaveConversation)
}
