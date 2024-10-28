package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Driver(router *gin.Engine) {
	group := router.Group("texapp/driver/")

	group.GET("/", services.GetDrivers)
	group.GET("/:id", services.SingleDriver)
	group.POST("/", services.CreateDriver)
	group.PUT("/:id", services.UpdateDriver)
	group.DELETE("/:id", services.DeleteDriver)
}
