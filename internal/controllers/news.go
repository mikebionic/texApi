package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func News(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/news/")
	{
		group.GET("/", services.GetNews)
		group.GET("/:id", services.GetNewsByID)
		group.POST("/", middlewares.GuardAdmin, services.CreateNews)
		group.PUT("/:id", middlewares.GuardAdmin, services.UpdateNews)
		group.DELETE("/:id", middlewares.GuardAdmin, services.DeleteNews)
	}
}
