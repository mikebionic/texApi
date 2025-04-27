package chat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/fileUtils"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

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

func (r *Repository) SaveMessage(msg *Message) (int, error) {
	ctx := context.Background()

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

	// TODO: Can this be optimized with RabbitMQ?
	_, err = r.db.Exec(ctx, `
		UPDATE tbl_conversation 
		SET last_message_id = $1, message_count = message_count + 1, last_activity = CURRENT_TIMESTAMP 
		WHERE id = $2
	`, messageID, msg.ConversationID)

	if err != nil {
		log.Printf("Error updating conversation stats: %v", err)
	}

	// TODO: Can this be optimized with RabbitMQ?
	// Update unread counts for all members except sender
	_, err = r.db.Exec(ctx, `
		UPDATE tbl_conversation_member
		SET unread_count = unread_count + 1
		WHERE conversation_id = $1 AND user_id != $2 AND active = 1 AND deleted = 0
	`, msg.ConversationID, msg.SenderID)
	if err != nil {
		log.Printf("Error updating unread counts: %v", err)
	}

	//// Create notifications for members (except sender)
	//// TODO: DO WE REALLY NEED THIS???
	//_, err = r.db.Exec(ctx, `
	//	INSERT INTO tbl_notification (user_id, conversation_id, message_id, notification_type, content)
	//	SELECT cm.user_id, $1, $2, $3, $4
	//	FROM tbl_conversation_member cm
	//	WHERE cm.conversation_id = $1 AND cm.user_id != $5 AND cm.active = 1 AND cm.deleted = 0
	//`, msg.ConversationID, messageID, "new_message", msg.Content, msg.SenderID)
	//
	//if err != nil {
	//	log.Printf("Error creating notifications: %v", err)
	//}

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

	// TODO: Can this be optimized with RabbitMQ?
	_, err = r.db.Exec(context.Background(), `
		UPDATE tbl_conversation 
		SET last_message_id = $1, message_count = message_count + 1, last_activity = CURRENT_TIMESTAMP 
		WHERE id = $2
	`, messageID, msg.ConversationID)

	if err != nil {
		log.Printf("Error updating conversation stats: %v", err)
	}

	// TODO: Can this be optimized with RabbitMQ?
	// Update unread counts for all members except sender
	_, err = r.db.Exec(context.Background(), `
		UPDATE tbl_conversation_member
		SET unread_count = unread_count + 1
		WHERE conversation_id = $1 AND user_id != $2 AND active = 1 AND deleted = 0
	`, msg.ConversationID, msg.SenderID)
	if err != nil {
		log.Printf("Error updating unread counts: %v", err)
	}

	//// TODO: tbl_notification not used, Also optional
	//// Create notifications for members (except sender)
	//_, err = tx.Exec(context.Background(), `
	//	INSERT INTO tbl_notification (user_id, conversation_id, message_id, notification_type, content)
	//	SELECT cm.user_id, $1, $2, $3, $4
	//	FROM tbl_conversation_member cm
	//	WHERE cm.conversation_id = $1 AND cm.user_id != $5 AND cm.active = 1 AND cm.deleted = 0
	//`, msg.ConversationID, messageID, "new_message", msg.Content, msg.SenderID)
	//
	//if err != nil {
	//	log.Printf("Error creating notifications: %v", err)
	//}

	return messageID, nil
}

// TODO: Probably not required route.
func (r *Repository) GetConversation(conversationID int) (*Conversation, error) {
	ctx := context.Background()
	query := `
		SELECT *
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

func (r *Repository) GetConversations(userID int) (*[]Conversation, error) {
	query := `
		SELECT 
		    c.*,
		   cm.unread_count,
		   (SELECT content FROM tbl_message WHERE id = c.last_message_id) as last_message
		FROM tbl_conversation c
		JOIN tbl_conversation_member cm ON c.id = cm.conversation_id
		WHERE cm.user_id = $1 AND cm.deleted = 0 AND c.deleted = 0
		ORDER BY c.last_activity DESC
	`

	var conversations []Conversation
	err := pgxscan.Select(context.Background(), r.db, &conversations, query, userID)
	return &conversations, err
}

func (r *Repository) GetConversationMembers(conversationID int) ([]Member, error) {
	query := `
		SELECT cm.user_id, cm.is_admin, cm.nickname, cm.privileges,
		       cm.last_read_message_id, cm.joined_at, cm.unread_count,
		       cm.notification_preference, cm.muted_until,
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

// NOTE!: You should send messageIDs if you want to remove messages separately.
// Else for deletion of whole conversation we put nil (all messages)
func (r *Repository) RemoveMessages(conversationID int, memberIDs, messageIDs []int, forEveryone bool) (err error) {
	if conversationID < 1 {
		return fmt.Errorf("invalid conversation ID")
	}
	if messageIDs == nil && memberIDs == nil {
		return fmt.Errorf("cannot remove messages without a message ID and member ID")
	}

	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	now := time.Now()
	jsonbData, err := prepareUsersJSONB(memberIDs, now)
	if err != nil {
		return fmt.Errorf("failed to prepare JSONB data: %w", err)
	}

	query := `
        UPDATE tbl_message 
        SET deleted_for = COALESCE(deleted_for, '[]'::jsonb) || $1::jsonb, 
            updated_at = CURRENT_TIMESTAMP
    `

	if forEveryone {
		query += `, deleted = 1, active = 0`
	}

	var commandTag pgconn.CommandTag

	if messageIDs == nil {
		query += ` WHERE conversation_id = $2`
		commandTag, err = tx.Exec(ctx, query, jsonbData, conversationID)
	} else {
		query += ` WHERE conversation_id = $2 AND id = ANY($3)`
		commandTag, err = tx.Exec(ctx, query, jsonbData, conversationID, messageIDs)
	}

	if err != nil {
		return fmt.Errorf("database execution error: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		fmt.Println("Warning: No rows were updated")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Function for admin and message sender_id
func (r *Repository) DeleteMessageOfOwner(messageID, userID int, isAdmin bool) error {
	ctx := context.Background()
	if !isAdmin {
		var senderID int
		err := r.db.QueryRow(ctx, `
            SELECT sender_id FROM tbl_message 
            WHERE id = $1 AND active = 1 AND deleted = 0
        `, messageID).Scan(&senderID)

		if err != nil {
			return err
		}

		if senderID != userID {
			return errors.New("only the sender or an admin can delete this message")
		}
	}

	_, err := r.db.Exec(ctx, `
        UPDATE tbl_message 
        SET deleted = 1, active = 0, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `, messageID)

	return err
}

func (r *Repository) AddConversationMembers(conversationID int, memberIDs []int) error {
	if conversationID < 1 {
		return fmt.Errorf("invalid conversation ID")
	}
	if len(memberIDs) == 0 {
		return fmt.Errorf("empty member IDs")
	}

	query := `
		INSERT INTO tbl_conversation_member 
			(conversation_id, user_id)
		VALUES ($1, $2)
	`

	errList := make([]string, 0)
	for _, memberID := range memberIDs {
		_, err := r.db.Exec(context.Background(), query, conversationID, memberID)
		if err != nil {
			errList = append(errList, err.Error())
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("failed to add conversation members: %w", errors.New(strings.Join(errList, ", ")))
	}

	query = `UPDATE tbl_conversation SET 
		member_count = member_count + $2,
		last_activity = CURRENT_TIMESTAMP 
        WHERE id = $1`
	commandTag, err := r.db.Exec(context.Background(), query, conversationID, len(memberIDs)-len(errList))
	if err != nil {
		return fmt.Errorf("failed to update member count of conversation: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows were updated while updating member_count")
	}

	return nil
}

func (r *Repository) RemoveConversationMembers(conversationID int, memberIDs []int) error {
	if conversationID < 1 {
		return fmt.Errorf("invalid conversation ID")
	}
	if len(memberIDs) == 0 {
		return fmt.Errorf("empty member IDs")
	}

	query := `
		UPDATE tbl_conversation_member 
		SET left_at = CURRENT_TIMESTAMP, 
		    active = 0,
		    deleted = 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE conversation_id = $1 
		AND user_id = ANY($2)
	`

	commandTag, err := r.db.Exec(context.Background(), query, conversationID, memberIDs)
	if err != nil {
		return fmt.Errorf("failed to remove conversation members: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows were updated")
	}

	query = `UPDATE tbl_conversation SET 
            member_count = member_count - $2,
            last_activity = CURRENT_TIMESTAMP 
            WHERE id = $1`
	commandTag, err = r.db.Exec(context.Background(), query, conversationID, len(memberIDs))
	if err != nil {
		return fmt.Errorf("failed to update member count of conversation: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows were updated")
	}

	return nil
}

func (r *Repository) RemoveConversation(conversationID int) error {
	query := `
		UPDATE tbl_conversation 
		SET deleted = 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted = 0
		RETURNING id
	`
	var deletedID int
	err := db.DB.QueryRow(context.Background(), query, conversationID).Scan(&deletedID)
	return err
}

func (r *Repository) GetConversationMessages(conversationID, userID, limit, offset int) ([]MessageDetails, error) {
	ctx := context.Background()
	query := `
		SELECT m.*,
               TRIM(COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '')) AS sender_name,
			   p.image_url as sender_avatar
		FROM tbl_message m
		JOIN tbl_user u ON m.sender_id = u.id
		JOIN tbl_company p ON u.company_id = p.id
		WHERE m.conversation_id = $1 AND m.active = 1 AND m.deleted = 0
		AND NOT EXISTS (
			SELECT 1
			FROM jsonb_array_elements(deleted_for) AS elem
			WHERE (elem->>'user_id')::INT = $4
		)
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var messages []MessageDetails
	err := pgxscan.Select(ctx, r.db, &messages, query, conversationID, limit, offset, userID)
	if err != nil {
		return nil, err
	}

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

func (r *Repository) UpdateConversation(conversationID, creatorID int, conv Conversation) error {
	query := `
        UPDATE tbl_conversation SET
            chat_type = COALESCE($1, chat_type), 
            title = COALESCE($2, title), 
            description = COALESCE($3, description), 
            is_public = COALESCE($4, is_public),  
            public_url = COALESCE($5, public_url),
            theme_color = COALESCE($6, theme_color), 
            image_url = COALESCE($7, image_url), 
            background_image_url = COALESCE($8, background_image_url), 
            background_blur = COALESCE($9, background_blur),
            auto_delete_duration = COALESCE($10, auto_delete_duration), 
            invite_token = COALESCE($11, invite_token),
            updated_at = CURRENT_TIMESTAMP, last_activity = CURRENT_TIMESTAMP
        WHERE creator_id = $12 AND id = $13
        `
	tag, err := r.db.Exec(context.Background(), query,
		conv.ChatType,
		conv.Title,
		conv.Description,
		conv.IsPublic,
		conv.PublicURL,
		conv.ThemeColor,
		conv.ImageURL,
		conv.BackgroundImageURL,
		conv.BackgroundBlur,
		conv.AutoDeleteDuration,
		conv.InviteToken,
		creatorID,
		conversationID)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("no rows were updated")
	}
	return nil
}

func (r *Repository) UpdateConversationMember(conversationID int, req UpdateMemberRequest) error {
	baseQuery := `UPDATE tbl_conversation_member SET `

	var setClauses []string
	var params []interface{}

	params = append(params, conversationID, req.UserID)
	paramCount := 3

	if req.Nickname != nil {
		setClauses = append(setClauses, fmt.Sprintf("nickname = $%d", paramCount))
		params = append(params, req.Nickname)
		paramCount++
	}

	if req.IsAdmin != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_admin = $%d", paramCount))
		params = append(params, req.IsAdmin)
		paramCount++
	}

	if req.Privileges != nil {
		setClauses = append(setClauses, fmt.Sprintf("privileges = $%d", paramCount))
		params = append(params, req.Privileges)
		paramCount++
	}

	if req.NotificationPref != nil {
		setClauses = append(setClauses, fmt.Sprintf("notification_preference = $%d", paramCount))
		params = append(params, req.NotificationPref)
		paramCount++
	}

	if req.MutedUntil != nil {
		setClauses = append(setClauses, fmt.Sprintf("muted_until = $%d", paramCount))
		params = append(params, req.MutedUntil)
		paramCount++
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := baseQuery + strings.Join(setClauses, ", ") +
		" WHERE conversation_id = $1 AND user_id = $2"

	commandTag, err := db.DB.Exec(context.Background(), query, params...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no rows were updated")
	}

	return nil
}

func (r *Repository) CreateConversation(creatorID int, conv CreateConversation) (int, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var conversationID int
	err = tx.QueryRow(ctx, `
            INSERT INTO tbl_conversation (
                    chat_type, title, description, creator_id,  is_public, member_count,
                	theme_color, background_image_url, background_blur
            ) VALUES (
                    $1, $2, $3, $4, $5, $6, $7, $8, $9
            ) RETURNING id
    `, conv.ChatType, conv.Title, conv.Description, creatorID, conv.IsPublic, len(conv.Members)+1,
		conv.ThemeColor, conv.BackgroundImageURL, conv.BackgroundBlur).Scan(&conversationID)

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

	// TODO WARNING: THIS SHOULD check that members are not blocked or blocking user. notify accordingly
	// Query the tbl_contact and check creatorID as a blocking or blocked user
	// Query and avoid adding not existing users, companys
	for _, memberID := range conv.Members {
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

			////// TODO: Do we need this? message_id should be nullable then
			//// Create notification for this member
			//_, err = tx.Exec(ctx, `
			//    INSERT INTO tbl_notification (user_id, conversation_id, notification_type, content)
			//    VALUES ($1, $2, $3, $4)
			//`, memberID, conversationID, "new_conversation", "You were added to conversation "+title)
			//
			//if err != nil {
			//	log.Printf("Error creating invitation notification: %v", err)
			//}
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
