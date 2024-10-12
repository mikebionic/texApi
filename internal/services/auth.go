package services

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"texApi/config"
	"texApi/internal/repositories"
	"texApi/pkg/utils"

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

	accessToken := utils.CreateToken(user.ID, config.ENV.ACCESS_TIME, config.ENV.ACCESS_KEY, "user.RoleId")
	refreshToken := utils.CreateToken(user.ID, config.ENV.REFRESH_TIME, config.ENV.REFRESH_KEY, "user.RoleId")
	//err := repositories.ManageToken(user.ID, refreshToken)
	//
	//if err != nil {
	//	log.Println(err.Error())
	//	response := utils.FormatErrorResponse("Error creating token", err.Error())
	//	ctx.JSON(http.StatusInternalServerError, response)
	//	return
	//}

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

func RegisterRequest(ctx *gin.Context) {
	roleID := ctx.GetHeader("RoleID")
	credentials := ctx.GetHeader("Credentials")
	registerMethod := ctx.GetHeader("RegisterMethod")
	if roleID == "" || credentials == "" || registerMethod == "" {
		response := utils.FormatErrorResponse("Invalid Request, missing required params", "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if !(registerMethod == "email" || registerMethod == "phone") {
		response := utils.FormatErrorResponse("Wrong method, use 'email' or 'phone'", "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := repositories.GetUser(credentials, registerMethod)
	if err != nil {
		response := utils.FormatErrorResponse("Server couldn't process your request", "")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if user.ID > 0 {
		response := utils.FormatErrorResponse(fmt.Sprintf("A user with this %s already exists", registerMethod), "")
		ctx.JSON(http.StatusOK, response)
		return
	}
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
