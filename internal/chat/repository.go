package chat

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
	"texApi/internal/dto"
	"texApi/pkg/fileUtils"
)

type Repository struct {
	db *pgxpool.Pool // Make sure this matches your database.DB type
}

// Update constructor if needed
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// GetUserConversations returns all conversations a user is a member of
func (r *Repository) GetUserConversations(userID int) ([]int, error) {
	ctx := context.Background()
	query := `
		SELECT conversation_id 
		FROM tbl_conversation_member 
		WHERE user_id = $1 AND active = 1 AND deleted = 0
	`

	var conversationIDs []int
	err := pgxscan.Select(ctx, r.db, &conversationIDs, query, userID)
	if err != nil {
		return nil, err
	}

	return conversationIDs, nil
}

// CanAccessConversation checks if a user can access a specific conversation
func (r *Repository) CanAccessConversation(userID, conversationID int) bool {
	ctx := context.Background()
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM tbl_conversation_member 
			WHERE user_id = $1 AND conversation_id = $2 AND active = 1 AND deleted = 0
		)
	`

	var canAccess bool
	err := r.db.QueryRow(ctx, query, userID, conversationID).Scan(&canAccess)
	if err != nil {
		log.Printf("Error checking conversation access: %v", err)
		return false
	}

	return canAccess
}

// SaveMessage saves a message to the database and returns the message ID
func (r *Repository) SaveMessage(msg *Message) (int, error) {
	ctx := context.Background()

	//// Don't allow empty messages
	////&& msg.MediaID == 0 && msg.StickerID == 0 {
	//if msg.Content == "" {
	//	return 0, errors.New("message content cannot be empty")
	//}

	query := `
		INSERT INTO tbl_message (
			conversation_id, sender_id, message_type, content, 
			reply_to_id, forwarded_from_id, media_id, sticker_id, is_silent
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`

	replyToID := msg.ReplyToID
	forwardedFrom := msg.ForwardedFrom

	var messageID int
	err := r.db.QueryRow(
		ctx, query,
		msg.ConversationID, msg.SenderID, msg.MessageType, msg.Content,
		replyToID, forwardedFrom, msg.MediaID, msg.StickerID, msg.IsSilent,
	).Scan(&messageID)

	if err != nil {
		return 0, err
	}

	// Update conversation last activity and message count
	// TODO: DO WE REALLY NEED THIS???
	_, err = r.db.Exec(ctx, `
		UPDATE tbl_conversation 
		SET last_message_id = $1, message_count = message_count + 1, last_activity = CURRENT_TIMESTAMP 
		WHERE id = $2
	`, messageID, msg.ConversationID)

	if err != nil {
		log.Printf("Error updating conversation stats: %v", err)
	}
	return messageID, nil
}

func (r *Repository) SaveMessageTx(tx pgx.Tx, msg *Message) (int, error) {
	query := `
		INSERT INTO tbl_message (
			conversation_id, sender_id, message_type, content, 
			reply_to_id, forwarded_from_id, media_id, sticker_id, is_silent
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`

	replyToID := msg.ReplyToID
	forwardedFrom := msg.ForwardedFrom

	var messageID int
	err := tx.QueryRow(
		context.Background(), query,
		msg.ConversationID, msg.SenderID, msg.MessageType, msg.Content,
		replyToID, forwardedFrom, msg.MediaID, msg.StickerID, msg.IsSilent,
	).Scan(&messageID)
	if err != nil {
		return 0, err
	}

	// TODO: NOTE! OPTIONAL
	_, err = tx.Exec(context.Background(), `
		UPDATE tbl_conversation 
		SET last_message_id = $1, message_count = message_count + 1, last_activity = CURRENT_TIMESTAMP 
		WHERE id = $2
	`, messageID, msg.ConversationID)
	if err != nil {
		log.Printf("Error updating conversation stats: %v", err)
	}
	return messageID, nil
}

func (r *Repository) GetConversation(conversationID int) (*Conversation, error) {
	ctx := context.Background()
	query := `
		SELECT id, uuid, chat_type, title, description, creator_id, 
			   theme_color, image_url, background_image_url, background_blur,
			   member_count, message_count, auto_delete_duration
		FROM tbl_conversation
		WHERE id = $1 AND active = 1 AND deleted = 0
	`

	var conversation Conversation
	err := pgxscan.Get(ctx, r.db, &conversation, query, conversationID)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *Repository) GetConversationMembers(conversationID int) ([]Member, error) {
	query := `
		SELECT cm.user_id, cm.is_admin, cm.nickname, cm.joined_at,
		       u.username, p.first_name, p.last_name, p.company_name, p.image_url
		FROM tbl_conversation_member cm
		JOIN tbl_user u ON cm.user_id = u.id
		JOIN tbl_company p ON u.company_id = p.id
		WHERE cm.conversation_id = $1 AND cm.active = 1 AND cm.deleted = 0
	`

	var members []Member
	err := pgxscan.Select(context.Background(), r.db, &members, query, conversationID)
	return members, err
}

func (r *Repository) GetConversationMessages(conversationID, userID, limit, offset int) ([]MessageDetails, error) {
	ctx := context.Background()
	query := `
		SELECT m.id, m.uuid, m.conversation_id, m.sender_id, m.message_type, 
			   m.content, m.reply_to_id, m.forwarded_from_id, m.media_id, 
			   m.sticker_id, m.is_edited, m.is_pinned, m.is_silent, m.created_at,
               TRIM(COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '')) AS sender_name,
			   p.image_url as sender_avatar
		FROM tbl_message m
		JOIN tbl_user u ON m.sender_id = u.id
		JOIN tbl_company p ON u.company_id = p.id
		WHERE m.conversation_id = $1 AND m.active = 1 AND m.deleted = 0
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var messages []MessageDetails
	err := pgxscan.Select(ctx, r.db, &messages, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	// TODO: this one always updates last_read_message_id???
	_, err = r.db.Exec(ctx, `
		UPDATE tbl_conversation_member 
		SET last_read_message_id = (
			SELECT MAX(id) FROM tbl_message 
			WHERE conversation_id = $1 AND active = 1 AND deleted = 0
		),
		unread_count = 0
		WHERE conversation_id = $1 AND user_id = $2 AND active = 1 AND deleted = 0
	`, conversationID, userID)
	if err != nil {
		log.Printf("Error marking messages as read: %v", err)
	}

	// TODO: check this
	// If we have messages, fetch the related media for each message
	if len(messages) > 0 {
		var messageIDs []interface{}
		messageIDToIndex := make(map[int]int)

		for i, msg := range messages {
			messageIDs = append(messageIDs, msg.ID)
			messageIDToIndex[msg.ID] = i
		}

		placeholders := make([]string, len(messageIDs))
		for i := range messageIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}

		var allMedia []struct {
			dto.MediaMain
			MessageID int `db:"message_id"`
		}

		mediaQueryWithJoin := fmt.Sprintf(`
            SELECT m.*, mm.message_id 
            FROM tbl_media m
            JOIN tbl_message_media mm ON m.id = mm.media_id
            WHERE mm.message_id IN (%s) AND m.deleted = 0
            ORDER BY mm.message_id, mm.sort_order
        `, strings.Join(placeholders, ","))

		err = pgxscan.Select(ctx, r.db, &allMedia, mediaQueryWithJoin, messageIDs...)
		if err != nil {
			log.Printf("Error fetching media for messages: %v", err)
			// Continue without media if there's an error
		} else {
			for _, media := range allMedia {
				generatedURL := fileUtils.GenerateMediaURL(media.UUID, media.Filename)
				media.URL = generatedURL["url"]
				media.ThumbURL = generatedURL["thumb_url"]
				if idx, ok := messageIDToIndex[media.MessageID]; ok {

					if messages[idx].Media == nil {
						messages[idx].Media = &[]dto.MediaMain{}
					}
					*messages[idx].Media = append(*messages[idx].Media, media.MediaMain)
				}
			}
		}
	}
	return messages, nil
}

func (r *Repository) CreateConversation(creatorID int, title, description string, chatType string, members []int) (int, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var conversationID int
	err = tx.QueryRow(ctx, `
            INSERT INTO tbl_conversation (
                    chat_type, title, description, creator_id, member_count
            ) VALUES (
                    $1, $2, $3, $4, $5
            ) RETURNING id
    `, chatType, title, description, creatorID, len(members)+1).Scan(&conversationID)

	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, `
            INSERT INTO tbl_conversation_member (
                    conversation_id, user_id, is_admin
            ) VALUES (
                    $1, $2, true
            )
    `, conversationID, creatorID)

	if err != nil {
		return 0, err
	}

	//creatorName, err := GetCreatorName(creatorID)

	for _, memberID := range members {
		if memberID != creatorID {
			_, err = tx.Exec(ctx, `
				INSERT INTO tbl_conversation_member (
					conversation_id, user_id, is_admin
				) VALUES (
					$1, $2, false
				)
			`, conversationID, memberID)

			if err != nil {
				return 0, err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return conversationID, nil
}

func (r *Repository) GetCreatorName(userID int) (string, error) {
	var creatorName string
	query := `
	  SELECT TRIM(COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, ''))
	  FROM tbl_user u
	  JOIN tbl_company p ON u.company_id = p.id
	  WHERE u.id = $1
	`
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&creatorName)
	if err != nil {
		log.Printf("Error getting creator name: %v", err)
		creatorName = "Someone"
	}
	return creatorName, nil
}
