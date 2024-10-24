package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Vehicle(router *gin.Engine) {
	group := router.Group("texapp/vehicle/")

	group.GET("/:id", services.GetVehicle)
	group.POST("", services.CreateVehicle)
	group.PUT("/:id", services.UpdateVehicle)
	group.DELETE("/:id", services.DeleteVehicle)
}
