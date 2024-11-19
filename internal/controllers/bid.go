package controllers

import (
	"texApi/config"
	"texApi/internal/services"

	"github.com/gin-gonic/gin"
)

func Bid(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/bid/")

	group.GET("/request/", services.GetRequestBids)
	group.GET("/user/", services.GetUserBids)
	group.POST("/create/", services.CreateBid)

	// Only company of request can approve or decline
	group.POST("/state/", services.ChangeBidState)

	// Only user or admin can, if not approved
	group.DELETE("/:id", services.DeleteUserBid)
}
