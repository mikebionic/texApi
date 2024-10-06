package services

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"texApi/config"
	"texApi/internal/repositories"
	"texApi/pkg/utils"
)

func UserLogin(ctx *gin.Context) {
	loginMethod := ctx.GetHeader("LoginMethod")
	if loginMethod == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Login Method"})
		return
	}
	username, password, err := utils.ExtractBasicAuth(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := repositories.GetUser(username, loginMethod)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if config.ENV.ENCRYPT_PASSWORDS > 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
			return
		}
	} else {
		if user.Password != password {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
			return
		}
	}

	session := sessions.Default(ctx)
	session.Set("userID", user.ID)
	session.Set("role", user.RoleID)
	session.Save()

	ctx.JSON(200, gin.H{"message": "Login successful", "user": user})
}

func UserGetMe(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("userID")

	if userID == nil {
		ctx.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	user := repositories.GetUserById(userID.(int))
	if user.ID == 0 {
		ctx.JSON(404, gin.H{"message": "User not found"})
		return
	}

	ctx.JSON(200, user)
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.JSON(200, gin.H{"message": "Logged out successfully"})
}

func RegisterRequest(ctx *gin.Context) {
	loginMethod := ctx.GetHeader("LoginMethod")
	roleID := ctx.GetHeader("RoleID")
	credentials := ctx.GetHeader("Credentials")
	registerMethod := ctx.GetHeader("RegisterMethod")
	//subRoleID := ctx.GetHeader("SubRoleID")
	if loginMethod == "" || roleID == "" || credentials == "" || registerMethod == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request, missing required params"})
		return
	}

	if !(registerMethod == "email" || registerMethod == "phone") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong method, use 'email' or 'phone'"})
		return
	}

	// check email or phone number
	// register method

}

func GetOAuthCallbackFunction(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Set("provider", provider)
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	fmt.Println(user)
	dbuser, err := repositories.GetUser(user.Email, "email")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	if dbuser.ID == 0 {
		fmt.Println("REGISTER USER")
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
