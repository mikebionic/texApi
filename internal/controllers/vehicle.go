package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Vehicle(router *gin.Engine) {
	group := router.Group("texapp/vehicle/")

	group.GET("/", services.GetVehicles)
	group.GET("/:id", services.GetVehicle)
	group.POST("/", services.CreateVehicle)
	group.PUT("/:id", services.UpdateVehicle)
	group.DELETE("/:id", services.DeleteVehicle)

	// Vehicle Brand routes
	group.GET("/brand", services.GetVehicleBrands)
	group.GET("/brand/:id", services.SingleVehicleBrand)
	group.POST("/brand", services.CreateVehicleBrand)
	group.PUT("/brand/:id", services.UpdateVehicleBrand)
	group.DELETE("/brand/:id", services.DeleteVehicleBrand)

	// Vehicle Type routes
	group.GET("/type", services.GetVehicleTypes)
	group.GET("/type/:id", services.SingleVehicleType)
	group.POST("/type", services.CreateVehicleType)
	group.PUT("/type/:id", services.UpdateVehicleType)
	group.DELETE("/type/:id", services.DeleteVehicleType)

	// Vehicle Model routes
	group.GET("/model", services.GetVehicleModels)
	group.GET("/model/:id", services.SingleVehicleModel)
	group.POST("/model", services.CreateVehicleModel)
	group.PUT("/model/:id", services.UpdateVehicleModel)
	group.DELETE("/model/:id", services.DeleteVehicleModel)
}
