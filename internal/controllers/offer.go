package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Offer(router *gin.Engine) {
	group := router.Group("texapp/offer/")
	group.Use(middlewares.Guard)

	group.GET("/", services.GetOfferList)
	group.GET("/:id", services.GetOffer)
	group.POST("/", services.CreateOffer)
	group.PUT("/:id", services.UpdateOffer)
	group.DELETE("/:id", services.DeleteOffer)
}
