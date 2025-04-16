package chat

import (
	"log"
	"sync"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Map of conversation IDs to sets of clients
	rooms map[int]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan *Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for concurrent access to maps
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[int]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	log.Printf("Client registered: UserID=%d", client.userID)

	// Add client to rooms (conversations)
	for _, conversationID := range client.conversations {
		if _, ok := h.rooms[conversationID]; !ok {
			h.rooms[conversationID] = make(map[*Client]bool)
		}
		h.rooms[conversationID][client] = true
		log.Printf("Client joined room: UserID=%d, ConversationID=%d", client.userID, conversationID)
	}
}

// unregisterClient removes a client from the hub and all rooms
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)

		// Remove from all rooms
		for conversationID, clients := range h.rooms {
			if _, ok := clients[client]; ok {
				delete(h.rooms[conversationID], client)
				log.Printf("Client left room: UserID=%d, ConversationID=%d", client.userID, conversationID)
			}

			// Clean up empty rooms
			if len(h.rooms[conversationID]) == 0 {
				delete(h.rooms, conversationID)
			}
		}

		close(client.send)
		log.Printf("Client unregistered: UserID=%d", client.userID)
	}
}

// broadcastMessage sends a message to all clients in a specific room
func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	log.Printf("Broadcasting message: ConversationID=%d, SenderID=%d, Type=%s",
		message.ConversationID, message.SenderID, message.MessageType)

	// Get clients in the conversation room
	clients, ok := h.rooms[message.ConversationID]
	if !ok {
		log.Printf("WARNING: No clients found in room for ConversationID=%d", message.ConversationID)
		return
	}

	log.Printf("Found %d clients in room", len(clients))

	// Send message to all clients in the room
	for client := range clients {
		if client.userID == message.SenderID && message.MessageType == "direct" {
			continue
		}
		log.Printf("Attempting to send message to UserID=%d", client.userID)

		select {
		case client.send <- message:
			log.Printf("Message sent to UserID=%d", client.userID)
		default:
			log.Printf("Failed to send message to UserID=%d (send channel full)", client.userID)
			// If client's send buffer is full, we assume it's disconnected
			h.mu.RUnlock()
			go h.unregisterClient(client)
			h.mu.RLock()
		}
	}
}

func (h *Hub) AddClientToRoom(client *Client, conversationID int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[conversationID]; !ok {
		h.rooms[conversationID] = make(map[*Client]bool)
	}
	h.rooms[conversationID][client] = true

	// Add to client's tracked conversations list if not already there
	found := false
	for _, id := range client.conversations {
		if id == conversationID {
			found = true
			break
		}
	}

	if !found {
		client.conversations = append(client.conversations, conversationID)
	}

	log.Printf("Client added to room: UserID=%d, ConversationID=%d", client.userID, conversationID)
}

func (h *Hub) RemoveClientFromRoom(client *Client, conversationID int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[conversationID]; ok {
		delete(clients, client)
		log.Printf("Client removed from room: UserID=%d, ConversationID=%d", client.userID, conversationID)

		if len(clients) == 0 {
			delete(h.rooms, conversationID)
		}
	}

	// Remove from client's tracked conversations
	for i, id := range client.conversations {
		if id == conversationID {
			client.conversations = append(client.conversations[:i], client.conversations[i+1:]...)
			break
		}
	}
}

func (h *Hub) RouteMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Get clients in the specific conversation room
	clients, ok := h.rooms[message.ConversationID]
	if !ok {
		log.Printf("No clients found for conversation %d", message.ConversationID)
		return
	}

	// Send message only to clients in this specific conversation
	for client := range clients {
		// Avoid sending the message back to the sender
		if client.userID != message.SenderID {
			select {
			case client.send <- message:
				log.Printf("Message routed to user %d", client.userID)
			default:
				log.Printf("Failed to route message to user %d (channel full)", client.userID)
			}
		}
	}
}

func (h *Hub) GetOnlineUsersInConversation(conversationID int) []int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var onlineUsers []int
	if clients, ok := h.rooms[conversationID]; ok {
		for client := range clients {
			onlineUsers = append(onlineUsers, client.userID)
		}
	}

	return onlineUsers
}

// Online and Offile status sending
func (h *Hub) TrackUserStatus(client *Client, status bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conversationID := range client.conversations {
		if clients, ok := h.rooms[conversationID]; ok {
			statusMessage := &Message{
				MessageType:    "user_status",
				ConversationID: conversationID,
				OnlineStatus: &OnlineStatus{
					UserID:   client.userID,
					IsOnline: status,
				},
			}

			// Broadcast status to all clients in the room
			for otherClient := range clients {
				select {
				case otherClient.send <- statusMessage:
				// Successfully sent
				default:
					log.Printf("Could not send status message to client %d", otherClient.userID)
				}
			}
		}
	}
}

// TODO: Document this but no usages
// Periodic cleanup in hub
func (h *Hub) StartConnectionCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			h.cleanupStaleConnections()
		}
	}()
}
func (h *Hub) cleanupStaleConnections() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conversationID, clients := range h.rooms {
		for client := range clients {
			// Check if send channel is full or closed
			select {
			case <-client.send:
				// Active connection, do nothing
			default:
				// Remove stale client
				delete(clients, client)
				close(client.send)
				log.Printf("Cleaned up stale client from conversation %d", conversationID)
			}
		}

		if len(clients) == 0 {
			delete(h.rooms, conversationID)
		}
	}
}
