package controllers

import (
	"texApi/pkg/middlewares"
	"texApi/pkg/services"

	"github.com/gin-gonic/gin"
)

func Services(router *gin.Engine) {
	group := router.Group("texapp/services")

	group.GET("", services.GetServices)
	group.GET("/list", services.GetServiceList)
	group.GET("/:id", middlewares.Guard, services.GetService)
	group.POST("", middlewares.Guard, services.CreateService)
	group.POST("/:id/image", middlewares.Guard, services.SetServiceImage)
	group.PUT("", middlewares.Guard, services.UpdateService)
	group.DELETE("/:id", middlewares.Guard, services.DeleteService)
}
