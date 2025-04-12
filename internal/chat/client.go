package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

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
	hub       *Hub
	conn      *websocket.Conn
	send      chan *Message
	userID    int
	companyID int
	//TODO: can add more info in JWT for "Claims"
	// Username      string `json:"username"`
	conversations []int // List of conversation IDs this client is part of
	onDisconnect  func()
}

// NewClient creates a new client
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

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump(repository *Repository) {
	defer func() {
		c.hub.unregister <- c
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

		// Parse incoming message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// Validate message fields
		if msg.ConversationID == 0 {
			log.Printf("Invalid message: missing conversation ID")
			continue
		}

		// Set the sender ID from the authenticated user
		msg.SenderID = c.userID

		if msg.TypingStatus {
			// Don't save typing indicators to database
			c.hub.RouteMessage(&msg)
			continue
		}

		// Validate conversation access
		if !repository.CanAccessConversation(c.userID, msg.ConversationID) {
			log.Printf("User %d cannot access conversation %d", c.userID, msg.ConversationID)
			continue
		}

		// Save message to database
		msgID, err := repository.SaveMessage(&msg)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Store message ID for potential replies or reactions
		log.Printf("Message saved with ID: %d", msgID)

		//// Broadcast message to all clients in the conversation
		//c.hub.broadcast <- &msg

		// Route message through hub to specific conversation clients
		c.hub.RouteMessage(&msg)
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()

		// Call onDisconnect callback if set
		if c.onDisconnect != nil {
			c.onDisconnect()
		}

		// Unregister from hub
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
