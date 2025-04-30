package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func OfferResponse(router *gin.Engine) {
	offerResponseGroup := router.Group(config.ENV.API_PREFIX + "/offer-response/").Use(middlewares.Guard)
	{
		offerResponseGroup.GET("/", services.GetDetailedOfferResponseList)
		offerResponseGroup.GET("/:id", services.GetOfferResponse)
		offerResponseGroup.POST("/", services.CreateOfferResponse)
		offerResponseGroup.PUT("/:id", services.UpdateOfferResponse) // Accept Decline Offer Response
		offerResponseGroup.DELETE("/:id", services.DeleteOfferResponse)
	}
}
