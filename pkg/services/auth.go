package services

import (
	"log"
	"texApi/config"
	repo "texApi/pkg/repositories"
	"texApi/pkg/schemas/request"
	"texApi/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(ctx *gin.Context) {
	var admin request.LoginForm

	validationError := ctx.BindJSON(&admin)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	findedAdmin := repo.GetAdmin(admin.Phone)

	compareError := bcrypt.CompareHashAndPassword(
		[]byte(findedAdmin.Password), []byte(admin.Password),
	)

	if compareError != nil {
		ctx.JSON(400, gin.H{"message": "Invalid password"})
		return
	}

	accessToken := utils.CreateToken(
		findedAdmin.ID, config.ENV.ACCESS_TIME,
		config.ENV.ACCESS_KEY, "admin",
	)

	refreshToken := utils.CreateToken(
		findedAdmin.ID, config.ENV.REFRESH_TIME,
		config.ENV.REFRESH_KEY, "admin",
	)

	ctx.JSON(200, gin.H{
		"access_token": accessToken, "refresh_token": refreshToken,
	})
}

func UserLogin(ctx *gin.Context) {
	var user request.LoginForm

	validationError := ctx.BindJSON(&user)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	findedUser := repo.GetUserForLogin(user.Phone)

	if findedUser.ID == 0 {
		ctx.JSON(400, gin.H{"message": "There is no user with this phone"})
		return
	}

	compareError := bcrypt.CompareHashAndPassword(
		[]byte(findedUser.Password), []byte(user.Password),
	)

	if compareError != nil {
		ctx.JSON(400, gin.H{"message": "Invalid password"})
		return
	}

	accessToken := utils.CreateToken(
		findedUser.ID, config.ENV.ACCESS_TIME,
		config.ENV.ACCESS_KEY, "user",
	)

	refreshToken := utils.CreateToken(
		findedUser.ID, config.ENV.REFRESH_TIME,
		config.ENV.REFRESH_KEY, "user",
	)

	err := repo.SetUserNotificationToken(user.NotificationToken, findedUser.ID)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          findedUser,
	})
}

func UserGetMe(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)

	if role != "user" {
		ctx.AbortWithStatus(403)
		return
	}

	userID := ctx.MustGet("id").(int)

	user := repo.GetUserMe(userID)

	ctx.JSON(200, user)
}

// func UserVerify(ctx *gin.Context) {
// 	var body request.UserVerify

// 	validationError := ctx.BindJSON(&body)

// 	if validationError != nil {
// 		ctx.JSON(400, validationError.Error())
// 		return
// 	}

// 	res, reqErr := http.Get(config.ENV.SMS_API + "/" + body.SmsID)

// 	if reqErr != nil {
// 		log.Println(reqErr.Error())
// 		ctx.JSON(500, gin.H{"message": "couldn't send to verify"})
// 		return
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		ctx.JSON(400, gin.H{"message": "wrong message ID"})
// 		return
// 	}

// 	var resBody response.OTP

// 	decodeErr := json.NewDecoder(res.Body).Decode(&resBody)

// 	if decodeErr != nil {
// 		log.Println(decodeErr.Error())
// 		ctx.JSON(500, gin.H{"message": "Error with decoding response body"})
// 		return
// 	}

// 	// verify message
// 	if body.Password != resBody.Message {
// 		ctx.JSON(400, gin.H{"message": "invalid password"})
// 		return
// 	}

// 	verifyErr := repo.VerifyUser(body.Phone)

// 	if verifyErr != nil {
// 		log.Println(verifyErr.Error())
// 		ctx.JSON(500, verifyErr.Error())
// 		return
// 	}

// 	user := repo.GetUserForLogin(body.Phone)

// 	accessToken := utils.CreateToken(
// 		user.ID, config.ENV.ACCESS_TIME, config.ENV.ACCESS_KEY, "user",
// 	)

// 	refreshToken := utils.CreateToken(
// 		user.ID, config.ENV.REFRESH_TIME, config.ENV.REFRESH_KEY, "user",
// 	)

// 	ctx.JSON(200, gin.H{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 		"user":          user,
// 	})
// }

// func UserForgetPassword(ctx *gin.Context) {
// 	var reqBody request.ForgetPasswordForm

// 	validationError := ctx.BindJSON(&reqBody)

// 	if validationError != nil {
// 		ctx.JSON(400, validationError.Error())
// 		return
// 	}

// 	phone := repo.CheckUserExist(reqBody.Phone)

// 	if phone == "" {
// 		ctx.JSON(400, gin.H{"message": "There is no user with this phone"})
// 		return
// 	}

// 	opt, err := utils.SendOTP(reqBody.Phone)

// 	if err != nil {
// 		log.Println(err.Error())
// 		ctx.JSON(500, err.Error())
// 		return
// 	}

// 	ctx.JSON(200, gin.H{"id": opt.ID})
// }

// func UserNewPassword(ctx *gin.Context) {
// 	var reqBody request.UserNewPassword

// 	validationError := ctx.BindJSON(&reqBody)

// 	if validationError != nil {
// 		ctx.JSON(400, validationError.Error())
// 		return
// 	}

// 	res, reqErr := http.Get(config.ENV.SMS_API + "/" + reqBody.ID)

// 	if reqErr != nil {
// 		log.Println(reqErr.Error())
// 		ctx.JSON(500, gin.H{"message": "couldn't send to verify"})
// 		return
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		ctx.JSON(400, gin.H{"message": "wrong message ID"})
// 		return
// 	}

// 	var resBody response.OTP

// 	decodeErr := json.NewDecoder(res.Body).Decode(&resBody)

// 	if decodeErr != nil {
// 		log.Println(decodeErr.Error())
// 		ctx.JSON(500, gin.H{"message": "Error with decoding response body"})
// 		return
// 	}

// 	if reqBody.OTP != resBody.Message {
// 		ctx.JSON(400, gin.H{"message": "invalid password"})
// 		return
// 	}

// 	hashedPassword, _ := bcrypt.GenerateFromPassword(
// 		[]byte(reqBody.Password), 10,
// 	)

// 	user := repo.UpdateUserPassword(string(hashedPassword), reqBody.Phone)

// 	accessToken := utils.CreateToken(
// 		user.ID, config.ENV.ACCESS_TIME, config.ENV.ACCESS_KEY, "user",
// 	)

// 	refreshToken := utils.CreateToken(
// 		user.ID, config.ENV.REFRESH_TIME, config.ENV.REFRESH_KEY, "user",
// 	)

// 	ctx.JSON(200, gin.H{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 		"user":          user,
// 	})
// }

func WorkerLogin(ctx *gin.Context) {
	var worker request.LoginForm

	validationError := ctx.BindJSON(&worker)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	findedWorker := repo.GetWorkerForLogin(worker.Phone)

	compareError := bcrypt.CompareHashAndPassword(
		[]byte(findedWorker.Password), []byte(worker.Password),
	)

	if compareError != nil {
		ctx.JSON(400, gin.H{"message": "Invalid password"})
		return
	}

	accessToken := utils.CreateToken(
		findedWorker.ID, config.ENV.ACCESS_TIME,
		config.ENV.ACCESS_KEY, "worker",
	)

	refreshToken := utils.CreateToken(
		findedWorker.ID, config.ENV.REFRESH_TIME,
		config.ENV.REFRESH_KEY, "worker",
	)

	ctx.JSON(200, gin.H{
		"access_token": accessToken, "refresh_token": refreshToken,
		"id": findedWorker.ID,
	})
}

func RefreshToken(ctx *gin.Context) {
	var refreshToken request.RefreshTokenForm

	validationError := ctx.BindJSON(&refreshToken)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		refreshToken.RefreshToken, claims, func(
			t *jwt.Token,
		) (interface{}, error) {
			return []byte(config.ENV.REFRESH_KEY), nil
		},
	)

	if err != nil {
		ctx.AbortWithStatus(401)
		return
	}

	prepareToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   int(claims["id"].(float64)),
		"role": claims["role"],
		"exp":  time.Now().Add(config.ENV.ACCESS_TIME).Unix(),
	})

	finalToken, _ := prepareToken.SignedString([]byte(config.ENV.ACCESS_KEY))

	ctx.JSON(200, gin.H{"access_token": finalToken})
}

func RefreshNotificationToken(ctx *gin.Context) {
	var notificationToken request.NotificationToken

	validationError := ctx.BindJSON(&notificationToken)

	if validationError != nil {
		ctx.JSON(400, validationError.Error())
		return
	}

	id := ctx.MustGet("id").(int)
	role := ctx.MustGet("role").(string)

	if role == "admin" || role == "worker" {
		ctx.AbortWithStatus(403)
		return
	}

	repo.SetUserNotificationToken(notificationToken.Token, id)

	ctx.JSON(200, gin.H{"message": "Has been updated"})
}
