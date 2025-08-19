package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"texApi/pkg/utils"
	"time"

	"github.com/gorilla/websocket"
)

// TODO NOTE: Move to ENV.Config ?
const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 4096
)

type Client struct {
	hub           *Hub
	conn          *websocket.Conn
	send          chan *Message
	userID        int
	companyID     int
	conversations []int
	onDisconnect  func()
}

func NewClient(hub *Hub, conn *websocket.Conn, userID, companyID int, conversations []int) *Client {
	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan *Message, 256),
		userID:        userID,
		companyID:     companyID,
		conversations: conversations,
		onDisconnect:  nil, // Optional callback
	}
}

func (c *Client) ReadPump(repository *Repository) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in ReadPump: %v", r)
			c.SendError("Internal error", "An unexpected error occurred")
			debug.PrintStack()
		}
		c.hub.unregister <- c
		if c.onDisconnect != nil {
			c.onDisconnect()
		}
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket close: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			c.SendError("Internal error", fmt.Sprintf("Invalid message format: %v", err))
			continue
		}

		if msg.ConversationID == 0 {
			c.SendError("Internal error", "Invalid message: missing conversation ID")
			continue
		}

		msg.SenderID = c.userID

		switch msg.Type {
		case MessageTypeMessage:
			c.handleSendMessage(msg, repository)

		case MessageTypeText:
			c.handleSendMessage(msg, repository)

		case MessageTypeMessageRead:
			c.handleMessageRead(msg, repository)

		case MessageTypeTyping:
			c.hub.RouteMessage(&msg)

		case MessageTypeSendingFile:
			c.hub.RouteMessage(&msg)

		case MessageTypeSticker:
			c.hub.RouteMessage(&msg)

		case MessageTypeNotification:
			c.hub.RouteMessage(&msg)

		default:
			c.SendError("Internal error", fmt.Sprintf("Unknown message type: %s", msg.Type))
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (c *Client) handleSendMessage(msg Message, repository *Repository) {
	if !repository.CanAccessConversation(c.userID, msg.ConversationID) {
		log.Printf("User %d cannot access conversation %d", c.userID, msg.ConversationID)
		return
	}

	msgID, err := repository.SaveMessage(&msg)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		return
	}

	msg.ID = msgID
	log.Printf("Message saved with ID: %d", msgID)

	c.hub.RouteMessage(&msg)

	// // backup function, or for offline users
	// go func() {
	// 	senderName := "Unknown"
	// 	if msg.SenderName != nil {
	// 		senderName = *msg.SenderName
	// 	}

	// 	if err := notification.SendNotificationToConversation(
	// 		msg.ConversationID,
	// 		msg.SenderID,
	// 		senderName,
	// 		msg.Content,
	// 	); err != nil {
	// 		log.Printf("Error sending Firebase notification for conversation %d: %v", msg.ConversationID, err)
	// 	}
	// }()
}

func (c *Client) handleMessageRead(msg Message, repository *Repository) {

	if !repository.CanAccessConversation(c.userID, msg.ConversationID) {
		c.SendError("Access denied", fmt.Sprintf("User %d cannot access conversation %d", c.userID, msg.ConversationID))
		return
	}

	err := repository.SetMessageRead(msg.ID, c.userID, msg.ConversationID)
	if err != nil {
		c.SendError("Failed to mark message as read", err.Error())
		return
	}

	c.hub.RouteMessage(&msg)
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()

		if c.onDisconnect != nil {
			c.onDisconnect()
		}

		c.hub.unregister <- c
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			json.NewEncoder(w).Encode(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendError(message string, errorMsg string) {
	errorResponse := utils.FormatErrorResponse(message, errorMsg)

	errorMessage := &Message{
		MessageCommon: MessageCommon{
			ConversationID: 0,
			SenderID:       0,
			Content:        message,
		},
	}

	errorMessage.Type = MessageTypeError

	extras := make(map[string]interface{})
	extras["error"] = errorResponse
	errorMessage.Extras = &extras

	select {
	case c.send <- errorMessage:
		log.Printf("Error message sent to user %d: %s", c.userID, message)
	default:
		log.Printf("Failed to send error to user %d (channel full)", c.userID)
	}
}
