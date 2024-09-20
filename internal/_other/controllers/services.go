package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/_other/services"
	"texApi/pkg/middlewares"
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
