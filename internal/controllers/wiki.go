package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Wiki(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/wiki/")
	{
		group.GET("/", services.GetWikis)
		group.GET("/slug/:slug", services.GetWikiBySlug)
		group.GET("/categories", services.GetWikiCategories)

		group.POST("/", middlewares.GuardAdmin, services.CreateWiki)
		group.PUT("/:id", middlewares.GuardAdmin, services.UpdateWiki)
		group.DELETE("/:id", middlewares.GuardAdmin, services.DeleteWiki)
	}
}
