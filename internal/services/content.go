package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"texApi/internal/dto"
	repo "texApi/internal/repositories"
)

func GetContents(ctx *gin.Context) {
	ctID, err := strconv.Atoi(ctx.GetHeader("ContentTypeId"))
	fmt.Println(ctID)
	contents, err := repo.GetContents(ctID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contents)
}

func GetContent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "ID must be positive integer number"})
		return
	}

	content := repo.GetContent(id)
	ctx.JSON(http.StatusOK, content)
}

func CreateContent(ctx *gin.Context) {
	var content dto.CreateContent
	validationError := ctx.BindJSON(&content)
	if validationError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": validationError.Error()})
		return
	}

	cID := repo.CreateContent(content)
	if cID == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Can not create content"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": cID})
}

func DeleteContent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "ID must be positive integer number"})
		return
	}
	err = repo.DeleteContent(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id, "message": "Successfully deleted"})
}
