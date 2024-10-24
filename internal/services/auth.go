package services

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"texApi/config"
	"texApi/internal/dto"
	"texApi/internal/repositories"
	"texApi/pkg/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

func UserLogin(ctx *gin.Context) {
	loginType := ctx.GetHeader("LoginType")
	if loginType == "" {
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

	user, err := repositories.GetUser(username, loginType)
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

	accessToken := utils.CreateToken(user.ID, user.RoleID)
	refreshToken := utils.CreateToken(user.ID, user.RoleID)
	err = repositories.ManageToken(user.ID, refreshToken, "create")
	if err != nil {
		log.Println(err.Error())
		response := utils.FormatErrorResponse("Error creating token", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.FormatResponse("Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
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
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	response := utils.FormatResponse("Logged out successfully", nil)
	ctx.JSON(http.StatusOK, response)
}

func RefreshToken(ctx *gin.Context) {
	var refreshToken dto.RefreshTokenForm

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
		"id":     int(claims["id"].(float64)),
		"roleID": claims["roleID"].(int),
		"exp":    time.Now().Add(config.ENV.ACCESS_TIME).Unix(),
	})
	finalToken, _ := prepareToken.SignedString([]byte(config.ENV.ACCESS_KEY))

	response := utils.FormatResponse("Token updated", gin.H{
		"access_token":  finalToken,
		"refresh_token": refreshToken,
	})
	ctx.JSON(http.StatusOK, response)
}

func RegisterRequest(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.GetHeader("RoleID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Role ID is required", ""))
		return
	}
	credentials := ctx.GetHeader("Credentials")
	registerType := ctx.GetHeader("RegisterType")
	if roleID == 0 || credentials == "" || registerType == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid Request, missing required params", ""))
		return
	}

	if ok, msg := utils.ValidateCredential(registerType, credentials); !ok {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse(msg, ""))
		return
	}

	user, _ := repositories.GetUser(credentials, registerType)
	if user.ID > 0 {
		if user.Verified == 1 && user.Deleted == 0 {
			response := utils.FormatErrorResponse(fmt.Sprintf("A user with this %s already exists", registerType), "")
			ctx.JSON(http.StatusOK, response)
			return
		}
	}

	otp, _ := utils.GenerateOTP()

	_, err = repositories.SaveUserWithOTP(user.ID, roleID, registerType, credentials, otp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating token", err.Error()))
		return
	}

	//// TODO: change this, it's development mode:
	ctx.JSON(http.StatusOK, utils.FormatResponse(otp, ""))
	return
}

func ValidateOTP(ctx *gin.Context) {
	promptOTP := ctx.GetHeader("OTP")
	credentials := ctx.GetHeader("Credentials")
	registerType := ctx.GetHeader("RegisterType")
	if err := repositories.ValidateOTPAndTime(registerType, credentials, promptOTP); err != nil {
		response := utils.FormatErrorResponse(err.Error(), "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.FormatResponse("OTP check success", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func Register(ctx *gin.Context) {
	credentials := ctx.GetHeader("Credentials")
	registerType := ctx.GetHeader("RegisterType")

	var user dto.CreateUser
	validationError := ctx.BindJSON(&user)
	if validationError != nil {
		response := utils.FormatErrorResponse("Invalid request body", validationError.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser, err := repositories.GetUser(credentials, registerType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("User Not Verified!", err.Error()))
		return
	}
	// verifying that request time is valid
	parsedTime, err := time.Parse("2006-01-02 15:04:05", currentUser.VerifyTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error parsing time", err.Error()))
		return
	}
	expirationTime := parsedTime.Add(15 * time.Minute)
	if time.Now().After(expirationTime) {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Register time expired. Start over!", ""))
		return
	}
	if currentUser.Verified == 0 {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("User Not Verified!", ""))
		return
	}
	user.Verified = currentUser.Verified
	user.Active = currentUser.Active
	user.RoleID = currentUser.RoleID
	if registerType == "email" {
		user.Email = currentUser.Email
	} else {
		user.Phone = currentUser.Phone
	}
	userID, err := repositories.UpdateUser(user, currentUser.ID)
	if err != nil {
		response := utils.FormatErrorResponse("Cannot register user", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	accessToken := utils.CreateToken(userID, user.RoleID)
	refreshToken := utils.CreateToken(userID, user.RoleID)
	err = repositories.ManageToken(userID, refreshToken, "create")
	if err != nil {
		log.Println(err.Error())
		response := utils.FormatErrorResponse("User created, but found error creating token, try logging in now", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("User created successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	ctx.JSON(http.StatusOK, response)
	return
}

func GetOAuthCallbackFunction(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Set("provider", provider)
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		response := utils.FormatErrorResponse("Unauthorized", err.Error())
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	dbuser, err := repositories.GetUser(user.Email, "email")
	if err != nil {
		response := utils.FormatErrorResponse("Unauthorized", err.Error())
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	if dbuser.ID == 0 {
		// Handle user registration here if needed
	} else {
		session := sessions.Default(ctx)
		session.Set("userID", dbuser.ID)
		session.Set("role", dbuser.RoleID)
		session.Save()
	}

	http.Redirect(ctx.Writer, ctx.Request, "http://localhost:7000/texapp/auth/oauth/testfront/", http.StatusFound)
}

func OAuthLogout(ctx *gin.Context) {
	gothic.Logout(ctx.Writer, ctx.Request)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func OAuthProvider(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
	res := ctx.Writer
	req := ctx.Request

	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		tmpl, err := template.New("user").Parse(userTemplate)
		if err != nil {
			log.Println("Error parsing template:", err)
			ctx.String(http.StatusInternalServerError, "Template parsing error")
			return
		}
		tmpl.Execute(res, gothUser)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func OAuthFront(ctx *gin.Context) {
	providers := []string{"google", "github", "facebook"}
	providersMap := map[string]string{
		"google":   "Google",
		"github":   "GitHub",
		"facebook": "Facebook",
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		log.Println("Error parsing template:", err)
		ctx.String(http.StatusInternalServerError, "Template parsing error")
		return
	}
	tmpl.Execute(ctx.Writer, gin.H{
		"Providers":    providers,
		"ProvidersMap": providersMap,
	})
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/texapp/auth/oauth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/texapp/auth/oauth/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
