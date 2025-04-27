package repo

import (
	"context"
	"fmt"
	"log"
	db "texApi/database"
	"texApi/internal/dto"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func CreateCompanyShort(company dto.CompanyCreateShort) (int, error) {
	var result []struct{ ID int }

	query := `
        INSERT INTO tbl_company (
            user_id, role, role_id, first_name, last_name, company_name, 
            about, phone, email, image_url
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        ) RETURNING id
    `

	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&result,
		query,
		company.UserID, company.Role, company.RoleID,
		company.FirstName, company.LastName, company.CompanyName,
		company.About, company.Phone, company.Email,
		company.ImageURL,
	)

	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, fmt.Errorf("user creation failed")
	}
	return result[0].ID, nil
}

func CreateCompany(company dto.CompanyCreate) (int, error) {
	var result []struct{ ID int }

	query := `
        INSERT INTO tbl_company (
            user_id, role, role_id, plan, first_name, last_name, company_name, 
            about, phone, email, address, image_url, 
            show_avatar, show_bio, show_last_seen, show_phone_number,
            self_destruct_duration, avatar_exceptions, bio_exceptions, 
            last_seen_exceptions, phone_number_exceptions,
            receive_calls, invite_group, notifications_chat, 
            notifications_group, notifications_story, notifications_reactions,
            receive_calls_exceptions, invite_group_exceptions
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 
            $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23,
            $24, $25, $26, $27, $28, $29
        ) RETURNING id
    `

	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&result,
		query,
		company.UserID, company.Role, company.RoleID, company.Plan,
		company.FirstName, company.LastName, company.CompanyName,
		company.About, company.Phone, company.Email,
		company.Address, company.ImageURL,
		company.ConfirmationRequest,
		company.ShowAvatar, company.ShowBio,
		company.ShowLastSeen, company.ShowPhoneNumber,
		company.SelfDestructDuration,
		company.AvatarExceptions, company.BioExceptions,
		company.LastSeenExceptions, company.PhoneNumberExceptions,
		company.ReceiveCalls, company.InviteGroup,
		company.NotificationsChat, company.NotificationsGroup,
		company.NotificationsStory, company.NotificationsReactions,
		company.ReceiveCallsExceptions, company.InviteGroupExceptions,
	)

	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, fmt.Errorf("user creation failed")
	}
	return result[0].ID, nil
}

func UpdateUserCompanyID(userID, companyID int) error {
	var id int
	query := `
        UPDATE tbl_user SET company_id = $1
        WHERE id = $2 RETURNING id
    `
	err := db.DB.QueryRow(
		context.Background(), query, companyID, userID,
	).Scan(&id)
	if err != nil {
		log.Println("UpdateUserCompanyID: ", err.Error())
	}
	return err
}
