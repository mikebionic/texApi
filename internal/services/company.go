package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func CreateCompany(ctx *gin.Context) {
	var company dto.CompanyCreate

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateCompany,
		company.UserID, company.Name, company.Address, company.Phone, company.Email, company.LogoURL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating company", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func GetCompany(ctx *gin.Context) {
	id := ctx.Param("id")

	var company dto.CompanyCreate
	err := db.DB.QueryRow(
		context.Background(),
		queries.GetCompany,
		id,
	).Scan(&company.ID, &company.UserID, &company.Name, &company.Address, &company.Phone, &company.Email, &company.LogoURL, &company.CreatedAt, &company.UpdatedAt, &company.Active, &company.Deleted)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company data", company))
}

func UpdateCompany(ctx *gin.Context) {
	id := ctx.Param("id")
	var company dto.CompanyUpdate

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateCompany,
		id, company.Name, company.Address, company.Phone, company.Email, company.LogoURL,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating company", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

func DeleteCompany(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}
