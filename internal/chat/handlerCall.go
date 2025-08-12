package chat

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"texApi/config"
	"texApi/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) CreateCallRoom(c *gin.Context) {
	userID := c.MustGet("id").(int)
	userIDStr := strconv.Itoa(userID)
	profileID := strconv.Itoa(c.MustGet("companyID").(int))

	var req CreateCallRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	if !h.repository.CanAccessConversation(userID, req.ConversationID) {
		c.JSON(http.StatusForbidden, utils.FormatErrorResponse("Access denied to this conversation", ""))
		return
	}

	if !contains(req.UserIDs, userIDStr) {
		req.UserIDs = append(req.UserIDs, userIDStr)
	}
	if !contains(req.ProfileIDs, profileID) {
		req.ProfileIDs = append(req.ProfileIDs, profileID)
	}

	if req.Duration == 0 {
		req.Duration = 60 // 1 hour default
	}
	if req.MaxUser == 0 {
		req.MaxUser = len(req.UserIDs)
	}

	hex, err := h.generateUniqueHex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to generate unique hex", err.Error()))
		return
	}

	callRoom, err := h.repository.CreateCallRoom(&CallRoom{
		ConversationID: req.ConversationID,
		MaxUser:        req.MaxUser,
		UserIDs:        req.UserIDs,
		ProfileIDs:     req.ProfileIDs,
		Title:          req.Title,
		Hex:            hex,
		Duration:       req.Duration,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create call room", err.Error()))
		return
	}

	callRoom.JoinURL = fmt.Sprintf("%s/%s/call-room/join/%s",
		config.ENV.API_SERVER_URL, config.ENV.API_PREFIX, callRoom.UUID)
	callRoom.JitsiURL = fmt.Sprintf("%s/%s", config.ENV.JITSI_URL, callRoom.Hex)

	var msg Message

	msg.ConversationID = req.ConversationID
	msg.SenderID = userID
	msg.MessageType = "notification"
	msg.Content = callRoom.JitsiURL
	msg.CreatedAt = time.Now()

	go func() {
		h.hub.mu.RLock()
		defer h.hub.mu.RUnlock()
		h.hub.RouteMessage(&msg)
	}()

	c.JSON(http.StatusCreated, utils.FormatResponse("Call room created", callRoom))
}

func (h *APIHandler) JoinCallRoom(c *gin.Context) {
	userID := c.MustGet("id").(int)
	userIDStr := strconv.Itoa(userID)
	uuid := c.Param("uuid")

	profileID := strconv.Itoa(c.MustGet("companyID").(int))

	callRoom, err := h.repository.GetCallRoomByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.FormatErrorResponse("Call room not found", err.Error()))
		return
	}

	if callRoom.Active != 1 || callRoom.Deleted != 0 {
		c.JSON(http.StatusForbidden, utils.FormatErrorResponse("Call room is not active", ""))
		return
	}

	hasAccess := contains(callRoom.UserIDs, userIDStr) || contains(callRoom.ProfileIDs, profileID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, utils.FormatErrorResponse("Access denied", ""))
		return
	}

	callRoom.JoinURL = fmt.Sprintf("%s/%s/call-room/join/%s",
		config.ENV.API_SERVER_URL, config.ENV.API_PREFIX, callRoom.UUID)
	callRoom.JitsiURL = fmt.Sprintf("%s/%s", config.ENV.JITSI_URL, callRoom.Hex)

	c.JSON(http.StatusOK, utils.FormatResponse("Call room access granted", callRoom))
}

func (h *APIHandler) EndCallRoom(c *gin.Context) {
	userID := c.MustGet("id").(int)
	userIDStr := strconv.Itoa(userID)
	uuid := c.Param("uuid")

	callRoom, err := h.repository.GetCallRoomByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.FormatErrorResponse("Call room not found", err.Error()))
		return
	}

	if !contains(callRoom.UserIDs, userIDStr) {
		c.JSON(http.StatusForbidden, utils.FormatErrorResponse("Access denied", ""))
		return
	}

	err = h.repository.EndCallRoom(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to end call room", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.FormatResponse("Call room ended", nil))
}

func (r *Repository) CreateCallRoom(callRoom *CallRoom) (*CallRoom, error) {
	query := `
		INSERT INTO tbl_call_room (
			conversation_id, max_user, user_ids, profile_ids, 
			title, hex, duration, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, uuid, created_at, updated_at, active, deleted
	`

	err := r.db.QueryRow(
		context.Background(), query,
		callRoom.ConversationID,
		callRoom.MaxUser,
		callRoom.UserIDs,
		callRoom.ProfileIDs,
		callRoom.Title,
		callRoom.Hex,
		callRoom.Duration,
	).Scan(
		&callRoom.ID,
		&callRoom.UUID,
		&callRoom.CreatedAt,
		&callRoom.UpdatedAt,
		&callRoom.Active,
		&callRoom.Deleted,
	)

	return callRoom, err
}

func (r *Repository) GetCallRoomByUUID(uuid string) (*CallRoom, error) {
	query := `
		SELECT id, uuid, conversation_id, max_user, user_ids, profile_ids,
           title, hex, duration, created_at, updated_at, active, deleted
		FROM tbl_call_room
		WHERE uuid = $1
	`

	var callRoom CallRoom
	err := r.db.QueryRow(context.Background(), query, uuid).Scan(
		&callRoom.ID,
		&callRoom.UUID,
		&callRoom.ConversationID,
		&callRoom.MaxUser,
		&callRoom.UserIDs,
		&callRoom.ProfileIDs,
		&callRoom.Title,
		&callRoom.Hex,
		&callRoom.Duration,
		&callRoom.CreatedAt,
		&callRoom.UpdatedAt,
		&callRoom.Active,
		&callRoom.Deleted,
	)

	return &callRoom, err
}

func (r *Repository) EndCallRoom(uuid string) error {
	query := `
		UPDATE tbl_call_room 
		SET active = 0, deleted = 1, updated_at = CURRENT_TIMESTAMP 
		WHERE uuid = $1
	`

	_, err := r.db.Exec(context.Background(), query, uuid)
	return err
}

func (r *Repository) HexExists(hex string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tbl_call_room WHERE hex = $1)`
	var exists bool
	err := r.db.QueryRow(context.Background(), query, hex).Scan(&exists)
	return exists, err
}

func (h *APIHandler) generateUniqueHex() (string, error) {
	for attempts := 0; attempts < 10; attempts++ {
		// Generate 30-character hex (15 bytes)
		bytes := make([]byte, 15)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		hex := hex.EncodeToString(bytes)

		// Check hex among last calls
		exists, err := h.repository.HexExists(hex)
		if err != nil {
			return "", err
		}
		if !exists {
			return hex, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique hex after 10 attempts")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
