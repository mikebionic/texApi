package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Geofence(router *gin.Engine) {
	gpsGroup := router.Group(config.ENV.API_PREFIX + "/geofence/")
	{
		gpsGroup.POST("/", middlewares.Guard, services.CreateGeofence)
		gpsGroup.GET("/:geofence_id", middlewares.Guard, services.GetGeofence)
		gpsGroup.GET("/", middlewares.Guard, services.GetGeofenceList)
		gpsGroup.PUT("/:geofence_id", middlewares.Guard, services.UpdateGeofence)
		gpsGroup.DELETE("/:geofence_id", middlewares.Guard, services.DeleteGeofence)
		//gpsGroup.GET("/events", middlewares.Guard, services.GetGeofenceEvents)
	}
}
