package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func GPS(router *gin.Engine) {
	gpsGroup := router.Group(config.ENV.API_PREFIX + "/gps/")
	{
		// Location recording endpoints
		gpsGroup.POST("/location", middlewares.Guard, services.RecordGPSLocation)
		gpsGroup.POST("/batch-location", middlewares.Guard, services.RecordBatchGPSLocation)

		// Location retrieval endpoints
		gpsGroup.GET("/current-location/:vehicle_id", middlewares.Guard, services.GetCurrentVehicleLocation)
		gpsGroup.GET("/vehicle-history", middlewares.Guard, services.GetVehicleLocationHistory)
		gpsGroup.GET("/driver-history", middlewares.Guard, services.GetDriverLocationHistory)

		// Trip management
		gpsGroup.GET("/trip/:trip_id", middlewares.Guard, services.GetTrip)
		gpsGroup.GET("/trips", middlewares.Guard, services.GetTripList)
		gpsGroup.POST("/trip", middlewares.Guard, services.CreateTrip)
		gpsGroup.PUT("/trip/:trip_id", middlewares.Guard, services.UpdateTrip)

		// Geofencing
		gpsGroup.POST("/geofence", middlewares.Guard, services.CreateGeofence)
		gpsGroup.GET("/geofence/:geofence_id", middlewares.Guard, services.GetGeofence)
		gpsGroup.GET("/geofence", middlewares.Guard, services.GetGeofenceList)
		gpsGroup.PUT("/geofence/:geofence_id", middlewares.Guard, services.UpdateGeofence)
		gpsGroup.DELETE("/geofence/:geofence_id", middlewares.Guard, services.DeleteGeofence)
		//gpsGroup.GET("/geofence/events", middlewares.Guard, services.GetGeofenceEvents)
		//
		//// Analytics and reporting
		//gpsGroup.GET("/analytics/vehicle/:vehicle_id", middlewares.Guard, services.GetVehicleAnalytics)
		//gpsGroup.GET("/analytics/driver/:driver_id", middlewares.Guard, services.GetDriverAnalytics)
		//
		//// Map data
		//gpsGroup.GET("/map-data", middlewares.Guard, services.GetMapData)
	}
}
