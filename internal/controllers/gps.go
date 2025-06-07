package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
)

func GPS(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/gps/")
	{
		group.POST("/trip/start/", services.StartTrip)
		group.POST("/trip/end/", services.EndTrip)
		group.GET("/trip/", services.GetTrips)

		group.POST("/log/", services.CreateGPSLogs)
		group.GET("/info/", services.GetGPSLogs)
		group.GET("/info/position/", services.GetLastPositions)
	}
}
