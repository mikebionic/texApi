package services

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"texApi/internal/_other/schemas/request"
	"texApi/internal/repositories"
)

func UserLogin(ctx *gin.Context) {
	var loginForm request.LoginForm

	if err := ctx.BindJSON(&loginForm); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user := repositories.GetUserByPhone(loginForm.Phone)
	if user.ID == 0 {
		ctx.JSON(400, gin.H{"message": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		ctx.JSON(400, gin.H{"message": "Incorrect password"})
		return
	}

	session := sessions.Default(ctx)
	session.Set("userID", user.ID)
	session.Set("role", "user")
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
