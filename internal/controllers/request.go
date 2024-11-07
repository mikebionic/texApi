package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Request(router *gin.Engine) {
	group := router.Group("texapp/request/")

	//user as a company
	group.POST("/", services.CreateRequest)
	group.PUT("/:id", services.UpdateRequest)

	//user's companies
	group.GET("/company/", services.GetCompanyRequests)

	// TODO: Implement pagination
	// Everyone see this in BIDS section
	group.GET("/", services.GetRequests)

	//Only user's created ones, or admin
	group.DELETE("/:id", services.DeleteRequest)
}
