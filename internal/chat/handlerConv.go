package chat

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) GetConversations(c *gin.Context) {
	userID := c.MustGet("id").(int)

	conversations, err := h.repository.GetConversations(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversations", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.FormatResponse("", conversations))
}

func (h *APIHandler) CreateConversation(c *gin.Context) {
	userID := c.MustGet("id").(int)

	var req CreateConversation

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Title is required", "Title == ''"))
		return
	}

	if req.ChatType != "direct" && req.ChatType != "group" && req.ChatType != "channel" {
		req.ChatType = "direct" // Default to direct
	}

	conversationID, err := h.repository.CreateConversation(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create conversation", err.Error()))
		return
	}

	conversation, err := h.repository.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversation details", err.Error()))
		return
	}

	creatorName, _ := h.repository.GetCreatorName(userID)

	for _, memberID := range req.Members {
		if memberID != userID {
			inviteMsg := &Message{
				MessageCommon: MessageCommon{
					MessageType:    "system",
					SenderID:       userID,
					SenderName:     &creatorName,
					ConversationID: conversationID,
					Content:        fmt.Sprintf("You were added to %s", req.Title),
				},
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

			for client := range h.hub.clients {
				if client.userID == memberID || client.userID == userID {
					h.hub.AddClientToRoom(client, conversationID)
				}
			}
		}
	}

	c.JSON(http.StatusCreated, utils.FormatResponse("", conversation))
}

func (h *APIHandler) UpdateConversation(c *gin.Context) {
	userID := c.MustGet("id").(int)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	var req Conversation
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	err = h.repository.UpdateConversation(conversationID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update conversation", err.Error()))
		return
	}

	conversation, err := h.repository.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversation details", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.FormatResponse("Successfully updated conversation", conversation))
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

func (h *APIHandler) DeleteConversation(c *gin.Context) {
	userID := c.MustGet("id").(int)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	forEveryone, err := strconv.Atoi(c.Query("everyone"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	members, err := h.repository.GetConversationMembers(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversation members", err.Error()))
		return
	}

	conversation, err := h.repository.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to fetch conversation details", err.Error()))
		return
	}
	if *conversation.ChatType == "channel" {
		if userID != *conversation.CreatorID {
			err = h.repository.RemoveConversationMembers(conversationID, []int{userID})
			if err != nil {
				c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to remove messages for everyone of conversation channel", err.Error()))
				return
			}
			for client := range h.hub.clients {
				if client.userID == userID {
					h.hub.RemoveClientFromRoom(client, conversationID)
				}
			}
			c.JSON(http.StatusCreated, utils.FormatResponse("You've left the channel", conversation))
			return
		}
	}

	memberIDs := ExtractMemberIDs(members)
	if forEveryone > 0 {
		err = h.repository.RemoveMessages(conversationID, memberIDs, nil, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to remove messages for everyone of conversation", err.Error()))
			return
		}
		err = h.repository.RemoveConversationMembers(conversationID, memberIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to remove members for everyone of conversation", err.Error()))
			return
		}
		err = h.repository.RemoveConversation(conversationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to remove conversation", err.Error()))
			return
		}

		for _, member := range members {
			for client := range h.hub.clients {
				if client.userID == member.UserID || client.userID == userID {
					h.hub.RemoveClientFromRoom(client, conversationID)
				}
			}
		}

	} else {
		err = h.repository.RemoveMessages(conversationID, []int{userID}, nil, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to remove messages of conversation", err.Error()))
			return
		}
		err = h.repository.RemoveConversationMembers(conversationID, []int{userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to remove members for everyone of conversation", err.Error()))
			return
		}

		for client := range h.hub.clients {
			if client.userID == userID {
				h.hub.RemoveClientFromRoom(client, conversationID)
			}
		}
	}

	c.JSON(http.StatusCreated, utils.FormatResponse("Conversation deleted successfully!", conversationID))
	return
}

func (h *APIHandler) UpdateConversationMember(c *gin.Context) {
	requestingUserID := c.MustGet("id").(int)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}

	var req UpdateMemberRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	isAdmin, err := h.repository.IsConversationAdmin(requestingUserID, conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to check admin status", err.Error()))
		return
	}

	if !isAdmin {
		req.IsAdmin = nil
		req.Privileges = nil
		req.Nickname = nil
	}

	err = h.repository.UpdateConversationMember(conversationID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update conversation member", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.FormatResponse("Successfully updated conversation member", req))
}

func (h *APIHandler) AddRemoveConversationMembers(c *gin.Context) {
	requestingUserID := c.MustGet("id").(int)
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
		return
	}
	action := c.Query("action")

	var req MemberIDsRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request payload", err.Error()))
		return
	}

	isAdmin, err := h.repository.IsConversationAdmin(requestingUserID, conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to check admin status", err.Error()))
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, utils.FormatErrorResponse("Unauthorized", "Only admins can add members to this conversation"))
		return
	}

	if action == "add" {
		err = h.repository.AddConversationMembers(conversationID, req.MemberIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to add conversation members", err.Error()))
			return
		}

		for _, memberID := range req.MemberIDs {
			for client := range h.hub.clients {
				if client.userID == memberID {
					h.hub.AddClientToRoom(client, conversationID)
				}
			}
		}

	} else if action == "remove" {
		err = h.repository.RemoveConversationMembers(conversationID, req.MemberIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to remove conversation members", err.Error()))
			return
		}

		for _, memberID := range req.MemberIDs {
			for client := range h.hub.clients {
				if client.userID == memberID {
					h.hub.RemoveClientFromRoom(client, conversationID)
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid action", action))
		return
	}
	c.JSON(http.StatusCreated, utils.FormatResponse("Successfully managed members in conversation", req.MemberIDs))
}
