package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"texApi/config"
	"texApi/internal/dto"
	"texApi/internal/repo"
	"texApi/pkg/smtp"
	"texApi/pkg/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/huandu/xstrings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func sendNewLoginNotification(userID int, content string, data interface{}) {
	payload := map[string]interface{}{
		"userID":  userID,
		"content": content,
		"extras":  data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal notification payload: %s", err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:%s/%s/ws-notification/", config.ENV.API_PORT, config.ENV.API_PREFIX),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Failed to create notification request: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(config.ENV.SYSTEM_HEADER, config.ENV.API_SECRET)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send notification request %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf(fmt.Sprintf("Notification API returned non-OK status: %d", resp.StatusCode), nil)
	}
}

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

	user, err := repo.GetUser(username, credType)
	if err != nil {
		response := utils.FormatErrorResponse("User not found", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if config.ENV.ENCRYPT_PASSWORDS {
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Login failed, Invalid credentials", err.Error()))
			return
		}
	} else {
		if user.Password != password {
			ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Login failed", "Invalid credentials"))
			return
		}
	}

	if user.Role == "driver" {
		if CheckDriverNotBlocked(ctx, user.DriverID) == false {
			return
		}
	}

	accessToken, refreshToken, accessExp := utils.CreateToken(user.ID, user.RoleID, user.CompanyID, user.DriverID, user.Role)
	deviceName, deviceModel, deviceFirmware, appName, appVersion := ExtractDeviceInfo(ctx)

	refreshExp := time.Now().Add(config.ENV.REFRESH_TIME)
	session := dto.CreateSessionInput{
		UserID:         user.ID,
		CompanyID:      user.CompanyID,
		RefreshToken:   refreshToken,
		ExpiresAt:      refreshExp,
		DeviceName:     deviceName,
		DeviceModel:    deviceModel,
		DeviceFirmware: deviceFirmware,
		AppName:        appName,
		AppVersion:     appVersion,
		UserAgent:      ctx.GetHeader("User-Agent"),
		IPAddress:      ctx.ClientIP(),
		LoginMethod:    "password",
	}

	_, err = repo.CreateSession(session)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating session", "Please try again"))
		return
	}

	go sendNewLoginNotification(user.ID, fmt.Sprintf("%s, your account has been logged in from a new device", user.Username), session)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    accessExp,
		"user":          user,
	}))
}

func UserGetMe(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", ""))
		return
	}

	user, err := repo.GetUserById(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("User retrieved successfully", user))
}

func RefreshToken(ctx *gin.Context) {
	var payload dto.RefreshTokenForm
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid body", err.Error()))
		return
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		payload.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.ENV.REFRESH_KEY), nil
		},
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Invalid token", "The refresh token is invalid or expired"))
		return
	}

	session, err := repo.GetSessionByRefreshToken(payload.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Invalid session", "Session not found or expired"))
		return
	}
	if time.Now().Add(config.ENV.TZAddHours).After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Session expired", ""))
		return
	}

	user, err := repo.GetUserById(session.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("User not found", "Cannot refresh token for invalid user"))
		return
	}

	if user.Role == "driver" {
		if CheckDriverNotBlocked(ctx, user.DriverID) == false {
			return
		}
	}

	accessToken, refreshToken, accessExp := utils.CreateToken(user.ID, user.RoleID, user.CompanyID, user.DriverID, user.Role)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Token refreshed", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    accessExp,
		"user":          user,
	}))
}

func Logout(ctx *gin.Context) {
	var payload dto.RefreshTokenForm
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request", "Missing refresh token"))
		return
	}

	err := repo.InvalidateSession(payload.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Logout failed", "Error invalidating session"))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Logged out successfully", nil))
}

func LogoutAllSessions(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)

	err := repo.InvalidateAllUserSessions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Logout failed", "Error invalidating sessions"))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("All sessions terminated successfully", nil))
}

func ForgotPassword(ctx *gin.Context) {
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	user, err := repo.GetUser(credentials, credType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse(fmt.Sprintf("A user with this %s not found", credType), err.Error()))
		return
	}

	otp := utils.GenerateOTP()
	userID, err := repo.SaveUserWithOTP(user.ID, user.RoleID, user.Verified, credType, credentials, otp, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating token", err.Error()))
		return
	}

	if credType == "email" {
		if err = smtp.SendOTPEmail(credentials, otp); err != nil {
			fmt.Printf("Error sending email: %v\n", err)
		}
	} else if credType == "phone" {
		if err = utils.SendOTPSMS(credentials, otp, utils.DetectDeviceFirmware(ctx.GetHeader("X-Device-Firmware"))); err != nil {
			log.Printf("Error sending SMS: %v\n", err)
		}
	}

	if config.ENV.API_DEBUG {
		ctx.JSON(http.StatusOK, utils.FormatResponse("OTP sent successfully", gin.H{
			"otp":     otp,
			"user_id": userID,
		}))
	} else {
		ctx.JSON(http.StatusOK, utils.FormatResponse("OTP sent successfully", gin.H{
			"user_id": userID,
		}))
	}
}

func UpdatePasswordOTP(ctx *gin.Context) {
	promptOTP := ctx.GetHeader("OTP")
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	user, err := repo.GetUser(credentials, credType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("User Not Verified!", err.Error()))
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", user.VerifyTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error parsing time", err.Error()))
		return
	}

	if time.Now().Add(config.ENV.TZAddHours).After(parsedTime.Add(15*time.Minute)) || promptOTP != user.OTPKey {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Register time expired or wrong token!", ""))
		return
	}

	password := ctx.GetHeader("Password")
	if password == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Password is required", ""))
		return
	}

	_, err = repo.UserUpdate(dto.UserUpdateAuth{Password: &password}, user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Couldn't update company", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Password successfully updated", nil))
}

func RegisterRequest(ctx *gin.Context) {
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")
	role := ctx.GetHeader("Role")
	roleID, err := strconv.Atoi(ctx.GetHeader("RoleID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Role ID is required", ""))
		return
	}
	if !(role == "sender" || role == "carrier") || roleID < 3 || credentials == "" || credType == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid Request, invalid or missing required header params: Role, Credentials, CredType", ""))
		return
	}

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	user, _ := repo.GetUser(credentials, credType)
	if user.ID > 0 && user.Verified == 1 && user.Deleted == 0 {
		ctx.JSON(http.StatusOK, utils.FormatErrorResponse(fmt.Sprintf("A user with this %s already exists", credType), ""))
		return
	}

	otp := utils.GenerateOTP()
	userID, err := repo.SaveUserWithOTP(user.ID, roleID, 0, credType, credentials, otp, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating token", err.Error()))
		return
	}

	if credType == "email" {
		if err = smtp.SendOTPEmail(credentials, otp); err != nil {
			fmt.Printf("Error sending email: %v\n", err)
		}
	} else if credType == "phone" {
		if err = utils.SendOTPSMS(credentials, otp, utils.DetectDeviceFirmware(ctx.GetHeader("X-Device-Firmware"))); err != nil {
			log.Printf("Error sending SMS: %v\n", err)
		}

	}

	if config.ENV.API_DEBUG {
		ctx.JSON(http.StatusOK, utils.FormatResponse("OTP sent successfully", gin.H{
			"otp":     otp,
			"user_id": userID,
		}))
	} else {
		ctx.JSON(http.StatusOK, utils.FormatResponse("OTP sent successfully", gin.H{
			"user_id": userID,
		}))
	}
}

func ValidateOTP(ctx *gin.Context) {
	promptOTP := ctx.GetHeader("OTP")
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	err := repo.ValidateOTPAndTime(credType, credentials, promptOTP)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("OTP check success", nil))
}

func Register(ctx *gin.Context) {
	promptOTP := ctx.GetHeader("OTP")
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	var user dto.CreateUser
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	currentUser, err := repo.GetUser(credentials, credType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Unable to create user!", err.Error()))
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", currentUser.VerifyTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error parsing time", err.Error()))
		return
	}

	if time.Now().Add(config.ENV.TZAddHours).After(parsedTime.Add(15*time.Minute)) || promptOTP != currentUser.OTPKey {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Register time expired or wrong token!", ""))
		return
	}

	if currentUser.Verified == 0 {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("User Not Verified!", ""))
		return
	}

	// TODO: check, this is unused
	//user.Verified = currentUser.Verified
	//user.Active = currentUser.Active
	//user.RoleID = currentUser.RoleID
	if currentUser.RoleID > 3 {
		user.Role = "carrier"
	} else if currentUser.RoleID == 3 {
		user.Role = "sender"
	}

	if credType == "email" {
		user.Email = currentUser.Email
		user.Phone = ""
	} else if credType == "phone" {
		user.Phone = currentUser.Phone
		user.Email = ""
	} else {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request credType", ""))
		return
	}

	userID, err := repo.UpdateUser(user, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Cannot register user", err.Error()))
		return
	}

	accessToken, refreshToken, accessExp := utils.CreateToken(userID, user.RoleID, user.CompanyID, user.DriverID, user.Role)
	err = repo.ManageToken(userID, refreshToken, "create")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("User created, but found error creating token, try logging in now", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("User created successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    accessExp,
	}))
}

func UserUpdate(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", ""))
		return
	}

	var userData dto.UserUpdateAuth
	if err := ctx.BindJSON(&userData); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	_, err := repo.UserUpdate(userData, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Couldn't update company", err.Error()))
		return
	}

	user, err := repo.GetUserById(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Company successfully updated", user))
}
func OTPLoginRequest(ctx *gin.Context) {
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	roleID, err := strconv.Atoi(ctx.GetHeader("RoleID"))
	if err != nil {
		roleID = 3
	}
	role := ctx.GetHeader("Role")
	if !(role == "sender" || role == "carrier" || role == "driver") || roleID < 3 || credentials == "" || credType == "" {
		role = "sender"
		roleID = 3
	}

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	otp := utils.GenerateOTP()
	if len(otp) == 0 {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error generating OTP", ""))
		return
	}

	dbUser, err := repo.GetUser(credentials, credType)
	var userID int
	if err == nil {
		userID, err = repo.SaveUserWithOTP(dbUser.ID, dbUser.RoleID, dbUser.Verified, credType, credentials, otp, dbUser.Role)
	} else {
		newUser := dto.CreateUser{
			Password: xstrings.Shuffle(fmt.Sprintf("%s%s", utils.GenerateOTP(6), credentials)),
			Role:     role,
			RoleID:   roleID,
			Active:   1,
			OTP:      &otp,
		}
		if credType == "email" {
			newUser.Username = fmt.Sprintf("%s_%s", strings.Split(credentials, "@")[0], utils.GenerateOTP(6))
			newUser.Email = credentials
			newUser.Phone = ""
		}
		if credType == "phone" {
			newUser.Username = fmt.Sprintf("user_%s_%s", credentials, utils.GenerateOTP(6))
			newUser.Email = ""
			newUser.Phone = credentials
		}

		userID, err = repo.CreateUser(newUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating user", err.Error()))
			return
		}

		newCompany := dto.CompanyCreateShort{
			UserID:      userID,
			Role:        role,
			RoleID:      roleID,
			FirstName:   newUser.Username,
			LastName:    "",
			CompanyName: newUser.Username,
			Email:       newUser.Email,
			ImageURL:    "",
		}
		companyID, err := repo.CreateCompanyShort(newCompany)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating company", err.Error()))
		}
		repo.UpdateUserCompanyID(userID, companyID)
		dbUser, _ = repo.GetUser(credentials, credType)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating login token", err.Error()))
		return
	}

	if credType == "email" {
		if err = smtp.SendOTPEmail(credentials, otp); err != nil {
			fmt.Printf("Error sending email: %v\n", err)
		}
	} else if credType == "phone" {
		if err = utils.SendOTPSMS(credentials, otp, utils.DetectDeviceFirmware(ctx.GetHeader("X-Device-Firmware"))); err != nil {
			log.Printf("Error sending SMS: %v\n", err)
		}
	}

	if config.ENV.API_DEBUG {
		ctx.JSON(http.StatusOK, utils.FormatResponse("OTP sent successfully", gin.H{
			"otp":     otp,
			"user_id": userID,
		}))
	} else {
		ctx.JSON(http.StatusOK, utils.FormatResponse("OTP sent successfully", gin.H{
			"user_id": userID,
		}))
	}
}
func OTPLogin(ctx *gin.Context) {
	credentials := ctx.GetHeader("Credentials")
	credType := ctx.GetHeader("CredType")
	promptOTP := ctx.GetHeader("OTP")

	if ok, msg := utils.ValidateCredential(credType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	if promptOTP == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("OTP is required", ""))
		return
	}

	user, err := repo.GetUser(credentials, credType)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found", ""))
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", user.VerifyTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error validating OTP", err.Error()))
	}

	if time.Now().Add(config.ENV.TZAddHours).After(parsedTime.Add(15*time.Minute)) || promptOTP != user.OTPKey {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid or expired OTP", ""))
		return
	}

	if user.Role == "driver" {
		if CheckDriverNotBlocked(ctx, user.DriverID) == false {
			return
		}
	}

	accessToken, refreshToken, accessExp := utils.CreateToken(user.ID, user.RoleID, user.CompanyID, user.DriverID, user.Role)
	deviceName, deviceModel, deviceFirmware, appName, appVersion := ExtractDeviceInfo(ctx)

	refreshExp := time.Now().Add(config.ENV.REFRESH_TIME)
	session := dto.CreateSessionInput{
		UserID:         user.ID,
		CompanyID:      user.CompanyID,
		RefreshToken:   refreshToken,
		ExpiresAt:      refreshExp,
		DeviceName:     deviceName,
		DeviceModel:    deviceModel,
		DeviceFirmware: deviceFirmware,
		AppName:        appName,
		AppVersion:     appVersion,
		UserAgent:      ctx.GetHeader("User-Agent"),
		IPAddress:      ctx.ClientIP(),
		LoginMethod:    "otp",
	}

	_, err = repo.CreateSession(session)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating session", "Please try again"))
		return
	}

	go sendNewLoginNotification(user.ID, fmt.Sprintf("%s, your account has been logged in from a new device", user.Username), session)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    accessExp,
		"user":          user,
	}))
}

func BeginOAuth(ctx *gin.Context) {
	state, err := GenerateStateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to generate state token", err.Error()))
		return
	}

	roleID, err := strconv.Atoi(ctx.GetHeader("RoleID"))
	if err != nil {
		roleID = 3
	}
	role := ctx.GetHeader("Role")
	if !(role == "sender" || role == "carrier" || role == "driver") || roleID < 3 {
		role = "sender"
		roleID = 3
		//ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid Request, invalid or missing required header params: Role, Credentials, CredType", ""))
		//return
	}

	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Set("role", role)
	session.Set("roleID", roleID)
	err = session.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to save session", err.Error()))
		return
	}

	authUrl := config.ENV.GoogleOAuthConfig.AuthCodeURL(state)

	ctx.Redirect(http.StatusTemporaryRedirect, authUrl)
	return
}

func BeginOAuthMobile(ctx *gin.Context) {
	var request struct {
		IDToken string `json:"id_token"`
		Role    string `json:"role"`
		RoleID  int    `json:"role_id"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request", err.Error()))
		return
	}

	userInfo, err := VerifyGoogleIDToken(request.IDToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Invalid ID token", err.Error()))
		return
	}

	AuthenticateOAuthUser(ctx, userInfo, request.Role, request.RoleID)
}

func CompleteOAuth(ctx *gin.Context) {
	// For CSRF
	session := sessions.Default(ctx)
	state := session.Get("state")
	if state != ctx.Query("state") {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid state parameter", "Session states mismatch"))
		return
	}

	role, ok := session.Get("role").(string)
	if !ok || !(role == "sender" || role == "carrier" || role == "driver" || role == "admin") {
		role = "sender"
	}

	roleID, ok := session.Get("roleID").(int)
	if !ok {
		roleID = 3
	}

	code := ctx.Query("code")
	token, err := config.ENV.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Code exchange failed", err.Error()))
		return
	}

	userInfo, err := GetUserInfo(token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Failed to get user info", err.Error()))
		return
	}

	AuthenticateOAuthUser(ctx, userInfo, role, roleID)
}

func AuthenticateOAuthUser(ctx *gin.Context, userInfo map[string]interface{}, role string, roleID int) {
	authUser := GetOAuthUserFromInfo(userInfo)
	DbUser, _ := repo.GetUser(authUser.Email, "email")
	var userID int
	var companyID int
	var driverID int
	if DbUser.ID == 0 {
		newUser := dto.CreateUser{
			Username: fmt.Sprintf("%s%s", strings.Split(authUser.Email, "@")[0], utils.GenerateOTP(6)),
			Password: xstrings.Shuffle(fmt.Sprintf("%s%s", utils.GenerateOTP(6), authUser.Email)),
			Email:    authUser.Email,
			Role:     role,
			RoleID:   roleID,
			Active:   1,
			Meta:     fmt.Sprint(userInfo),
		}

		userID, err := repo.CreateUser(newUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating user", err.Error()))
			return
		}

		newCompany := dto.CompanyCreateShort{
			UserID:      userID,
			Role:        role,
			RoleID:      roleID,
			FirstName:   authUser.FirstName,
			LastName:    authUser.LastName,
			CompanyName: authUser.Name,
			Email:       authUser.Email,
			About:       "",
			Phone:       "",
			ImageURL:    authUser.AvatarURL,
		}
		companyID, err = repo.CreateCompanyShort(newCompany)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating company", err.Error()))
			return
		}
		repo.UpdateUserCompanyID(userID, companyID)
		DbUser, _ = repo.GetUser(authUser.Email, "email")
	} else {
		userID = DbUser.ID
		companyID = DbUser.CompanyID
		driverID = DbUser.DriverID
	}

	accessToken, refreshToken, accessExp := utils.CreateToken(userID, roleID, companyID, driverID, DbUser.Role)
	deviceName, deviceModel, deviceFirmware, appName, appVersion := ExtractDeviceInfo(ctx)

	refreshExp := time.Now().Add(config.ENV.REFRESH_TIME)
	dbSession := dto.CreateSessionInput{
		UserID:         DbUser.ID,
		CompanyID:      DbUser.CompanyID,
		RefreshToken:   refreshToken,
		ExpiresAt:      refreshExp,
		DeviceName:     deviceName,
		DeviceModel:    deviceModel,
		DeviceFirmware: deviceFirmware,
		AppName:        appName,
		AppVersion:     appVersion,
		UserAgent:      ctx.GetHeader("User-Agent"),
		IPAddress:      ctx.ClientIP(),
		LoginMethod:    "oauth_" + authUser.Provider,
	}

	_, err := repo.CreateSession(dbSession)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating session", err.Error()))
		return
	}

	go sendNewLoginNotification(userID, fmt.Sprintf("%s, your account has been logged in from a new device", DbUser.Username), dbSession)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    accessExp,
		"user":          DbUser,
	}))
}

func GetOAuthUserFromInfo(userInfo map[string]interface{}) dto.OAuthUser {
	email, _ := userInfo["email"].(string)
	firstName, _ := userInfo["given_name"].(string)
	lastName, _ := userInfo["family_name"].(string)
	name, _ := userInfo["name"].(string)
	avatarURL, _ := userInfo["picture"].(string)

	authUser := dto.OAuthUser{
		RawData:   userInfo,
		Provider:  "Google",
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Name:      name,
		AvatarURL: avatarURL,
	}
	return authUser
}

func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err = json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func VerifyGoogleIDToken(idToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to verify token: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenInfo map[string]interface{}
	if err = json.Unmarshal(data, &tokenInfo); err != nil {
		return nil, err
	}

	if tokenInfo["aud"] != config.ENV.GoogleOAuthConfig.ClientID {
		return nil, fmt.Errorf("token was not issued by this application")
	}

	return tokenInfo, nil
}
