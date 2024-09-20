package services

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	repo "texApi/internal/_other/repositories"
	"texApi/internal/_other/schemas/request"
)

func GetAboutUsAll(ctx *gin.Context) {
	aboutUs, err := repo.GetAboutUsAll()

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, aboutUs)
}

func GetAboutUs(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	aboutUsID, error := strconv.Atoi(idStr)

	if error != nil || aboutUsID == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	aboutUs := repo.GetAboutUs(aboutUsID)

	ctx.JSON(200, aboutUs)
}

func GetAboutUsForUser(ctx *gin.Context) {
	aboutUs := repo.GetAboutUsForUser()
	ctx.JSON(200, aboutUs)
}

func CreateAboutUs(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var aboutUs request.CreateAboutUs

	validationError := ctx.BindJSON(&aboutUs)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	aboutUsID, err := repo.CreateAboutUs(aboutUs)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(201, gin.H{"id": aboutUsID})
}

func UpdateAboutUs(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var aboutUs request.UpdateAboutUs

	validationError := ctx.BindJSON(&aboutUs)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	err := repo.UpdateAboutUs(aboutUs)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func DeleteAboutUs(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	aboutUsID, error := strconv.Atoi(idStr)

	if error != nil || aboutUsID == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	err := repo.DeleteAboutUs(aboutUsID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully deleted"})
}
