package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub        *Hub
	repository *Repository
	upgrader   websocket.Upgrader
	jwtSecret  []byte
}

func NewWebSocketHandler(hub *Hub, repository *Repository, jwtSecret []byte) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		repository: repository,
		jwtSecret:  jwtSecret,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				//TODO: Add proper origin check in production
				// origin := r.Header.Get("Origin")
				// return origin == "http://localhost:3000"
				return true // Accept all origins
			},
		},
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userID := c.MustGet("id").(int)
	companyID := c.MustGet("companyID").(int)

	conversations, err := h.repository.GetUserConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error loading user conversations", err.Error()))
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		c.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Invalid token", err.Error()))
		return
	}

	client := NewClient(h.hub, conn, userID, companyID, conversations)
	h.hub.TrackUserStatus(client, true)
	h.hub.register <- client

	for _, conversationID := range conversations {
		h.hub.AddClientToRoom(client, conversationID)
		log.Printf("Added client UserID=%d to ConversationID=%d", client.userID, conversationID)
	}

	go client.ReadPump(h.repository)
	go client.WritePump()

	client.onDisconnect = func() {
		h.hub.TrackUserStatus(client, false)
	}
}

//// TODO NOTE: API endpoint, unused
//func (h *WebSocketHandler) MessageRead(c *gin.Context) {
//	userID := c.MustGet("id").(int)
//
//	conversationID, err := strconv.Atoi(c.Param("conversation_id"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
//		return
//	}
//
//	messageID, err := strconv.Atoi(c.Param("message_id"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid message ID", err.Error()))
//		return
//	}
//
//	if !h.repository.CanAccessConversation(userID, conversationID) {
//		c.JSON(http.StatusForbidden, utils.FormatErrorResponse("Access denied to this conversation", ""))
//		return
//	}
//
//	err = h.repository.SetMessageRead(messageID, userID, conversationID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to mark message as read", err.Error()))
//		return
//	}
//
//	readMsg := Message{
//		MessageCommon: MessageCommon{
//			ID:             &messageID,
//			ConversationID: conversationID,
//			SenderID:       userID,
//		},
//	}
//	msgType := MessageType("message_read")
//	readMsg.Type = &msgType
//
//	h.hub.RouteMessage(&readMsg)
//
//	c.JSON(http.StatusOK, gin.H{
//		"status":          "message_read",
//		"conversation_id": conversationID,
//		"message_id":      messageID,
//	})
//}

// TODO: this is unused, because Connect already joins User to all conversations
// HandleJoinConversation adds a client to a conversation room
func (h *WebSocketHandler) HandleJoinConversation(c *gin.Context) {
	userID := c.MustGet("id").(int)
	conversationID, err := strconv.Atoi(c.Param("conversation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	// Check if user can access the conversation
	if !h.repository.CanAccessConversation(userID, conversationID) {
		c.JSON(http.StatusForbidden, utils.FormatErrorResponse("Access denied to this conversation", err.Error()))
		return
	}

	// Find client in hub
	var targetClient *Client
	h.hub.mu.RLock()
	for client := range h.hub.clients {
		if client.userID == userID {
			targetClient = client
			break
		}
	}
	h.hub.mu.RUnlock()

	if targetClient == nil {
		c.JSON(http.StatusNotFound, utils.FormatErrorResponse("Client not connected", ""))
		return
	}

	h.hub.AddClientToRoom(targetClient, conversationID)
	c.Writer.WriteHeader(http.StatusOK)

	json.NewEncoder(c.Writer).Encode(map[string]string{
		"status":          "joined",
		"conversation_id": strconv.Itoa(conversationID),
	})
}

func (h *WebSocketHandler) HandleLeaveConversation(c *gin.Context) {
	userID := c.MustGet("id").(int)
	conversationID, err := strconv.Atoi(c.Param("conversation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	var targetClient *Client
	h.hub.mu.RLock()
	for client := range h.hub.clients {
		if client.userID == userID {
			targetClient = client
			break
		}
	}
	h.hub.mu.RUnlock()

	if targetClient == nil {
		c.JSON(http.StatusNotFound, utils.FormatErrorResponse("Client not connected", ""))
		return
	}

	h.hub.RemoveClientFromRoom(targetClient, conversationID)

	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(map[string]string{
		"status":          "left",
		"conversation_id": strconv.Itoa(conversationID),
	})
}

func (h *WebSocketHandler) GetOnlineUsers(conversationID int) []int {
	var onlineUsers []int

	h.hub.mu.RLock()
	defer h.hub.mu.RUnlock()

	if clients, ok := h.hub.rooms[conversationID]; ok {
		for client := range clients {
			onlineUsers = append(onlineUsers, client.userID)
		}
	}

	return onlineUsers
}
