package services

import (
	"log"
	"strconv"
	repo "texApi/internal/_other/repositories"
	"texApi/internal/_other/schemas/request"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetWorkers(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	pageStr := ctx.Query("page")
	countStr := ctx.Query("count")

	pageInt, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(countStr)

	// set default values if not presented
	if pageInt == 0 {
		pageInt = 1
	}

	if limit == 0 {
		limit = 20
	}

	offset := pageInt*limit - limit

	workers, err := repo.GetWorkers(offset, limit)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, workers)
}

func GetWorker(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" && role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	worker := repo.GetWorker(id)

	ctx.JSON(200, worker)
}

func GetWorkerMe(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" {
		ctx.AbortWithStatus(403)
		return
	}

	id := ctx.MustGet("id").(int)

	worker := repo.GetWorker(id)

	ctx.JSON(200, worker)
}

func CreateWorker(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var worker request.CreateWorker

	validationError := ctx.BindJSON(&worker)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	exist := repo.CheckWorkerExist(worker.Phone)

	if exist != "" {
		ctx.JSON(400, "This phone number is already exists")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(worker.Password), 10,
	)
	worker.Password = string(hashedPassword)

	workerID, err := repo.CreateWorker(worker)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(201, gin.H{"id": workerID})
}

func UpdateWorker(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" && role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	var worker request.UpdateWorker

	validationError := ctx.BindJSON(&worker)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	if worker.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword(
			[]byte(worker.Password), 10,
		)

		worker.Password = string(hashedPassword)
	}

	err := repo.UpdateWorker(worker)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func DeleteWorker(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" && role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	err := repo.DeleteWorker(id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully deleted"})
}

func SetWorkerImage(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "worker" && role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	if id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	var workerID int

	if id == 0 {
		workerID = ctx.MustGet("id").(int)
	} else {
		workerID = id
	}

	filename := utils.WriteImage(ctx, "workers/")

	if filename == "" {
		ctx.JSON(400, gin.H{"message": "Only jpeg, jpg, png, svg, webp"})
		return
	}

	err := repo.SetWorkerImage(workerID, filename)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Has been saved"})
}
