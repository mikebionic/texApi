package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Offer(router *gin.Engine) {
	group := router.Group("texapp/request/")

	//user as a company
	group.POST("/", services.CreateOffer)
	group.PUT("/:id", services.UpdateOffer)

	//user's companies
	group.GET("/company/", services.GetCompanyOffers)

	// TODO: Implement pagination
	// Everyone see this in BIDS section
	group.GET("/", services.GetOffers)

	//Only user's created ones, or admin
	group.DELETE("/:id", services.DeleteOffer)
}
