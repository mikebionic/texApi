package services

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	repo "texApi/internal/_other/repositories"
	"texApi/internal/_other/schemas/request"
)

func GetSubscriptions(ctx *gin.Context) {
	subscriptions, err := repo.GetSubscriptions()

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, subscriptions)
}

func GetSubscription(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	subscription := repo.GetSubscription(id)

	ctx.JSON(200, subscription)
}

func CreateSubscription(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var subscription request.CreateSubscription

	validationError := ctx.BindJSON(&subscription)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	subID, err := repo.CreateSubscription(subscription)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(201, gin.H{"id": subID})
}

func UpdateSubscription(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var subscription request.UpdateSubscription

	validationError := ctx.BindJSON(&subscription)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	err := repo.UpdateSubscription(subscription)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func DeleteSubscription(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	err := repo.DeleteSubscription(id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully deleted"})
}
