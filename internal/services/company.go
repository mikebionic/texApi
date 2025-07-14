package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/internal/repo"
	"texApi/pkg/utils"
	"time"
)

func GetCompanyList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	filters := map[string]interface{}{
		"user_id":              ctx.Query("user_id"),
		"role":                 ctx.Query("role"),
		"role_id":              ctx.Query("role_id"),
		"plan":                 ctx.Query("plan"),
		"plan_active":          ctx.Query("plan_active"),
		"company_name":         ctx.Query("company_name"),
		"first_name":           ctx.Query("first_name"),
		"last_name":            ctx.Query("last_name"),
		"phone":                ctx.Query("phone"),
		"email":                ctx.Query("email"),
		"country":              ctx.Query("country"),
		"country_id":           ctx.Query("country_id"),
		"city_id":              ctx.Query("city_id"),
		"verified":             ctx.Query("verified"),
		"confirmation_request": ctx.Query("confirmation_request"),
		"entity":               ctx.Query("entity"),
		"featured":             ctx.Query("featured"),
		"rating":               ctx.Query("rating"),
		"partner":              ctx.Query("partner"),
		"blocked":              ctx.Query("blocked"),
		"active":               ctx.DefaultQuery("active", "1"),
	}

	createdStart := ctx.Query("created_start")
	createdEnd := ctx.Query("created_end")
	updatedStart := ctx.Query("updated_start")
	updatedEnd := ctx.Query("updated_end")
	lastActiveStart := ctx.Query("last_active_start")
	lastActiveEnd := ctx.Query("last_active_end")

	orderBy := ctx.DefaultQuery("order_by", "id")
	orderDir := ctx.DefaultQuery("order_dir", "DESC")

	search := ctx.Query("search")

	stmt := `
        SELECT 
            c.*,
            COUNT(*) OVER() as total_count,
            json_agg(DISTINCT d.*) FILTER (WHERE d.id IS NOT NULL) as drivers,
            json_agg(DISTINCT v.*) FILTER (WHERE v.id IS NOT NULL) as vehicles
        FROM tbl_company c
        LEFT JOIN tbl_driver d ON c.id = d.company_id AND d.deleted = 0
        LEFT JOIN tbl_vehicle v ON c.id = v.company_id AND v.deleted = 0
        WHERE c.deleted = 0
    `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	for key, value := range filters {
		if value != "" && value != nil {
			if key == "company_name" || key == "first_name" || key == "last_name" || key == "phone" || key == "email" {
				whereClauses = append(whereClauses, fmt.Sprintf("c.%s ILIKE $%d", key, argCounter))
				args = append(args, "%"+value.(string)+"%")
			} else {
				whereClauses = append(whereClauses, fmt.Sprintf("c.%s = $%d", key, argCounter))
				args = append(args, value)
			}
			argCounter++
		}
	}

	if search != "" {
		searchClause := fmt.Sprintf(`(
            c.company_name ILIKE $%d OR 
            c.first_name ILIKE $%d OR 
            c.last_name ILIKE $%d OR 
            c.phone ILIKE $%d OR 
            c.email ILIKE $%d OR 
            c.about ILIKE $%d
        )`, argCounter, argCounter, argCounter, argCounter, argCounter, argCounter)
		whereClauses = append(whereClauses, searchClause)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern)
		argCounter++
	}

	if createdStart != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.created_at >= $%d", argCounter))
		startTime, _ := time.Parse(time.RFC3339, createdStart)
		args = append(args, startTime)
		argCounter++
	}
	if createdEnd != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.created_at <= $%d", argCounter))
		endTime, _ := time.Parse(time.RFC3339, createdEnd)
		args = append(args, endTime)
		argCounter++
	}

	if updatedStart != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.updated_at >= $%d", argCounter))
		startTime, _ := time.Parse(time.RFC3339, updatedStart)
		args = append(args, startTime)
		argCounter++
	}
	if updatedEnd != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.updated_at <= $%d", argCounter))
		endTime, _ := time.Parse(time.RFC3339, updatedEnd)
		args = append(args, endTime)
		argCounter++
	}

	if lastActiveStart != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.last_active >= $%d", argCounter))
		startTime, _ := time.Parse(time.RFC3339, lastActiveStart)
		args = append(args, startTime)
		argCounter++
	}
	if lastActiveEnd != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.last_active <= $%d", argCounter))
		endTime, _ := time.Parse(time.RFC3339, lastActiveEnd)
		args = append(args, endTime)
		argCounter++
	}

	if len(whereClauses) > 0 {
		stmt += " AND " + strings.Join(whereClauses, " AND ")
	}

	stmt += " GROUP BY c.id"

	stmt += fmt.Sprintf(" ORDER BY c.%s %s", orderBy, orderDir)

	stmt += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var companies []dto.CompanyDetails
	err := pgxscan.Select(ctx, db.DB, &companies, stmt, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Couldn't retrieve data", err.Error()))
		return
	}

	var totalCount int
	if len(companies) > 0 {
		totalCount = companies[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    companies,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company list", response))
}

func GetCompany(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Company ID parameter should be int", err.Error()))
	}
	company, err := repo.GetCompanyByID(id)
	ctx.JSON(http.StatusOK, utils.FormatResponse("Company details", company))
	return
}

// // TODO: should they be able to create only one company? check if company of that user exists?
func CreateCompany(ctx *gin.Context) {
	var company dto.CompanyCreate

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	userID := ctx.MustGet("id").(int)
	roleID := ctx.MustGet("roleID").(int)
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		company.UserID = userID
		company.RoleID = roleID
	}

	var companyID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateCompany,
		company.UserID,
		company.RoleID,
		company.CompanyName,
		company.FirstName,
		company.LastName,
		company.PatronymicName,
		company.Phone,
		company.Phone2,
		company.Phone3,
		company.Email,
		company.Email2,
		company.Email3,
		company.Meta,
		company.Meta2,
		company.Meta3,
		company.Address,
		company.Country,
		company.CountryID,
		company.CityID,
		company.ImageURL,
		company.Entity,
	).Scan(&companyID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating company", err.Error()))
		return
	}

	_, err = db.DB.Exec(
		context.Background(),
		queries.UpdateUserCompany,
		companyID,
		company.UserID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating user company", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created!", gin.H{"id": companyID}))
}

func UpdateCompany(ctx *gin.Context) {
	id := ctx.Param("id")
	var company dto.CompanyUpdate

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	stmt := queries.UpdateCompany

	userID := ctx.MustGet("id").(int)
	roleID := ctx.MustGet("roleID").(int)
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		company.UserID = &userID
		company.RoleID = &roleID
		company.Active = nil
		company.Deleted = nil
		stmt += ` WHERE (id = $1 AND user_id = $21) AND (active = 1 AND deleted = 0)`
	} else {
		stmt += ` WHERE id = $1`
	}

	stmt += ` RETURNING id;`

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		stmt,
		id,
		company.CompanyName,
		company.FirstName,
		company.LastName,
		company.PatronymicName,
		company.Phone,
		company.Phone2,
		company.Phone3,
		company.Email,
		company.Email2,
		company.Email3,
		company.Meta,
		company.Meta2,
		company.Meta3,
		company.Address,
		company.Country,
		company.CountryID,
		company.CityID,
		company.ImageURL,
		company.Entity,
		company.UserID,
		company.RoleID,
		company.Active,
		company.Deleted,
		company.Partner,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating company", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

func DeleteCompany(ctx *gin.Context) {
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Operation can't be done by user", ""))
		return
	}

	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteCompany,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting company", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}
