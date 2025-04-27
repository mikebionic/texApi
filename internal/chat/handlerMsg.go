package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"texApi/internal/dto"
	"texApi/internal/services"
	"texApi/pkg/fileUtils"
	"texApi/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

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

	messageIDs := make([]int, len(messages))
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	reactions, err := h.repository.GetMessageReactions(messageIDs)

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
	msg.CreatedAt = time.Now()

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
		MediaIDs []map[string]interface{} `json:"media_ids,omitempty"`
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

// Only conversation admin or a message sender can delete message route
func (h *APIHandler) DeleteMessageOfOwner(c *gin.Context) {
	userID := c.MustGet("id").(int)
	var req struct {
		MessageID int `json:"message_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	messageDetails, err := h.repository.GetMessageDetails(req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to get message details", err.Error()))
		return
	}

	isAdmin, err := h.repository.IsConversationAdmin(userID, messageDetails.ConversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to check admin status", err.Error()))
		return
	}

	err = h.repository.DeleteMessageOfOwner(req.MessageID, userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete message", err.Error()))
		return
	}

	notificationMsg := &Message{
		MessageCommon: MessageCommon{
			MessageType:    "delete",
			ConversationID: messageDetails.ConversationID,
			SenderID:       userID,
			Content:        "Message deleted",
		},
		Extras: &map[string]interface{}{"message_id": req.MessageID},
	}

	go h.hub.RouteMessage(notificationMsg)

	c.JSON(http.StatusCreated, utils.FormatResponse("Message deleted successfully", map[string]interface{}{
		"message_id": req.MessageID,
	}))
}
