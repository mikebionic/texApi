package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Driver(router *gin.Engine) {
	group := router.Group("texapp/driver/")

	group.GET("/", services.GetDriverList)
	group.GET("/:id", services.GetDriver)
	group.POST("/", services.CreateDriver)
	group.PUT("/:id", services.UpdateDriver)
	group.DELETE("/:id", services.DeleteDriver)

}
