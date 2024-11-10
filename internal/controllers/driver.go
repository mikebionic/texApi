package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Driver(router *gin.Engine) {
	group := router.Group("texapp/driver/")

	group.GET("/", services.GetDriverList)
	group.GET("/:id", services.GetDriver)
	group.POST("/", middlewares.Guard, services.CreateDriver)
	group.PUT("/:id", middlewares.Guard, services.UpdateDriver)
	group.DELETE("/:id", middlewares.Guard, services.DeleteDriver)

}
