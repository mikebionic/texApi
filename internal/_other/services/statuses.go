package services

import (
	"github.com/gin-gonic/gin"
	"log"
	repo "texApi/internal/_other/repositories"
)

func GetStatuses(ctx *gin.Context) {
	statuses, err := repo.GetStatuses()

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, statuses)
}
