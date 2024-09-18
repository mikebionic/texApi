package services

import (
	"log"
	"strconv"
	repo "texApi/pkg/repositories"
	"texApi/pkg/schemas/request"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetServices(ctx *gin.Context) {
	services, err := repo.GetServices()

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, services)
}

func GetServiceList(ctx *gin.Context) {
	services, err := repo.GetServiceList()

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, services)
}

func GetService(ctx *gin.Context) {
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

	service := repo.GetService(id)

	ctx.JSON(200, service)
}

func CreateService(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var service request.CreateService

	validationError := ctx.BindJSON(&service)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	serviceID, err := repo.CreateService(service)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(201, gin.H{"id": serviceID})
}

func UpdateService(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var service request.UpdateService

	validationError := ctx.BindJSON(&service)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	err := repo.UpdateService(service)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func DeleteService(ctx *gin.Context) {
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

	err := repo.DeleteService(id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully deleted"})
}

func SetServiceImage(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	filename := utils.WriteImage(ctx, "services/")

	if filename == "" {
		ctx.JSON(400, gin.H{"message": "Only jpeg, jpg, png, svg, webp"})
		return
	}

	err := repo.SetServiceImage(id, filename)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Has been saved"})
}
