package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Organization(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/organization/")
	{
		group.GET("/", services.GetOrganization)
		group.POST("/", middlewares.GuardAdmin, services.CreateOrganization)
		group.PUT("/:id", middlewares.GuardAdmin, services.UpdateOrganization)
		group.DELETE("/:id", middlewares.GuardAdmin, services.DeleteOrganization)
	}
}
