package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Version(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/version/")
	{
		group.GET("/latest/:platform", services.GetLatestVersion)
		group.GET("/check/:platform/:current_version", services.CheckForUpdates)
		group.GET("/", services.GetVersions)
		group.GET("/:uuid", services.GetVersionByID)

		group.POST("/", middlewares.GuardAdmin, services.CreateVersion)
		group.PUT("/:uuid", middlewares.GuardAdmin, services.UpdateVersion)
		group.DELETE("/:uuid", middlewares.GuardAdmin, services.DeleteVersion)
		group.POST("/:uuid/activate", middlewares.GuardAdmin, services.ActivateVersion)
		group.POST("/:uuid/deprecate", middlewares.GuardAdmin, services.DeprecateVersion)
	}
}
