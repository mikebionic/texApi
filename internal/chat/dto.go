package chat

import (
	"texApi/internal/dto"
	"time"
)

type Message struct {
	Content        string                  `json:"content"`
	ConversationID int                     `json:"conversation_id"` // Kind of ROOM
	SenderID       int                     `json:"sender_id"`       //UserID
	SenderName     *string                 `json:"sender_name,omitempty"`
	MessageType    string                  `json:"message_type"`
	ReplyToID      *int                    `json:"reply_to_id,omitempty"`
	ForwardedFrom  *int                    `json:"forwarded_from_id,omitempty"`
	MediaID        *int                    `json:"media_id,omitempty"`
	StickerID      *int                    `json:"sticker_id,omitempty"`
	IsSilent       *bool                   `json:"is_silent,omitempty"`
	CreatedAt      *time.Time              `json:"created_at,omitempty"`
	Extras         *map[string]interface{} `json:"extras,omitempty"`        // For additional metadata
	TypingStatus   bool                    `json:"typing_status,omitempty"` // true when typing, false when stopped
	OnlineStatus   *OnlineStatus           `json:"online_status,omitempty"`
	Media          *[]dto.MediaMain        `json:"media,omitempty"`
}

type OnlineStatus struct {
	UserID   int  `json:"user_id"`
	IsOnline bool `json:"is_online"`
}

type MessageDetails struct {
	ID                int              `json:"id"`
	UUID              string           `json:"uuid"`
	ConversationID    int              `json:"conversation_id"`
	ConversationTitle *string          `json:"conversation_title,omitempty"`
	ConversationType  *string          `json:"conversation_type,omitempty"`
	SenderID          int              `json:"sender_id"`
	MessageType       string           `json:"message_type"`
	Content           string           `json:"content"`
	ReplyToID         *int             `json:"reply_to_id,omitempty"`
	ForwardedFromID   *int             `json:"forwarded_from_id,omitempty"`
	MediaID           *int             `json:"media_id,omitempty"`
	StickerID         *int             `json:"sticker_id,omitempty"`
	IsEdited          *bool            `json:"is_edited"`
	IsPinned          *bool            `json:"is_pinned"`
	IsSilent          *bool            `json:"is_silent,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
	SenderName        string           `json:"sender_name"`
	SenderAvatar      string           `json:"sender_avatar"`
	Reactions         *[]Reaction      `json:"reactions,omitempty"`
	Media             *[]dto.MediaMain `json:"media,omitempty"` // Added for media attachments
}

type Reaction struct {
	MessageID  int     `json:"message_id"`
	UserID     int     `json:"user_id"`
	CompanyID  int     `json:"company_id"`
	Emoji      string  `json:"emoji"`
	SenderName *string `json:"sender_name,omitempty"`
}

type Conversation struct {
	ID                 int    `json:"id"`
	UUID               string `json:"uuid"`
	ChatType           string `json:"chat_type"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	CreatorID          int    `json:"creator_id"`
	ThemeColor         string `json:"theme_color"`
	ImageURL           string `json:"image_url"`
	BackgroundImageURL string `json:"background_image_url"`
	BackgroundBlur     int    `json:"background_blur"`
	MemberCount        int    `json:"member_count"`
	MessageCount       int    `json:"message_count"`
	AutoDeleteDuration int    `json:"auto_delete_duration"`
}

type Member struct {
	UserID      int        `json:"user_id"`
	IsAdmin     bool       `json:"is_admin"`
	Nickname    *string    `json:"nickname"`
	JoinedAt    *time.Time `json:"joined_at"`
	Username    *string    `json:"username"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	CompanyName *string    `json:"company_name"`
	ImageURL    *string    `json:"image_url"`
}
