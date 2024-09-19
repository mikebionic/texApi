package services

import (
	"log"
	"net/http"
	repo "texApi/pkg/repositories"

	"github.com/gin-gonic/gin"
)

func GetContents(ctx *gin.Context) {
	contents, err := repo.GetContents()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, contents)
}
