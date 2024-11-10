package services

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func GetCompanyList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		queries.GetCompanyWithRelations,
		perPage,
		offset,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var companies []dto.CompanyDetails
	var totalCount int

	for rows.Next() {
		var company dto.CompanyDetails
		var driversJSON, vehiclesJSON []byte

		err := rows.Scan(
			&company.ID, &company.UUID, &company.UserID, &company.RoleID,
			&company.CompanyName, &company.FirstName, &company.LastName,
			&company.PatronymicName, &company.Phone, &company.Phone2,
			&company.Phone3, &company.Email, &company.Email2, &company.Email3,
			&company.Meta, &company.Meta2, &company.Meta3, &company.Address,
			&company.Country, &company.CountryID, &company.CityID,
			&company.ImageURL, &company.Entity, &company.Featured,
			&company.Rating, &company.Partner, &company.SuccessfulOps,
			&company.CreatedAt, &company.UpdatedAt, &company.Active, &company.Deleted,
			&totalCount, &driversJSON, &vehiclesJSON,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}

		json.Unmarshal(driversJSON, &company.Drivers)
		json.Unmarshal(vehiclesJSON, &company.Vehicles)
		companies = append(companies, company)
	}

	response := dto.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    companies,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company list", response))
}

func GetCompany(ctx *gin.Context) {
	id := ctx.Param("id")

	var company dto.CompanyDetails
	var driversJSON, vehiclesJSON []byte

	err := db.DB.QueryRow(
		context.Background(),
		queries.GetCompanyByID,
		id,
	).Scan(
		&company.ID, &company.UUID, &company.UserID, &company.RoleID,
		&company.CompanyName, &company.FirstName, &company.LastName,
		&company.PatronymicName, &company.Phone, &company.Phone2,
		&company.Phone3, &company.Email, &company.Email2, &company.Email3,
		&company.Meta, &company.Meta2, &company.Meta3, &company.Address,
		&company.Country, &company.CountryID, &company.CityID,
		&company.ImageURL, &company.Entity, &company.Featured,
		&company.Rating, &company.Partner, &company.SuccessfulOps,
		&company.CreatedAt, &company.UpdatedAt, &company.Active, &company.Deleted,
		&driversJSON, &vehiclesJSON,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Company not found", err.Error()))
		return
	}

	json.Unmarshal(driversJSON, &company.Drivers)
	json.Unmarshal(vehiclesJSON, &company.Vehicles)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company details", company))
}

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
