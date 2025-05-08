package middlewares

import (
	"log"
	"net/http"
	"strings"
	"texApi/config"
	"texApi/internal/repo"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GuardURLParam(ctx *gin.Context) {
	token := ctx.Query("token")
	if len(token) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse("Unauthorized", "Token is missing"))
		return
	}
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

func SysGuard(ctx *gin.Context) {
	token := ctx.GetHeader(config.ENV.SYSTEM_HEADER)
	if token != config.ENV.API_SECRET {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			utils.FormatErrorResponse("Unauthorized", "This endpoint is for internal use only"))
		return
	}
	ctx.Next()
}

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

func UpdateLastActive(ctx *gin.Context) {
	authorization := ctx.Request.Header["Authorization"]
	if len(authorization) == 0 {
		ctx.Next()
		return
	}
	bearer := strings.Split(authorization[0], "Bearer ")
	if len(bearer) == 0 || len(bearer) == 1 {
		ctx.Next()
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
		ctx.Next()
		return
	}

	companyID := int(claims["companyID"].(float64))
	if companyID == 0 {
		ctx.Next()
		return
	}

	go func() {
		err := repo.UpdateUserLastActive(companyID)
		if err != nil {
			log.Print("Failed to update last seen", err)
			ctx.Next()
			return
		}
	}()

	ctx.Next()
}
