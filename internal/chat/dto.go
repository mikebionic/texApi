package chat

import (
	"encoding/json"
	"texApi/internal/dto"
	"time"
)

const (
	MessageTypeText         string = "text"
	MessageTypeMessage      string = "message"
	MessageTypeJoin         string = "user_joined"
	MessageTypeLeave        string = "user_left"
	MessageTypeError        string = "error"
	MessageTypeNotification string = "notification"
	MessageTypeDisconnect   string = "disconnect"
	MessageTypeTyping       string = "typing"
	MessageTypeSendingFile  string = "sending_file"
	MessageTypeSticker      string = "choosing_sticker"
	MessageTypeMessageRead  string = "message_read"
)

type Message struct {
	MessageCommon
	Type          string                  `json:"type"`
	ForwardedFrom *int                    `json:"forwarded_from_id"`
	Extras        *map[string]interface{} `json:"extras"` // For additional metadata
	OnlineStatus  *OnlineStatus           `json:"online_status,omitempty"`
}

type MessageCommon struct {
	ID             int              `json:"id"`
	ConversationID int              `json:"conversation_id"`
	SenderID       int              `json:"sender_id"`
	MessageType    string           `json:"message_type"`
	Content        string           `json:"content"`
	ReplyToID      *int             `json:"reply_to_id"`
	MediaID        *int             `json:"media_id"`
	StickerID      *int             `json:"sticker_id"`
	IsSilent       *bool            `json:"is_silent"`
	CreatedAt      time.Time        `json:"created_at"`
	SenderName     *string          `json:"sender_name"`
	Media          *[]dto.MediaMain `json:"media,omitempty"`
}

type OnlineStatus struct {
	UserID   int  `json:"user_id"`
	IsOnline bool `json:"is_online"`
}

type MessageDetails struct {
	MessageCommon
	UUID              string        `json:"uuid"`
	ConversationTitle *string       `json:"conversation_title,omitempty"`
	ConversationType  *string       `json:"conversation_type,omitempty"`
	ForwardedFromID   *int          `json:"forwarded_from_id,omitempty"`
	IsEdited          *bool         `json:"is_edited"`
	IsPinned          *bool         `json:"is_pinned"`
	IsDelivered       *bool         `json:"is_delivered"`
	SenderAvatar      *string       `json:"sender_avatar"`
	EditedAt          *time.Time    `json:"edited_at"`
	ReadAt            *time.Time    `json:"read_at"`
	ReadBy            *[]UsersJSONB `json:"read_by"`
	DeletedFor        *[]UsersJSONB `json:"deleted_for"`
	UpdatedAt         *time.Time    `json:"updated_at"`
	Active            *int          `json:"active,omitempty"`
	Deleted           *int          `json:"deleted,omitempty"`
	Reactions         *[]Reaction   `json:"reactions,omitempty"`
}

type Reaction struct {
	MessageID  int     `json:"message_id"`
	UserID     int     `json:"user_id"`
	CompanyID  int     `json:"company_id"`
	Emoji      string  `json:"emoji"`
	SenderName *string `json:"sender_name,omitempty"`
}

type CreateConversation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ChatType    string `json:"chat_type"`
	Members     []int  `json:"members"`
	IsPublic    *bool  `json:"is_public"`

	ThemeColor         *string `json:"theme_color"`
	BackgroundImageURL *string `json:"background_image_url"`
	BackgroundBlur     *int    `json:"background_blur"`
}
type UpdateMemberRequest struct {
	UserID           int        `json:"user_id" binding:"required"`
	Nickname         *string    `json:"nickname,omitempty"`
	IsAdmin          *bool      `json:"is_admin,omitempty"`
	Privileges       *[]string  `json:"privileges,omitempty"`
	NotificationPref *string    `json:"notification_preference,omitempty"`
	MutedUntil       *time.Time `json:"muted_until,omitempty"` // ISO 8601 format
}

type MemberIDsRequest struct {
	MemberIDs []int `json:"member_ids" binding:"required"`
}
type Conversation struct {
	ID                 int        `json:"id"`
	UUID               string     `json:"uuid"`
	ChatType           *string    `json:"chat_type"`
	Title              *string    `json:"title"`
	Description        *string    `json:"description"`
	CreatorID          *int       `json:"creator_id"`
	ThemeColor         *string    `json:"theme_color"`
	ImageURL           *string    `json:"image_url"`
	BackgroundImageURL *string    `json:"background_image_url"`
	BackgroundBlur     *int       `json:"background_blur"`
	LastMessageID      *int       `json:"last_message_id"`
	MemberCount        *int       `json:"member_count"`
	MessageCount       *int       `json:"message_count"`
	AutoDeleteDuration *int       `json:"auto_delete_duration"`
	InviteToken        *string    `json:"invite_token"`
	IsPublic           *bool      `json:"is_public"`
	PublicURL          *string    `json:"public_url"`
	LastActivity       *time.Time `json:"last_activity"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	Active             *int       `json:"active"`
	Deleted            *int       `json:"deleted"`

	UnreadCount *int     `json:"unread_count"`
	LastMessage *string  `json:"last_message"`
	MemberIDs   []int    `json:"member_ids"`
	Member      []Member `json:"members"`
}

type Member struct {
	ConversationID         *int       `json:"conversation_id,omitempty"`
	UserID                 int        `json:"user_id"`
	IsAdmin                bool       `json:"is_admin"`
	Nickname               *string    `json:"nickname"`
	Privileges             *[]string  `json:"privileges,omitempty"`
	LastReadMessageID      *int       `json:"last_read_message_id,omitempty"`
	UnreadCount            *int       `json:"unread_count,omitempty"`
	NotificationPreference *string    `json:"notification_preference,omitempty"`
	MutedUntil             *time.Time `json:"muted_until,omitempty"`
	JoinedAt               *time.Time `json:"joined_at"`
	LeftAt                 *time.Time `json:"left_at,omitempty"`

	Username    *string `json:"username"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	CompanyName *string `json:"company_name"`
	ImageURL    *string `json:"image_url"`
}

func ExtractMemberIDs(members []Member) []int {
	memberIDs := make([]int, len(members))
	for i, member := range members {
		memberIDs[i] = member.UserID
	}
	return memberIDs
}

type UsersJSONB struct {
	UserID int    `json:"user_id"`
	DT     string `json:"dt"`
}

func prepareUsersJSONB(memberIDs []int, timestamp time.Time) (string, error) {
	formattedTime := timestamp.Format("2006-01-02 15:04:05")

	deletedForList := make([]UsersJSONB, len(memberIDs))
	for i, id := range memberIDs {
		deletedForList[i] = UsersJSONB{
			UserID: id,
			DT:     formattedTime,
		}
	}

	jsonData, err := json.Marshal(deletedForList)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type CallRoom struct {
	ID             int       `json:"id"`
	UUID           string    `json:"uuid"`
	ConversationID int       `json:"conversation_id"`
	MaxUser        int       `json:"max_user"`
	UserIDs        []string  `json:"user_ids"`
	ProfileIDs     []string  `json:"profile_ids"`
	Title          string    `json:"title"`
	Hex            string    `json:"hex"`
	Duration       int       `json:"duration"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Active         int       `json:"active"`
	Deleted        int       `json:"deleted"`
	JoinURL        string    `json:"join_url,omitempty"`
	JitsiURL       string    `json:"jitsi_url,omitempty"`
}

type CreateCallRoomRequest struct {
	UserIDs        []string `json:"user_ids" binding:"required"`
	ConversationID int      `json:"conversation_id"`
	ProfileIDs     []string `json:"profile_ids"`
	Title          string   `json:"title"`
	Duration       int      `json:"duration"`
	MaxUser        int      `json:"max_user"`
}
