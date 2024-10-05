package services

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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

	user, err := repositories.GetUser(username, password, loginMethod)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
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
