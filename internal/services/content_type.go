package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	repo "texApi/internal/repositories"
)

func GetContentTypes(ctx *gin.Context) {
	withContent, err := strconv.Atoi(ctx.GetHeader("WithContent"))
	contentTypes, err := repo.GetContentTypes(withContent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contentTypes)
}
