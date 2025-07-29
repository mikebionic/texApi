package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func PriceQuote(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/price-quote/")
	{
		group.GET("/", services.GetPriceQuoteList)
		group.GET("/analyze", services.GetPriceQuoteWithOfferAnalysis)
		group.POST("/", middlewares.GuardAdmin, services.CreatePriceQuote)
		group.PUT("/:id", middlewares.GuardAdmin, services.UpdatePriceQuote)
		group.DELETE("/:id", middlewares.GuardAdmin, services.DeletePriceQuote)
	}
}
