package firebasePush

import (
	"context"
	"fmt"
	"log"
	db "texApi/database"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func SaveToken(userID int, req SaveNotificationTokenRequest) (err error) {
	var existingToken FirebaseToken
	checkQuery := `
		SELECT id FROM tbl_firebase_token 
		WHERE user_id = $1 AND token = $2
	`
	err = pgxscan.Get(context.Background(), db.DB, &existingToken, checkQuery, userID, req.NotificationToken)
	if err == nil {
		// if exists - make active
		updateQuery := `
			UPDATE tbl_firebase_token 
			SET active = 1, updated_at = NOW(), device_type = $2
			WHERE id = $1
		`
		_, err = db.DB.Exec(context.Background(), updateQuery, existingToken.ID, req.DeviceType)
		if err != nil {
			return
		}
	} else {
		insertQuery := `
			INSERT INTO tbl_firebase_token (user_id, token, device_type)
			VALUES ($1, $2, $3)
		`
		_, err = db.DB.Exec(context.Background(), insertQuery, userID, req.NotificationToken, req.DeviceType)
		if err != nil {
			return
		}
	}
	return
}

func GetUserTokens(userID int) ([]string, error) {
	var tokens []string
	query := `
		SELECT token FROM tbl_firebase_token 
		WHERE user_id = $1 AND active = 1 AND deleted = 0
	`

	rows, err := db.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying tokens: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			log.Printf("Error scanning token: %v", err)
			continue
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func removeInvalidToken(token string) {
	query := `UPDATE tbl_firebase_token SET active = 0, deleted = 1, updated_at = NOW() WHERE token = $1`
	_, err := db.DB.Exec(context.Background(), query, token)
	if err != nil {
		log.Printf("Error removing invalid token %s: %v", token, err)
	} else {
		log.Printf("Removed invalid token: %s", token)
	}
}

func RemoveInvalidUserTokens(userID int) {
	query := `UPDATE tbl_firebase_token SET active = 0, updated_at = NOW() WHERE user_id = $1`
	_, err := db.DB.Exec(context.Background(), query, userID)
	if err != nil {
		log.Printf("Error removing tokens for user_id %d: %v", userID, err)
	} else {
		log.Printf("Removed tokens for user_id: %d", userID)
	}
}

func GetConversationParticipants(conversationID, senderID int) ([]int, error) {
	var participants []int

	query := `
		SELECT DISTINCT user_id 
		FROM tbl_conversation_member 
		WHERE conversation_id = $1 AND user_id != $2 AND active = 1 AND deleted = 0
	`

	rows, err := db.DB.Query(context.Background(), query, conversationID, senderID)
	if err != nil {
		return nil, fmt.Errorf("error querying conversation participants: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Printf("Error scanning participant ID: %v", err)
			continue
		}
		participants = append(participants, userID)
	}

	return participants, nil
}
