package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func PackagingType(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/packaging-type/")

	group.GET("/", services.GetPackagingTypes)
	group.GET("/:id", services.GetPackagingType)
	group.POST("/", middlewares.GuardAdmin, services.CreatePackagingType)
	group.PUT("/:id", middlewares.GuardAdmin, services.UpdatePackagingType)
	group.DELETE("/:id", middlewares.GuardAdmin, services.DeletePackagingType)

}
