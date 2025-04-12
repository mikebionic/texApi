package chat

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"texApi/pkg/utils"
)

func (h *APIHandler) ConversationAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		conversationID, err := strconv.Atoi(idParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid conversation ID", err.Error()))
			return
		}

		userIDVal, exists := c.Get("id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", "User ID not found in context"))
			return
		}

		userID := userIDVal.(int)

		if !h.repository.CanAccessConversation(userID, conversationID) {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.FormatErrorResponse("You don't have access to this conversation", ""))
			return
		}

		c.Set("conversationID", conversationID)
		c.Set("userID", userID)

		c.Next()
	}
}
