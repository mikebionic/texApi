package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	repo "texApi/internal/repositories"
)

func GetContentTypes(ctx *gin.Context) {
	contentTypes, err := repo.GetContentTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"content_types": contentTypes})
}
