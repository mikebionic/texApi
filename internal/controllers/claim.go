package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Claim(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/claim/")
	{
		group.GET("/", services.GetFilteredClaims)
		group.POST("/new/", services.NewClaim)
		group.PUT("/", middlewares.GuardAdmin, services.UpdateClaim)
		group.DELETE("/", middlewares.GuardAdmin, services.DeleteClaim)
	}
}
