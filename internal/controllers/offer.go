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

	group.GET("/detailed/", services.GetDetailedOfferList)
	group.GET("/", services.GetOfferListUpdate)
	group.GET("/my/", services.GetMyOfferListUpdate)
	group.GET("/:id", services.GetOffer)
	group.POST("/", services.CreateOffer)
	group.PUT("/:id", services.UpdateOffer)
	group.DELETE("/:id", services.DeleteOffer)
}
