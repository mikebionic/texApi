package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func OfferResponse(router *gin.Engine) {
	offerResponseGroup := router.Group(config.ENV.API_PREFIX + "/offer-response/")
	{
		offerResponseGroup.GET("/", services.GetDetailedOfferResponseList)
		offerResponseGroup.GET("/:id", services.GetOfferResponse)
		offerResponseGroup.POST("/", middlewares.Guard, services.CreateOfferResponse)
		offerResponseGroup.PUT("/:id", middlewares.Guard, services.UpdateOfferResponse)
		offerResponseGroup.DELETE("/:id", middlewares.Guard, services.DeleteOfferResponse)
	}
}
