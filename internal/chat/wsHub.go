package chat

import (
	"log"
	"sync"
	"time"
)

type Hub struct {
	clients    map[*Client]bool
	rooms      map[int]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

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
			h.RouteMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	log.Printf("Client registered: UserID=%d", client.userID)

	for _, conversationID := range client.conversations {
		if _, ok := h.rooms[conversationID]; !ok {
			h.rooms[conversationID] = make(map[*Client]bool)
		}
		h.rooms[conversationID][client] = true
		log.Printf("Client joined room: UserID=%d, ConversationID=%d", client.userID, conversationID)
	}
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)

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

func (h *Hub) AddClientToRoom(client *Client, conversationID int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[conversationID]; !ok {
		h.rooms[conversationID] = make(map[*Client]bool)
	}
	h.rooms[conversationID][client] = true

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

	clients, ok := h.rooms[message.ConversationID]
	if !ok {
		log.Printf("No clients found for conversation %d", message.ConversationID)
		return
	}

	for client := range clients {
		if client.userID != message.SenderID {
			select {
			case client.send <- message:
				log.Printf("Message routed to user %d", client.userID)
			default:
				log.Printf("Failed to route message to user %d (channel full)", client.userID)
				go h.unregisterClient(client)
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
				MessageCommon: MessageCommon{
					MessageType:    "user_status",
					ConversationID: conversationID,
				},
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
