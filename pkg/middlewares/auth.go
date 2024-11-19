package middlewares

import (
	"net/http"
	"strings"
	"texApi/config"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Guard(ctx *gin.Context) {
	authorization := ctx.Request.Header["Authorization"]
	if len(authorization) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", ""))
		return
	}

	bearer := strings.Split(authorization[0], "Bearer ")
	if len(bearer) == 0 || len(bearer) == 1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", ""))
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
		ctx.AbortWithStatusJSON(http.StatusForbidden, utils.FormatErrorResponse("Forbidden", err.Error()))
		return
	}

	ctx.Set("id", int(claims["id"].(float64)))
	ctx.Set("roleID", int(claims["roleID"].(float64)))
	ctx.Set("companyID", int(claims["companyID"].(float64)))
	ctx.Set("role", claims["role"])
	ctx.Next()
}

func GuardAdmin(ctx *gin.Context) {
	authorization := ctx.Request.Header["Authorization"]
	if len(authorization) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", ""))
		return
	}

	bearer := strings.Split(authorization[0], "Bearer ")
	if len(bearer) == 0 || len(bearer) == 1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", ""))
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
		ctx.AbortWithStatusJSON(http.StatusForbidden, utils.FormatErrorResponse("Forbidden", err.Error()))
		return
	}

	if !(claims["role"] == "admin" || claims["role"] == "system") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse("Permission denied!", ""))
		return
	}

	ctx.Set("id", int(claims["id"].(float64)))
	ctx.Set("roleID", int(claims["roleID"].(float64)))
	ctx.Set("companyID", int(claims["companyID"].(float64)))
	ctx.Set("role", claims["role"])
	ctx.Next()
}
