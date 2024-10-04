package middlewares

import (
	"github.com/gin-contrib/sessions"
	"strings"
	"texApi/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Guard(ctx *gin.Context) {
	authorization := ctx.Request.Header["Authorization"]

	if len(authorization) == 0 {
		ctx.AbortWithStatus(401)
		return
	}

	bearer := strings.Split(authorization[0], "Bearer ")

	if len(bearer) == 0 || len(bearer) == 1 {
		ctx.AbortWithStatus(401)
		return
	}

	token := bearer[1]
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.ENV.ACCESS_KEY), nil
		},
	)

	if err != nil {
		ctx.AbortWithStatus(403)
		return
	}

	ctx.Set("id", int(claims["id"].(float64)))
	ctx.Set("role", claims["role"])
	ctx.Next()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get("userID")

		if userID == nil {
			ctx.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
