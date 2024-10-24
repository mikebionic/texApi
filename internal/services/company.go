package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func CreateCompany(ctx *gin.Context) {
	var company dto.CompanyCreate

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		`INSERT INTO tbl_company (user_id, name, address, phone, email, logo_url) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id`,
		company.UserID, company.Name, company.Address, company.Phone, company.Email, company.LogoURL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating company"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"company_id": id})
}

func GetCompany(ctx *gin.Context) {
	id := ctx.Param("id")

	var company dto.CompanyCreate
	err := db.DB.QueryRow(
		context.Background(),
		queries.GetCompany,
		id,
	).Scan(&company.UserID, &company.Name, &company.Address, &company.Phone, &company.Email, &company.LogoURL)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	ctx.JSON(http.StatusOK, company)
}

func UpdateCompany(ctx *gin.Context) {
	id := ctx.Param("id")
	var company dto.CompanyUpdate

	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateCompany,
		id, company.Name, company.Address, company.Phone, company.Email, company.LogoURL,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating company"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"company_id": updatedID})
}

func DeleteCompany(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteCompany,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting company"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Company deleted"})
}
