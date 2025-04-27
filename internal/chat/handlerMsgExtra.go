package chat

import (
	"fmt"
	"net/http"
	"strconv"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) SearchMessages(c *gin.Context) {
	userID := c.MustGet("id").(int)

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Search query is required", "Empty search query"))
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

	messages, err := h.repository.SearchMessages(userID, query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to search messages", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.FormatResponse("", messages))
}

func (h *APIHandler) PinMessage(c *gin.Context) {
	userID := c.MustGet("id").(int)

	var req struct {
		IsPinned  bool `json:"is_pinned"`
		MessageID int  `json:"message_id"`
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

	err = h.repository.PinMessage(req.MessageID, messageDetails.ConversationID, req.IsPinned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to pin/unpin message", err.Error()))
		return
	}

	messageDetails.IsPinned = &req.IsPinned
	notificationMsg := &Message{
		MessageCommon: MessageCommon{
			MessageType:    "pin",
			ConversationID: messageDetails.ConversationID,
			SenderID:       userID,
			Content:        fmt.Sprintf("Message has been %s", map[bool]string{true: "pinned", false: "unpinned"}[req.IsPinned]),
		},
		Extras: &map[string]interface{}{
			"message_details": messageDetails,
		},
	}

	go h.hub.RouteMessage(notificationMsg)

	c.JSON(http.StatusCreated, utils.FormatResponse("Message pin status updated", map[string]interface{}{
		"message_id": req.MessageID,
		"is_pinned":  req.IsPinned,
	}))
}

func (h *APIHandler) EditMessage(c *gin.Context) {
	userID := c.MustGet("id").(int)

	var req struct {
		Content   string `json:"content"`
		MessageID int    `json:"message_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Message content cannot be empty", ""))
		return
	}

	messageDetails, err := h.repository.GetMessageDetails(req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to get message details", err.Error()))
		return
	}

	err = h.repository.EditMessage(req.MessageID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to edit message", err.Error()))
		return
	}

	messageDetails.Content = req.Content
	isEdited := true
	messageDetails.IsEdited = &isEdited

	notificationMsg := &Message{
		MessageCommon: MessageCommon{
			MessageType:    "edit",
			ConversationID: messageDetails.ConversationID,
			SenderID:       userID,
			Content:        "Message edited",
		},
		Extras: &map[string]interface{}{
			"message_details": messageDetails,
		},
	}

	go h.hub.RouteMessage(notificationMsg)

	c.JSON(http.StatusCreated, utils.FormatResponse("Message edited successfully", messageDetails))
}

func (h *APIHandler) ReactToMessage(c *gin.Context) {
	userID := c.MustGet("id").(int)
	companyID := c.MustGet("companyID").(int)

	var req struct {
		Emoji     string `json:"emoji"`
		MessageID int    `json:"message_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	if req.Emoji == "" {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Emoji must be provided", ""))
		return
	}

	messageDetails, err := h.repository.GetMessageDetails(req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to get message details", err.Error()))
		return
	}

	err = h.repository.AddReaction(req.MessageID, userID, companyID, req.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to add reaction", err.Error()))
		return
	}

	reactions, err := h.repository.GetMessageReactions([]int{req.MessageID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to get reactions", err.Error()))
		return
	}

	reactionData := map[string]interface{}{
		"message_id":      req.MessageID,
		"user_id":         userID,
		"company_id":      companyID,
		"emoji":           req.Emoji,
		"conversation_id": messageDetails.ConversationID,
		"reactions":       reactions,
	}

	notificationMsg := &Message{
		MessageCommon: MessageCommon{
			MessageType:    "reaction",
			ConversationID: messageDetails.ConversationID,
			SenderID:       userID,
			Content:        "Reaction updated",
		},
		Extras: &reactionData,
	}

	go h.hub.RouteMessage(notificationMsg)

	c.JSON(http.StatusCreated, utils.FormatResponse("Reaction updated", reactionData))
}

func (h *APIHandler) GetPinnedMessages(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	pinnedMessages, err := h.repository.GetPinnedMessages(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch pinned messages", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.FormatResponse("", pinnedMessages))
}
