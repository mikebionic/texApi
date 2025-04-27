package services

import (
	"net/http"
	"strconv"
	"texApi/internal/dto"
	"texApi/internal/repo"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetPackagingTypes(ctx *gin.Context) {
	types, err := repo.GetPackagingTypes()
	if err != nil {
		response := utils.FormatErrorResponse("Failed to retrieve packaging types", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("Packaging types retrieved successfully", types)
	ctx.JSON(http.StatusOK, response)
}

func GetPackagingType(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response := utils.FormatErrorResponse("Invalid ID format", "ID must be a positive integer")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	packagingType, err := repo.GetPackagingType(id)
	if err != nil {
		response := utils.FormatErrorResponse("Packaging type not found", err.Error())
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	response := utils.FormatResponse("Packaging type retrieved successfully", packagingType)
	ctx.JSON(http.StatusOK, response)
}

func CreatePackagingType(ctx *gin.Context) {
	var packagingType dto.CreatePackagingType
	if err := ctx.BindJSON(&packagingType); err != nil {
		response := utils.FormatErrorResponse("Invalid request body", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id, err := repo.CreatePackagingType(packagingType)
	if err != nil {
		response := utils.FormatErrorResponse("Failed to create packaging type", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.FormatResponse("Packaging type created successfully", gin.H{"id": id})
	ctx.JSON(http.StatusCreated, response)
}

func UpdatePackagingType(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response := utils.FormatErrorResponse("Invalid ID format", "ID must be a positive integer")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var packagingType dto.CreatePackagingType
	if err := ctx.BindJSON(&packagingType); err != nil {
		response := utils.FormatErrorResponse("Invalid request body", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	updatedID, err := repo.UpdatePackagingType(packagingType, id)
	if err != nil {
		response := utils.FormatErrorResponse("Failed to update packaging type", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	updated, err := repo.GetPackagingType(updatedID)
	if err != nil {
		response := utils.FormatErrorResponse("Failed to retrieve updated packaging type", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.FormatResponse("Packaging type updated successfully", updated)
	ctx.JSON(http.StatusOK, response)
}

func DeletePackagingType(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response := utils.FormatErrorResponse("Invalid ID format", "ID must be a positive integer")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err := repo.DeletePackagingType(id); err != nil {
		response := utils.FormatErrorResponse("Failed to delete packaging type", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.FormatResponse("Packaging type deleted successfully", gin.H{"id": id})
	ctx.JSON(http.StatusOK, response)
}
