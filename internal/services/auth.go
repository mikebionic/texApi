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
	"texApi/internal/_others/schemas/request"
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
	loginMethod := ctx.GetHeader("LoginMethod")
	if loginMethod == "" {
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

	user, err := repositories.GetUser(username, loginMethod)
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
	registerMethod := ctx.GetHeader("RegisterMethod")
	if roleID == 0 || credentials == "" || registerMethod == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid Request, missing required params", ""))
		return
	}

	if !(registerMethod == "email" || registerMethod == "phone") {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Wrong method, use 'email' or 'phone'", ""))
		return
	}

	user, _ := repositories.GetUser(credentials, registerMethod)
	if user.ID > 0 {
		response := utils.FormatErrorResponse(fmt.Sprintf("A user with this %s already exists", registerMethod), "")
		ctx.JSON(http.StatusOK, response)
		return
	}

	// generate OTP
	otp := "1234"
	session := sessions.Default(ctx)
	session.Set("RequestRoleID", roleID)
	session.Set("Credentials", credentials)
	session.Set("RegisterMethod", registerMethod)
	session.Set("otp", otp)
	session.Set("otpValidated", 0)
	session.Save()
	// send email with generated password OTP
}

func ValidateOTP(ctx *gin.Context) {
	session := sessions.Default(ctx)
	otp := session.Get("otp")
	promptOTP := ctx.GetHeader("OTP")
	if otp != promptOTP {
		response := utils.FormatErrorResponse("Wrong OTP", "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	session.Set("otpValidated", 1)
	session.Save()
	response := utils.FormatResponse("OTP check success", nil)
	ctx.JSON(http.StatusOK, response)
	return
}

func Register(ctx *gin.Context) {
	session := sessions.Default(ctx)
	otpValidated, ok := session.Get("otpValidated").(int)
	if otpValidated == 0 || !ok {
		response := utils.FormatErrorResponse("OTP error", "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var user dto.CreateUser
	validationError := ctx.BindJSON(&user)
	if validationError != nil {
		response := utils.FormatErrorResponse("Invalid request body", validationError.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	//Credentials := session.Get("Credentials")
	//RegisterMethod := session.Get("RegisterMethod")
	// just in case for harder logics
	user.RoleID = session.Get("RequestRoleID").(int)
	user.Verified = 1
	user.Active = 1
	fmt.Println(user)
	userID, err := repositories.CreateUser(user)
	if userID == 0 || err != nil {
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

	session.Clear()
	session.Save()
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
