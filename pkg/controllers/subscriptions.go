package controllers

import (
	"texApi/pkg/middlewares"
	"texApi/pkg/services"

	"github.com/gin-gonic/gin"
)

func Subscriptions(router *gin.Engine) {
	group := router.Group("texapp/subscriptions")

	group.GET("", services.GetSubscriptions)
	group.GET("/:id", middlewares.Guard, services.GetSubscription)
	group.POST("", middlewares.Guard, services.CreateSubscription)
	group.PUT("", middlewares.Guard, services.UpdateSubscription)
	group.DELETE("/:id", middlewares.Guard, services.DeleteSubscription)
}
