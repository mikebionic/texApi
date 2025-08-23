package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Newsletter(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/newsletter/")

	group.POST("/subscribe/", services.Subscribe)

	group.GET("/", middlewares.GuardAdmin, services.GetNewsletterList)
	group.PUT("/:id", middlewares.GuardAdmin, services.UpdateNewsletter)
	group.DELETE("/:id", middlewares.GuardAdmin, services.DeleteNewsletter)
}
