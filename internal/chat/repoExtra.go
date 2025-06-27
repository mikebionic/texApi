package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *Repository) SearchMessages(userID int, searchQuery string, limit, offset int) ([]MessageDetails, error) {
	ctx := context.Background()
	query := `
        SELECT m.id, m.uuid, m.conversation_id, m.sender_id, m.message_type, 
               m.content, m.reply_to_id, m.forwarded_from_id, m.media_id, 
               m.sticker_id, m.is_edited, m.is_pinned, m.is_silent, m.created_at,
				TRIM(
					COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '') || 
					COALESCE(d.first_name,'') || ' ' || COALESCE(d.last_name,'') 
				) AS sender_name,
				COALESCE(p.image_url, d.image_url) AS sender_avatar,
			   c.title as conversation_title,
               c.chat_type as conversation_type
		FROM tbl_message m
		JOIN tbl_user u ON m.sender_id = u.id
		LEFT JOIN tbl_company p ON u.company_id = p.id
		LEFT JOIN tbl_driver d ON u.driver_id = d.id

        JOIN tbl_conversation c ON m.conversation_id = c.id
        JOIN tbl_conversation_member cm ON (c.id = cm.conversation_id AND cm.user_id = $1)
        WHERE 
          NOT EXISTS (
				SELECT 1
				FROM jsonb_array_elements(m.deleted_for) AS elem
				WHERE (elem->>'user_id')::INT = $4
			) AND
            m.deleted = 0 AND
            cm.deleted = 0 AND
            c.deleted = 0 AND
            (m.content ILIKE $2 OR c.title ILIKE $2 or 
                TRIM(COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '')) ILike $2)
        ORDER BY m.created_at DESC
        LIMIT $3 OFFSET $4
    `

	var messages []MessageDetails
	err := pgxscan.Select(ctx, r.db, &messages, query, userID, "%"+searchQuery+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *Repository) PinMessage(messageID, conversationID int, isPinned bool) error {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
            UPDATE tbl_message SET is_pinned = $3
            WHERE id = $1 AND conversation_id = $2
        `, messageID, conversationID, isPinned)

	return tx.Commit(ctx)
}

func (r *Repository) EditMessage(messageID, userID int, newContent string) error {
	ctx := context.Background()

	var senderID int
	err := r.db.QueryRow(ctx, `
        SELECT sender_id FROM tbl_message 
        WHERE id = $1 AND active = 1 AND deleted = 0
    `, messageID).Scan(&senderID)

	if err != nil {
		return err
	}
	if senderID != userID {
		return errors.New("only the sender can edit this message")
	}

	_, err = r.db.Exec(ctx, `
        UPDATE tbl_message 
        SET content = $1, is_edited = true, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2 AND sender_id = $3
    `, newContent, messageID, userID)

	return err
}

func (r *Repository) SetMessageRead(messageID, userID, conversationID int) error {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
        UPDATE tbl_conversation_member 
        SET last_read_message_id = $1, unread_count = 0,
        updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $2 AND conversation_id = $3
    `, messageID, userID, conversationID)

	if err != nil {
		return fmt.Errorf("failed to update conversation member: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no conversation member found to update")
	}

	_, err = tx.Exec(ctx, `
		UPDATE tbl_message 
		SET read_by = COALESCE(read_by, '[]'::jsonb) || $1::jsonb,
		read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND sender_id != $3
	`, []int{userID}, messageID, userID)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) AddReaction(messageID, userID, companyID int, emoji string) error {
	ctx := context.Background()

	var exists bool
	err := r.db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM tbl_message_reaction
            WHERE message_id = $1 AND company_id = $2 AND emoji = $3
        )
    `, messageID, companyID, emoji).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		// Remove reaction if it already exists (toggle behavior)
		_, err = r.db.Exec(ctx, `
            DELETE FROM tbl_message_reaction
            WHERE message_id = $1 AND company_id = $2 AND emoji = $3
        `, messageID, companyID, emoji)
	} else {
		// Add new reaction
		_, err = r.db.Exec(ctx, `
            INSERT INTO tbl_message_reaction (message_id, user_id, company_id, emoji)
            VALUES ($1, $2, $3, $4)
        `, messageID, userID, companyID, emoji)
	}

	return err
}

func (r *Repository) GetMessageDetails(messageID int) (*MessageDetails, error) {
	ctx := context.Background()
	query := `
        SELECT m.id, m.uuid, m.conversation_id, m.sender_id, m.message_type, 
               m.content, m.reply_to_id, m.forwarded_from_id, m.media_id, 
               m.sticker_id, m.is_edited, m.is_pinned, m.is_silent, m.created_at,
		TRIM(
				COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '') || 
				COALESCE(d.first_name,'') || ' ' || COALESCE(d.last_name,'') 
			) AS sender_name,
			COALESCE(p.image_url, d.image_url) AS sender_avatar
		FROM tbl_message m
		JOIN tbl_user u ON m.sender_id = u.id
		LEFT JOIN tbl_company p ON u.company_id = p.id
		LEFT JOIN tbl_driver d ON u.driver_id = d.id

        WHERE m.id = $1 AND m.active = 1
    `

	var message MessageDetails
	err := pgxscan.Get(ctx, r.db, &message, query, messageID)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *Repository) IsConversationAdmin(userID, conversationID int) (bool, error) {
	ctx := context.Background()
	var isAdmin bool

	err := r.db.QueryRow(ctx, `
        SELECT is_admin FROM tbl_conversation_member 
        WHERE user_id = $1 AND conversation_id = $2 AND active = 1 AND deleted = 0
    `, userID, conversationID).Scan(&isAdmin)

	if err != nil {
		return false, err
	}

	return isAdmin, nil
}

func (r *Repository) GetMessageReactions(messageIDs []int) ([]Reaction, error) {
	ctx := context.Background()
	query := `
        SELECT mr.message_id, mr.company_id, mr.emoji, mr.user_id,
                TRIM(COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '')) AS sender_name
        FROM tbl_message_reaction mr
        LEFT JOIN tbl_company p ON mr.company_id = p.id
        WHERE mr.message_id = ANY($1)
    `

	var reactions []Reaction
	err := pgxscan.Select(ctx, r.db, &reactions, query, messageIDs)
	if err != nil {
		return nil, err
	}

	return reactions, nil
}

func (r *Repository) GetPinnedMessages(conversationID int) ([]MessageDetails, error) {
	query := `
        SELECT m.id, m.uuid, m.conversation_id, m.sender_id, m.message_type, 
               m.content, m.reply_to_id, m.forwarded_from_id, m.media_id, 
               m.sticker_id, m.is_edited, m.is_pinned, m.is_silent, m.created_at,
		TRIM(
				COALESCE(p.first_name,'') || ' ' || COALESCE(p.last_name,'') || ' ' || COALESCE(p.company_name, '') || 
				COALESCE(d.first_name,'') || ' ' || COALESCE(d.last_name,'') 
			) AS sender_name,
			COALESCE(p.image_url, d.image_url) AS sender_avatar
		FROM tbl_message m
		JOIN tbl_user u ON m.sender_id = u.id
		LEFT JOIN tbl_company p ON u.company_id = p.id
		LEFT JOIN tbl_driver d ON u.driver_id = d.id
        WHERE m.conversation_id = $1 AND m.active = 1 AND m.is_pinned = true
        ORDER BY m.created_at DESC
    `

	var pinnedMessages []MessageDetails
	err := pgxscan.Select(context.Background(), r.db, &pinnedMessages, query, conversationID)
	return pinnedMessages, err
}
