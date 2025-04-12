package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"texApi/internal/dto"
	"texApi/internal/services"
	"texApi/pkg/fileUtils"
	"texApi/pkg/utils"
)

// APIHandler handles REST API requests for the chat system
type APIHandler struct {
	repository *Repository
	hub        *Hub
	jwtSecret  []byte
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(repository *Repository, hub *Hub, jwtSecret []byte) *APIHandler {
	return &APIHandler{
		repository: repository,
		hub:        hub,
		jwtSecret:  jwtSecret,
	}
}

type ApiError struct {
	Message string
	Status  int
}

func (e *ApiError) Error() string {
	return e.Message
}

func (h *APIHandler) GetConversations(c *gin.Context) {
	userID := c.MustGet("id").(int)

	query := `
		SELECT c.id, c.uuid, c.chat_type, c.title, c.description, 
		       c.image_url, c.theme_color, c.last_activity::TEXT, c.member_count, cm.unread_count,
		       (SELECT content FROM tbl_message WHERE id = c.last_message_id) as last_message
		FROM tbl_conversation c
		JOIN tbl_conversation_member cm ON c.id = cm.conversation_id
		WHERE cm.user_id = $1 AND cm.active = 1 AND cm.deleted = 0 AND c.active = 1 AND c.deleted = 0
		ORDER BY c.last_activity DESC
	`

	type ConversationListItem struct {
		ID           int     `json:"id"`
		UUID         string  `json:"uuid"`
		ChatType     string  `json:"chat_type"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		ImageURL     string  `json:"image_url"`
		ThemeColor   *string `json:"theme_color"`
		LastActivity *string `json:"last_activity"`
		MemberCount  int     `json:"member_count"`
		UnreadCount  *int    `json:"unread_count"`
		LastMessage  *string `json:"last_message"`
	}

	var conversations []ConversationListItem
	err := pgxscan.Select(context.Background(), h.repository.db, &conversations, query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to fetch conversations")
		return
	}

	c.JSON(http.StatusOK, utils.FormatResponse("", conversations))
}

func (h *APIHandler) CreateConversation(c *gin.Context) {
	userID := c.MustGet("id").(int)

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ChatType    string `json:"chat_type"`
		Members     []int  `json:"members"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Title is required", "Title == ''"))
		return
	}

	if req.ChatType != "direct" && req.ChatType != "group" && req.ChatType != "channel" {
		req.ChatType = "group" // Default to group
	}

	conversationID, err := h.repository.CreateConversation(
		userID, req.Title, req.Description, req.ChatType, req.Members,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create conversation", err.Error()))
		return
	}

	// Get created conversation details
	conversation, err := h.repository.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversation details", err.Error()))
		return
	}

	// TODO: Can this be from CLAIMS? add to jwt maybe?
	creatorName, _ := h.repository.GetCreatorName(userID)

	// For each invited member, send a WebSocket notification
	for _, memberID := range req.Members {
		if memberID != userID {
			inviteMsg := &Message{
				MessageType:    "system",
				SenderID:       userID,
				SenderName:     &creatorName,
				ConversationID: conversationID,
				Content:        fmt.Sprintf("You were added to %s", req.Title),
				Extras: &map[string]interface{}{
					"title":       req.Title,
					"description": req.Description,
					"chat_type":   req.ChatType,
					"members":     req.Members,
				},
			}

			h.hub.mu.RLock()
			for client := range h.hub.clients {
				if client.userID == memberID || client.userID == userID {
					select {
					case client.send <- inviteMsg:
						log.Printf("Sent conversation invite notification to user %d", memberID)
					default:
						log.Printf("Failed to send conversation invite to user %d", memberID)
					}
					break
				}
			}
			h.hub.mu.RUnlock()

			// And add to room
			for client := range h.hub.clients {
				if client.userID == memberID || client.userID == userID {
					h.hub.AddClientToRoom(client, conversationID)
				}
			}
		}
	}

	c.JSON(http.StatusCreated, utils.FormatResponse("", conversation))
}

func (h *APIHandler) GetConversation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	conversation, err := h.repository.GetConversation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversation details", err.Error()))
		return
	}

	members, err := h.repository.GetConversationMembers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversation members", err.Error()))
		return
	}

	response := struct {
		*Conversation
		Members []Member `json:"members"`
	}{
		Conversation: conversation,
		Members:      members,
	}
	c.JSON(http.StatusOK, utils.FormatResponse("", response))
}

func (h *APIHandler) GetMessages(c *gin.Context) {
	userID := c.MustGet("id").(int)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	limit := 50
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	messages, err := h.repository.GetConversationMessages(id, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch messages", err.Error()))
		return
	}

	// Get reactions for these messages
	query := `
		SELECT mr.message_id, mr.company_id, mr.emoji, mr.user_id,
		TRIM(COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '')) AS sender_name
		FROM tbl_message_reaction mr
		JOIN tbl_company p ON mr.company_id = p.id
		WHERE mr.message_id IN (
			SELECT id FROM tbl_message 
			WHERE conversation_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		)
	`

	var reactions []Reaction
	err = pgxscan.Select(context.Background(), h.repository.db, &reactions, query, id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch reactions", err.Error()))
		return
	}

	// Group reactions by message
	reactionsByMessage := make(map[int][]Reaction)
	for _, reaction := range reactions {
		reactionsByMessage[reaction.MessageID] = append(reactionsByMessage[reaction.MessageID], reaction)
	}

	type MessageWithReactions struct {
		MessageDetails
		Reactions []Reaction `json:"reactions"`
	}

	var messagesWithReactions []MessageWithReactions
	for _, msg := range messages {
		msgReactions := reactionsByMessage[msg.ID]
		messagesWithReactions = append(messagesWithReactions, MessageWithReactions{
			MessageDetails: msg,
			Reactions:      msgReactions,
		})
	}

	c.JSON(http.StatusOK, utils.FormatResponse("", messagesWithReactions))
}

func (h *APIHandler) SendMessage(c *gin.Context) {
	userID := c.MustGet("id").(int)
	companyID := c.MustGet("companyID").(int)

	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	jsonData := c.Request.FormValue("data")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("No message data provided", ""))
		return
	}

	var msg Message
	if err = json.Unmarshal([]byte(jsonData), &msg); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid JSON message data", err.Error()))
		return
	}

	msg.ConversationID = conversationID
	msg.SenderID = userID

	tx, err := h.repository.db.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	messageID, err := h.repository.SaveMessageTx(tx, &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to save message", err.Error()))
		return
	}
	var (
		mediaList      []dto.MediaCreate
		mediaRecord    dto.MediaCreate
		processedFiles []fileUtils.ProcessedFile
	)

	fileResults, err := services.ValidateAndProcessFiles(c, "message", "files")
	if err == nil && len(fileResults) > 0 {
		processedFiles, err = fileUtils.ProcessMediaFiles(fileResults)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to process media", err.Error()))
			return
		}

		for i, processedFile := range processedFiles {
			mediaRecord, err = services.SaveMediaToDatabase(c, tx, processedFile, userID, companyID, "message", messageID, i == 0)
			if err != nil {
				c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to save media to database", err.Error()))
				return
			}
			mediaList = append(mediaList, mediaRecord)
		}

		for _, media := range mediaList {
			generatedURL := fileUtils.GenerateMediaURL(media.UUID, media.Filename)
			mediaMain := dto.MediaMain{
				MediaCreate: media,
				URL:         generatedURL["url"],
				ThumbURL:    generatedURL["thumb_url"],
			}

			if msg.Media == nil {
				msg.Media = &[]dto.MediaMain{}
			}
			*msg.Media = append(*msg.Media, mediaMain)
		}
		msg.Extras = &map[string]interface{}{"files": processedFiles}
	}

	if err = tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to commit transaction", err.Error()))
		return
	}

	response := struct {
		ID int `json:"id"`
		*Message
		MediaIDs []map[string]interface{} `json:"media_ids"`
	}{
		ID:       messageID,
		Message:  &msg,
		MediaIDs: services.ExtractMediaIDs(mediaList),
	}

	go func() {
		h.hub.mu.RLock()
		defer h.hub.mu.RUnlock()
		h.hub.RouteMessage(&msg)
	}()

	c.JSON(http.StatusCreated, utils.FormatResponse("Message sent", response))
}
