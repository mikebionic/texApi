package services

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"texApi/internal/dto"
	repo "texApi/internal/repositories"
	"texApi/pkg/utils"
)

func GetContents(ctx *gin.Context) {
	ctID, err := strconv.Atoi(ctx.GetHeader("ContentTypeId"))
	contents, err := repo.GetContents(ctID)
	if err != nil {
		log.Println(err.Error())
		response := utils.FormatErrorResponse("Failed to retrieve contents", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("Contents retrieved successfully", contents)
	ctx.JSON(http.StatusOK, response)
}

func GetContent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		response := utils.FormatErrorResponse("ID must be positive integer number", "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	content := repo.GetContent(id)
	response := utils.FormatResponse("Content retrieved successfully", content)
	ctx.JSON(http.StatusOK, response)
}

func CreateContent(ctx *gin.Context) {
	var content dto.CreateContent
	validationError := ctx.BindJSON(&content)
	if validationError != nil {
		response := utils.FormatErrorResponse("Invalid request body", validationError.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	cID := repo.CreateContent(content)
	if cID == 0 {
		response := utils.FormatResponse("Cannot create content", "")
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("Content created successfully", gin.H{"id": cID})
	ctx.JSON(http.StatusCreated, response)
}

func UpdateContent(ctx *gin.Context) {
	var content dto.CreateContent
	validationError := ctx.BindJSON(&content)
	if validationError != nil {
		response := utils.FormatErrorResponse("Invalid request body", validationError.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response := utils.FormatErrorResponse("Invalid content ID", "Content ID must be a positive integer")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	updatedID, err := repo.UpdateContent(content, id)
	if err != nil {
		response := utils.FormatResponse("Cannot update content", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.FormatResponse("Content updated successfully", repo.GetContent(updatedID))
	ctx.JSON(http.StatusOK, response)
}

func DeleteContent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		response := utils.FormatErrorResponse("ID must be positive integer number", "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	err = repo.DeleteContent(id)
	if err != nil {
		response := utils.FormatErrorResponse("Failed to delete content", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("Successfully deleted", gin.H{"id": id})
	ctx.JSON(http.StatusOK, response)
}
