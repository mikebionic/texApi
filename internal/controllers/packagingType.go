package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func PackagingType(router *gin.Engine) {
	group := router.Group("texapp/packaging-type/")

	group.GET("/", services.GetPackagingTypes)
	group.GET("/:id", services.GetPackagingType)
	group.POST("/", middlewares.GuardAdmin, services.CreatePackagingType)
	group.PUT("/:id", middlewares.GuardAdmin, services.UpdatePackagingType)
	group.DELETE("/:id", middlewares.GuardAdmin, services.DeletePackagingType)

}
