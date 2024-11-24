package services

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/huandu/xstrings"
	"log"
	"net/http"
	"strconv"
	"texApi/config"
	"texApi/internal/dto"
	"texApi/internal/repositories"
	"texApi/pkg/smtp"
	"texApi/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

func UserLogin(ctx *gin.Context) {
	credType := ctx.GetHeader("CredType")
	if credType == "" {
		response := utils.FormatErrorResponse("Invalid Login Method", "")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}
	username, password, err := utils.ExtractBasicAuth(ctx.GetHeader("Authorization"))
	if err != nil {
		response := utils.FormatErrorResponse("Unauthorized", err.Error())
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	user, err := repositories.GetUser(username, credType)
	if err != nil {
		response := utils.FormatErrorResponse("User not found", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if config.ENV.ENCRYPT_PASSWORDS {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			response := utils.FormatErrorResponse("Login failed", "Invalid credentials")
			ctx.JSON(http.StatusUnauthorized, response)
			return
		}
	} else {
		if user.Password != password {
			response := utils.FormatErrorResponse("Login failed", "Invalid credentials")
			ctx.JSON(http.StatusUnauthorized, response)
			return
		}
	}

	accessToken, refreshToken, exp := utils.CreateToken(user.ID, user.RoleID, user.CompanyID, user.Role)
	err = repositories.ManageToken(user.ID, refreshToken, "create")
	if err != nil {
		response := utils.FormatErrorResponse("Error creating token", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.FormatResponse("Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"exp":           exp,
		"user":          user,
	})
	ctx.JSON(http.StatusOK, response)
}

func UserGetMe(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)

	if userID == 0 {
		response := utils.FormatErrorResponse("Unauthorized", "")
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	user := repositories.GetUserById(userID)
	if user.ID == 0 {
		response := utils.FormatErrorResponse("User not found", "")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	response := utils.FormatResponse("User retrieved successfully", user)
	ctx.JSON(http.StatusOK, response)
}

func Logout(ctx *gin.Context) {
	id := ctx.MustGet("id").(int)
	_, err := repositories.RemoveUserToken(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Couldn't logout", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.FormatResponse("Logged out successfully", id))
}

func RefreshToken(ctx *gin.Context) {
	var payload dto.RefreshTokenForm

	validationError := ctx.BindJSON(&payload)
	if validationError != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid body", validationError.Error()))
		return
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		payload.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.ENV.REFRESH_KEY), nil
		},
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Error occurred", err.Error()))
		return
	}

	idFloat, _ := claims["id"].(float64)
	id := int(idFloat)
	user := repositories.GetUserById(id)

	err = repositories.ManageToken(user.ID, payload.RefreshToken, "validate")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("User/Token is invalid", err.Error()))
		return
	}

	accessToken, _, exp := utils.CreateToken(user.ID, user.RoleID, user.CompanyID, user.Role)
	response := utils.FormatResponse("Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": payload.RefreshToken,
		"exp":           exp,
		"user":          user,
	})
	ctx.JSON(http.StatusCreated, response)
}

func ForgotPassword(ctx *gin.Context) {
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	user, _ := repositories.GetUser(credentials, credType)
	if user.ID == 0 {
		response := utils.FormatErrorResponse(fmt.Sprintf("A user with this %s not found", credType), "")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	otp, _ := utils.GenerateOTP()

	_, err := repositories.SaveUserWithOTP(user.ID, user.RoleID, user.Verified, credType, credentials, otp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating token", err.Error()))
		return
	}

	if credType == "email" {
		err := smtp.SendOTPEmail(credentials, otp)
		if err != nil {
			log.Fatalf("Error sending email: %v", err)
		}
	}

	//// TODO: change this, it's development mode:
	ctx.JSON(http.StatusOK, utils.FormatResponse(otp, ""))
	return
}

func UpdatePasswordOTP(ctx *gin.Context) {
	promptOTP := ctx.GetHeader("OTP")
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	currentUser, err := repositories.GetUser(credentials, credType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("User Not Verified!", err.Error()))
		return
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05", currentUser.VerifyTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error parsing time", err.Error()))
		return
	}
	expirationTime := parsedTime.Add(15 * time.Minute)
	if time.Now().After(expirationTime) || promptOTP != currentUser.OTPKey {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Register time expired or wrong token!", ""))
		return
	}

	userData := dto.ProfileUpdate{}

	password := ctx.GetHeader("Password")
	if password != "" {
		userData.Password = &password
	}

	_, err = repositories.ProfileUpdate(userData, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Couldn't update profile", err.Error()))
		return
	}

	ctx.JSON(http.StatusBadRequest, utils.FormatResponse("Password successfully updated", password))
	return
}

func RegisterRequest(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.GetHeader("RoleID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Role ID is required", ""))
		return
	}
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")
	if roleID < 3 || credentials == "" || credType == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid Request, missing required params", ""))
		return
	}

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	user, _ := repositories.GetUser(credentials, credType)
	if user.ID > 0 {
		if user.Verified == 1 && user.Deleted == 0 {
			response := utils.FormatErrorResponse(fmt.Sprintf("A user with this %s already exists", credType), "")
			ctx.JSON(http.StatusOK, response)
			return
		}
	}

	otp, _ := utils.GenerateOTP()

	_, err = repositories.SaveUserWithOTP(user.ID, roleID, 0, credType, credentials, otp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating token", err.Error()))
		return
	}

	if credType == "email" {
		err := smtp.SendOTPEmail(credentials, otp)
		if err != nil {
			log.Fatalf("Error sending email: %v", err)
		}
	}

	//// TODO: change this, it's development mode:
	ctx.JSON(http.StatusOK, utils.FormatResponse(otp, ""))
	return
}

func ValidateOTP(ctx *gin.Context) {
	promptOTP := ctx.GetHeader("OTP")
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")
	if err := repositories.ValidateOTPAndTime(credType, credentials, promptOTP); err != nil {
		response := utils.FormatErrorResponse(err.Error(), "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.FormatResponse("OTP check success", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func Register(ctx *gin.Context) {
	promptOTP := ctx.GetHeader("OTP")
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	var user dto.CreateUser
	validationError := ctx.BindJSON(&user)
	if validationError != nil {
		response := utils.FormatErrorResponse("Invalid request body", validationError.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser, err := repositories.GetUser(credentials, credType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Unable to create user!", err.Error()))
		return
	}

	// verifying that request time is valid
	parsedTime, err := time.Parse("2006-01-02 15:04:05", currentUser.VerifyTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error parsing time", err.Error()))
		return
	}
	expirationTime := parsedTime.Add(15 * time.Minute)
	if time.Now().After(expirationTime) || promptOTP != currentUser.OTPKey {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Register time expired or wrong token!", ""))
		return
	}
	if currentUser.Verified == 0 {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("User Not Verified!", ""))
		return
	}
	user.Verified = currentUser.Verified
	user.Active = currentUser.Active
	user.RoleID = currentUser.RoleID
	if currentUser.RoleID > 3 {
		user.Role = "carrier"
	} else if currentUser.RoleID == 3 {
		user.Role = "sender"
	}

	if credType == "email" {
		user.Email = currentUser.Email
		user.Phone = ""
	} else {
		user.Phone = currentUser.Phone
		user.Email = ""
	}
	userID, err := repositories.UpdateUser(user, currentUser.ID)
	if err != nil {
		response := utils.FormatErrorResponse("Cannot register user", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	accessToken, refreshToken, exp := utils.CreateToken(userID, user.RoleID, user.CompanyID, user.Role)
	err = repositories.ManageToken(userID, refreshToken, "create")
	if err != nil {
		response := utils.FormatErrorResponse("User created, but found error creating token, try logging in now", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("User created successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"exp":           exp,
	})
	ctx.JSON(http.StatusOK, response)
	return
}

func ProfileUpdate(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", ""))
		return
	}

	var userData dto.ProfileUpdate
	validationError := ctx.BindJSON(&userData)
	if validationError != nil {
		response := utils.FormatErrorResponse("Invalid request body", validationError.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, err := repositories.ProfileUpdate(userData, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Couldn't update profile", err.Error()))
		return
	}

	user := repositories.GetUserById(userID)
	if user.ID == 0 {
		response := utils.FormatErrorResponse("User not found", "")
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Profile successfully updated", user))
	return
}

func BeginOAuth(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func CompleteOAuth(ctx *gin.Context) {
	roleID, _ := strconv.Atoi(ctx.GetHeader("RoleID"))
	role := ctx.GetHeader("Role")
	fmt.Println(role, roleID)

	if roleID < 3 {
		roleID = 3
	}
	if roleID > 3 {
		role = "carrier"
	} else if roleID == 3 {
		role = "sender"
	}

	authUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", err.Error()))
		return
	}

	var userID int
	var user dto.CreateUser
	dbUser, _ := repositories.GetUser(user.Email, "email")
	if dbUser.ID == 0 {
		user.Role = role
		user.RoleID = roleID
		user.Email = authUser.Email
		user.Username = fmt.Sprintf("%s%s%s%s", authUser.Name, authUser.FirstName, authUser.LastName, authUser.UserID)
		user.Password = xstrings.Shuffle(fmt.Sprintf("%s%s", authUser.UserID, authUser.Email))
		user.OauthIDToken = authUser.UserID
		user.OauthProvider = authUser.Provider
		user.Verified = 1
		user.Active = 1
		user.Phone = ""

		userID, err = repositories.CreateUser(user)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Error adding user", err.Error()))
			return
		}
	} else {
		userID = dbUser.ID
	}

	accessToken, refreshToken, exp := utils.CreateToken(userID, user.RoleID, user.CompanyID, user.Role)
	err = repositories.ManageToken(userID, refreshToken, "create")
	if err != nil {
		response := utils.FormatErrorResponse("User created, but found error creating token, try logging in now", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("User created successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"exp":           exp,
	})
	ctx.JSON(http.StatusOK, response)
	return
}
