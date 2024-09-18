package services

import (
	"log"
	repo "texApi/pkg/repositories"

	"github.com/gin-gonic/gin"
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
