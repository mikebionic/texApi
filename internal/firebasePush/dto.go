package firebasePush

import "time"

type FirebaseToken struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Token      string    `json:"token"`
	DeviceType *string   `json:"device_type"`
	Meta       *string   `json:"meta"`
	Meta2      *string   `json:"meta2"`
	Meta3      *string   `json:"meta3"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Active     int       `json:"active"`
	Deleted    int       `json:"deleted"`
}

type SaveNotificationTokenRequest struct {
	NotificationToken string  `json:"notification_token" binding:"required"`
	DeviceType        *string `json:"device_type,omitempty"`
}

type NotificationPayload struct {
	SenderName     string  `json:"sender_name"`
	ConversationID int     `json:"conversation_id"`
	UserID         int     `json:"user_id"`
	Content        string  `json:"content"`
	Title          *string `json:"title"`
	CreatedAt      string  `json:"created_at"`
	Type           string  `json:"type"`
	IsSilent       int     `json:"is_silent"`
}
