package services

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	repo "texApi/internal/_other/repositories"
	"texApi/internal/_other/schemas/request"
)

func GetUsers(ctx *gin.Context) {
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

	users, err := repo.GetUsers(offset, limit)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, users)
}

func GetUser(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "user" && role != "admin" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	user := repo.GetUser(id)

	ctx.JSON(200, user)
}

// func CreateUser(ctx *gin.Context) {
// 	var user request.CreateUser

// 	validationError := ctx.BindJSON(&user)

// 	if validationError != nil {
// 		ctx.JSON(400, validationError.Error())
// 		return
// 	}

// 	exist := repo.CheckUserExistWithStatus(user.Phone)

// 	if exist.Phone != "" && exist.IsVerified == true {
// 		ctx.JSON(400, gin.H{"message": "This phone number is already exists"})
// 		return
// 	}

// 	otp, otpErr := utils.SendOTP(user.Phone)

// 	if otpErr != nil {
// 		log.Println(otpErr.Error())
// 		ctx.JSON(500, otpErr.Error())
// 		return
// 	}

// 	hashedPassword, _ := bcrypt.GenerateFromPassword(
// 		[]byte(user.Password), 10,
// 	)

// 	user.Password = string(hashedPassword)

// 	userCreateErr := repo.CreateUser(user)

// 	if userCreateErr != nil {
// 		log.Println(userCreateErr.Error())
// 		ctx.JSON(500, userCreateErr.Error())
// 		return
// 	}

// 	ctx.JSON(201, gin.H{"id": otp.ID})
// }

func UpdateUser(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "user" {
		ctx.AbortWithStatus(403)
		return
	}

	var user request.UpdateUser

	validationError := ctx.BindJSON(&user)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	if user.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword(
			[]byte(user.Password), 10,
		)
		user.Password = string(hashedPassword)
	}

	err := repo.UpdateUser(user)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully updated"})
}

func DeleteUser(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "user" {
		ctx.AbortWithStatus(403)
		return
	}

	idStr := ctx.Param("id")
	id, error := strconv.Atoi(idStr)

	if error != nil || id == 0 {
		ctx.JSON(400, gin.H{"message": "ID must be positive integer number"})
		return
	}

	err := repo.DeleteUser(id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, gin.H{"message": "Successfully deleted"})
}
