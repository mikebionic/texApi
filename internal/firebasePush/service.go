package firebasePush

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"texApi/config"
	"texApi/pkg/utils"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var firebaseClient *messaging.Client

func InitFirebase() error {
	opt := option.WithCredentialsFile(config.ENV.FIREBASE_ADMINSDK_FILE)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error getting messaging client: %v", err)
	}

	firebaseClient = client
	log.Println("Firebase initialized successfully")
	return nil
}

func SaveNotificationToken(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)

	var req SaveNotificationTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request", err.Error()))
		return
	}

	if err := SaveToken(userID, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Notification token saved successfully", req))
}

func SendNotificationToUser(userID int, payload NotificationPayload) error {
	if firebaseClient == nil {
		return fmt.Errorf("firebase client not initialized")
	}

	tokens, err := GetUserTokens(userID)
	if err != nil {
		return fmt.Errorf("error getting user tokens: %v", err)
	}

	if len(tokens) == 0 {
		log.Printf("No tokens found for user %d", userID)
		return nil
	}

	data := map[string]string{
		"sender_name":     payload.SenderName,
		"conversation_id": strconv.Itoa(payload.ConversationID),
		"user_id":         strconv.Itoa(payload.UserID),
		"content":         payload.Content,
		"type":            payload.Type,
		"app_name":        config.ENV.APP_NAME,
		"is_silent":       strconv.Itoa(payload.IsSilent),
		"created_at":      payload.CreatedAt,
	}

	defaultMessage := fmt.Sprintf("New message from %s", payload.SenderName)
	if payload.Title == nil {
		payload.Title = &defaultMessage
	}
	notification := &messaging.Notification{
		Title: *payload.Title,
		Body:  payload.Content,
	}

	if len(tokens) == 1 {
		message := &messaging.Message{
			Notification: notification,
			Data:         data,
			Token:        tokens[0],
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					ClickAction: "FLUTTER_NOTIFICATION_CLICK",
					ChannelID:   "messages",
				},
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Alert: &messaging.ApsAlert{
							Title: notification.Title,
							Body:  notification.Body,
						},
						Badge: func(i int) *int { return &i }(1),
						Sound: "default",
					},
				},
			},
		}

		response, err := firebaseClient.Send(context.Background(), message)
		if err != nil {
			log.Printf("Error sending single notification: %v", err)
			if isInvalidToken(err) {
				go removeInvalidToken(tokens[0])
			}
			return err
		}
		log.Printf("Successfully sent message to user %d: %s", userID, response)
	} else {
		multicastMessage := &messaging.MulticastMessage{
			Notification: notification,
			Data:         data,
			Tokens:       tokens,
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					ClickAction: "FLUTTER_NOTIFICATION_CLICK",
					ChannelID:   "messages",
				},
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Alert: &messaging.ApsAlert{
							Title: notification.Title,
							Body:  notification.Body,
						},
						Badge: func(i int) *int { return &i }(1),
						Sound: "default",
					},
				},
			},
		}

		br, err := firebaseClient.SendEachForMulticast(context.Background(), multicastMessage)
		if err != nil {
			log.Printf("Error sending multicast notification: %v", err)
			return err
		}

		log.Printf("Successfully sent %d messages to user %d", br.SuccessCount, userID)

		if br.FailureCount > 0 {
			var failedTokens []string
			for idx, resp := range br.Responses {
				if !resp.Success {
					failedTokens = append(failedTokens, tokens[idx])
					log.Printf("Failed to send to token: %s, error: %v", tokens[idx], resp.Error)

					if isInvalidToken(resp.Error) {
						go removeInvalidToken(tokens[idx])
					}
				}
			}
			log.Printf("List of tokens that caused failures: %v", failedTokens)
		}
	}

	return nil
}

func SendNotificationToConversation(conversationID, senderID int, senderName, content string) error {
	participants, err := GetConversationParticipants(conversationID, senderID)
	if err != nil {
		return fmt.Errorf("error getting conversation participants: %v", err)
	}

	payload := NotificationPayload{
		SenderName:     senderName,
		ConversationID: conversationID,
		Content:        content,
	}

	for _, participantID := range participants {
		go func(userID int) {
			if err := SendNotificationToUser(userID, payload); err != nil {
				log.Printf("Error sending notification to user %d: %v", userID, err)
			}
		}(participantID)
	}

	return nil
}
