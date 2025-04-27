package repo

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
)

func CreateSession(session dto.CreateSessionInput) (int, error) {
	var id int
	err := db.DB.QueryRow(context.Background(), `
        INSERT INTO tbl_sessions (
            user_id, company_id, refresh_token, expires_at, device_name, 
            device_model, device_firmware, app_name, app_version, user_agent,
            ip_address, login_method
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `,
		session.UserID, session.CompanyID, session.RefreshToken, session.ExpiresAt,
		session.DeviceName, session.DeviceModel, session.DeviceFirmware,
		session.AppName, session.AppVersion, session.UserAgent,
		session.IPAddress, session.LoginMethod).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetSessionByRefreshToken(refreshToken string) (dto.Session, error) {
	var session dto.Session
	err := pgxscan.Get(context.Background(), db.DB, &session, `
        SELECT * FROM tbl_sessions
        WHERE refresh_token = $1 AND expires_at > NOW() AND is_active = TRUE
    `, refreshToken)

	if err != nil {
		return dto.Session{}, err
	}

	_, err = db.DB.Exec(context.Background(), `
        UPDATE tbl_sessions
        SET last_used_at = NOW(), updated_at = NOW()
        WHERE id = $1
    `, session.ID)

	if err != nil {
		return dto.Session{}, err
	}

	return session, nil
}

func InvalidateSession(refreshToken string) error {
	_, err := db.DB.Exec(context.Background(), `
        UPDATE tbl_sessions
        SET is_active = FALSE, updated_at = NOW()
        WHERE refresh_token = $1
    `, refreshToken)

	return err
}

func InvalidateAllUserSessions(userID int) error {
	_, err := db.DB.Exec(context.Background(), `
        UPDATE tbl_sessions
        SET is_active = FALSE, updated_at = NOW()
        WHERE user_id = $1
    `, userID)

	return err
}

func GetUserSessions(userID int) ([]dto.Session, error) {
	var sessions []dto.Session
	err := pgxscan.Select(context.Background(), db.DB, &sessions, `
        SELECT * FROM tbl_sessions
        WHERE user_id = $1 AND is_active = TRUE
        ORDER BY created_at DESC
    `, userID)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func ListSessions(params dto.SessionListParams) ([]dto.SessionListItem, int, error) {
	baseQuery := `
        SELECT s.*, u.username, u.email, u.phone, u.role, count(*) OVER() as total_count 
        FROM tbl_sessions s
        LEFT JOIN tbl_user u ON s.user_id = u.id
        WHERE 1=1
    `

	conditions := []string{}
	args := []interface{}{}
	paramCount := 1

	if params.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("s.user_id = $%d", paramCount))
		args = append(args, *params.UserID)
		paramCount++
	}

	if params.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("s.company_id = $%d", paramCount))
		args = append(args, *params.CompanyID)
		paramCount++
	}

	if params.LoginMethod != nil {
		conditions = append(conditions, fmt.Sprintf("s.login_method = $%d", paramCount))
		args = append(args, *params.LoginMethod)
		paramCount++
	}

	if params.DeviceName != nil {
		conditions = append(conditions, fmt.Sprintf("s.device_name ILIKE $%d", paramCount))
		args = append(args, "%"+*params.DeviceName+"%")
		paramCount++
	}

	if params.AppName != nil {
		conditions = append(conditions, fmt.Sprintf("s.app_name = $%d", paramCount))
		args = append(args, *params.AppName)
		paramCount++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("s.is_active = $%d", paramCount))
		args = append(args, *params.IsActive)
		paramCount++
	}

	if params.CreatedFrom != nil {
		conditions = append(conditions, fmt.Sprintf("s.created_at >= $%d", paramCount))
		args = append(args, *params.CreatedFrom)
		paramCount++
	}

	if params.CreatedTo != nil {
		conditions = append(conditions, fmt.Sprintf("s.created_at <= $%d", paramCount))
		args = append(args, *params.CreatedTo)
		paramCount++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	validOrderColumns := map[string]bool{
		"id": true, "user_id": true, "login_method": true, "created_at": true, "last_used_at": true,
	}

	orderBy := "created_at"
	if validOrderColumns[params.OrderBy] {
		orderBy = params.OrderBy
	}

	orderDir := "DESC"
	if params.OrderDir == "ASC" {
		orderDir = "ASC"
	}

	baseQuery += fmt.Sprintf(" ORDER BY s.%s %s", orderBy, orderDir)

	offset := (params.Page - 1) * params.PerPage
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, params.PerPage, offset)

	var sessions []dto.SessionListItem
	err := pgxscan.Select(context.Background(), db.DB, &sessions, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	totalCount := 0
	if len(sessions) > 0 {
		totalCount = sessions[0].TotalCount
	}

	return sessions, totalCount, nil
}

func CleanExpiredSessions() (int, error) {
	tag, err := db.DB.Exec(context.Background(), `
        DELETE FROM tbl_sessions
        WHERE expires_at < NOW()
    `)

	if err != nil {
		return 0, err
	}

	return int(tag.RowsAffected()), nil
}

func GetSessionByID(sessionID int) (dto.Session, error) {
	var session dto.Session
	err := pgxscan.Get(context.Background(), db.DB, &session, `
        SELECT * FROM tbl_sessions
        WHERE id = $1
    `, sessionID)

	if err != nil {
		return dto.Session{}, err
	}

	return session, nil
}

func InvalidateSessionByID(sessionID int) error {
	_, err := db.DB.Exec(context.Background(), `
        UPDATE tbl_sessions
        SET is_active = FALSE, updated_at = NOW()
        WHERE id = $1
    `, sessionID)

	return err
}
