package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Offer(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/offer/")
	group.Use(middlewares.Guard)

	group.GET("/", services.GetOfferList)
	group.GET("/my/", services.GetMyOfferList)
	group.GET("/:id", services.GetOffer)
	group.POST("/", middlewares.Guard, services.CreateOffer)
	group.PUT("/:id", middlewares.Guard, services.UpdateOffer)
	group.DELETE("/:id", middlewares.Guard, services.DeleteOffer)
}
