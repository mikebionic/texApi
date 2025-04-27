package queries

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BuildFilteredQuery builds a SQL query with filters, ordering, and pagination
func BuildFilteredQuery(baseQuery string, filters map[string]interface{}, orderBy, orderDir string, validOrderColumns map[string]bool) (string, []interface{}, int) {
	query := baseQuery
	args := []interface{}{}
	paramCount := 0

	// Add filters
	for key, value := range filters {
		if value != nil && value != "" {
			if strings.Contains(key, "LIKE") {
				// Handle LIKE filters
				field := strings.Split(key, " ")[0]
				query += fmt.Sprintf(" AND LOWER(%s) LIKE LOWER($%d)", field, paramCount+1)
				args = append(args, "%"+value.(string)+"%")
			} else {
				// Handle exact match filters
				query += fmt.Sprintf(" AND %s = $%d", key, paramCount+1)
				args = append(args, value)
			}
			paramCount++
		}
	}

	// Add ordering
	if validOrderColumns[strings.ToLower(orderBy)] {
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
	}

	return query, args, paramCount
}

// ExecuteUpdate executes an update query and returns the number of affected rows
func ExecuteUpdate(ctx context.Context, conn *pgxpool.Pool, query string, args ...interface{}) (int64, error) {
	commandTag, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), nil
}

// ExecuteCreate executes an insert query and returns the generated ID
func ExecuteCreate(ctx context.Context, conn *pgxpool.Pool, query string, args ...interface{}) (int, error) {
	var id int
	err := conn.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateCompanyVerification updates the verification status of a company
func UpdateCompanyVerification(ctx context.Context, conn *pgxpool.Pool, companyID int, verified int) (int64, error) {
	query := `
		UPDATE tbl_company
		SET verified = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted = 0
	`
	return ExecuteUpdate(ctx, conn, query, verified, companyID)
}

// UpdateCompanyPlanActive updates the plan_active status of a company
func UpdateCompanyPlanActive(ctx context.Context, conn *pgxpool.Pool, companyID int, planActive int) (int64, error) {
	query := `
		UPDATE tbl_company
		SET plan_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted = 0
	`
	return ExecuteUpdate(ctx, conn, query, planActive, companyID)
}

// CheckExpiredPlans checks for expired plans and updates company status
func CheckExpiredPlans(ctx context.Context, conn *pgxpool.Pool) (int64, error) {
	query := `
		WITH expired_plans AS (
			SELECT company_id
			FROM tbl_plan_moves
			WHERE status = 'approved'
			AND valid_until < CURRENT_TIMESTAMP
			AND deleted = 0
		)
		UPDATE tbl_company p
		SET plan_active = 0, updated_at = CURRENT_TIMESTAMP
		FROM expired_plans e
		WHERE p.id = e.company_id
		AND p.plan_active = 1
		AND p.deleted = 0
	`
	return ExecuteUpdate(ctx, conn, query)
}

// ExtendPlanValidity extends the validity of a plan
func ExtendPlanValidity(ctx context.Context, conn *pgxpool.Pool, planMoveID int) (time.Time, error) {
	query := `
		UPDATE tbl_plan_moves
		SET valid_until = CURRENT_TIMESTAMP + INTERVAL '1 month', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted = 0
		RETURNING valid_until
	`
	var validUntil time.Time
	err := conn.QueryRow(ctx, query, planMoveID).Scan(&validUntil)
	if err != nil {
		return time.Time{}, err
	}
	return validUntil, nil
}
