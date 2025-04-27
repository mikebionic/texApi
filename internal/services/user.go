package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"

	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

// TODO: ALERT REFACTOR THIS?!!!!

func GetUserList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	search := ctx.Query("search")
	role := ctx.Query("role")
	active := ctx.Query("active")
	verified := ctx.Query("verified")
	orderBy := ctx.DefaultQuery("order_by", "id")
	orderDir := ctx.DefaultQuery("order_dir", "ASC")

	query := `
		SELECT *, count(*) OVER() as total_count 
		FROM tbl_user 
		WHERE deleted = 0
	`
	args := []interface{}{perPage, offset}
	paramCount := 2

	if search != "" {
		query += fmt.Sprintf(` AND (
			LOWER(username) LIKE LOWER($%d) OR
			LOWER(email) LIKE LOWER($%d) OR
			LOWER(phone) LIKE LOWER($%d)
		)`, paramCount+1, paramCount+1, paramCount+1)
		args = append(args, "%"+search+"%")
		paramCount++
	}

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", paramCount+1)
		args = append(args, role)
		paramCount++
	}

	if active != "" {
		query += fmt.Sprintf(" AND active = $%d", paramCount+1)
		args = append(args, active)
		paramCount++
	}

	if verified != "" {
		query += fmt.Sprintf(" AND verified = $%d", paramCount+1)
		args = append(args, verified)
		paramCount++
	}

	validOrderColumns := map[string]bool{
		"id": true, "username": true, "email": true, "role": true,
	}
	if validOrderColumns[strings.ToLower(orderBy)] {
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
	}

	query += " LIMIT $1 OFFSET $2"

	var users []dto.UserDetails
	err := pgxscan.Select(context.Background(), db.DB, &users, query, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	totalCount := 0
	if len(users) > 0 {
		totalCount = users[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    users,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("User list", response))
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		SELECT * 
		FROM tbl_user 
		WHERE id = $1 AND deleted = 0
	`

	var user dto.UserDetails
	err := pgxscan.Get(context.Background(), db.DB, &user, query, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("User details", user))
}

func GetUserRichInfo(ctx *gin.Context) {
	id := ctx.Param("id")

	userQuery := `
        SELECT * FROM tbl_user
        WHERE id = $1 AND deleted = 0
    `
	var user dto.UserDetails
	err := pgxscan.Get(context.Background(), db.DB, &user, userQuery, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found", err.Error()))
		return
	}

	richInfo := dto.UserRichInfo{
		User:    user,
		Company: dto.CompanyDetails{},
	}

	if user.CompanyID != 0 {
		companyQuery := `
			SELECT * FROM tbl_company
			WHERE id = $1
		`
		err = pgxscan.Get(context.Background(), db.DB, &richInfo.Company, companyQuery, user.CompanyID)
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("User rich info", richInfo))
}

func CreateUser(ctx *gin.Context) {
	var user dto.UserCreate

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// TODO: Hash password
	//hashedPassword, err := utils.HashPassword(user.Password)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Password hashing error", err.Error()))
	//	return
	//}
	hashedPassword := user.Password

	query := `
		INSERT INTO tbl_user (
			username, email, phone, password, role, meta, meta2, meta3, verified
		) VALUES (
			$1, $2, $3, $4, COALESCE($5, 'role')::role_t, COALESCE($6, 'meta'),
			COALESCE($7, 'meta2'),COALESCE($8, 'meta3'), COALESCE($9, 0)
		) RETURNING id
	`

	var userID int
	err := db.DB.QueryRow(
		context.Background(),
		query,
		user.Username, user.Email, user.Phone,
		hashedPassword, user.Role, user.Meta,
		user.Meta2, user.Meta3, user.Verified,
	).Scan(&userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating user", err.Error()))
		return
	}

	queryCompany := `
		INSERT INTO tbl_company (
			first_name, last_name, email, phone, role, user_id
		) VALUES (
			$1,  COALESCE($2, 'last_name'), $3, $4, COALESCE($5, 'role')::role_t, $6
		) RETURNING id
	`

	var companyID int
	err = db.DB.QueryRow(
		context.Background(),
		queryCompany,
		user.FirstName, user.LastName, user.Email, user.Phone, user.Role, userID,
	).Scan(&companyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating company", err.Error()))
	}

	_, err = db.DB.Exec(
		context.Background(),
		`UPDATE tbl_user SET company_id = $1 WHERE id = $2;`,
		companyID, userID,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating user's company_id", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created!", gin.H{"user_id": userID, "company_id": companyID}))
}

func UpdateUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var user dto.UserUpdate

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Optional password hashing
	var hashedPassword *string
	if user.Password != nil {
		//TODO: HASHED PASSWORD
		//hashed, err := utils.HashPassword(*user.Password)
		//if err != nil {
		//	ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Password hashing error", err.Error()))
		//	return
		//}
		//hashedPassword = &hashed
		hashedPassword = user.Password
	}

	query := `
		UPDATE tbl_user SET
			username = COALESCE($1, username),
			email = COALESCE($2, email),
			phone = COALESCE($3, phone),
			password = COALESCE($4, password),
			role = COALESCE($5, role)::role_t,
			active = COALESCE($6, active),
			verified = COALESCE($7, verified),
			meta = COALESCE($8,meta),
			meta2 = COALESCE($9,meta2),
			meta3 = COALESCE($10,meta3),
			company_id = COALESCE($11,company_id)
		WHERE id = $12 AND deleted = 0
	`

	commandTag, err := db.DB.Exec(
		context.Background(),
		query,
		user.Username, user.Email, user.Phone,
		hashedPassword, user.Role, user.Active,
		user.Verified, user.Meta, user.Meta2,
		user.Meta3, user.CompanyID, id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating user", err.Error()))
		return
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found or already deleted", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": id}))
}

func DeleteUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	query := `
		UPDATE tbl_user 
		SET deleted = 1
		WHERE id = $1 AND deleted = 0
	`

	commandTag, err := db.DB.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting user", err.Error()))
		return
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found or already deleted", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}
