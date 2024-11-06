package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Bid(router *gin.Engine) {
	group := router.Group("texapp/bid/")

	group.GET("/request/", services.GetRequestBids)
	group.GET("/user/", services.GetUserBids)
	group.POST("/create/", services.CreateBid)

	// Only company of request can approve or decline
	group.POST("/state/", services.ChangeBidState)

	// Only user or admin can, if not approved
	group.DELETE("/:id", services.DeleteUserBid)
}
