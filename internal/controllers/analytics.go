package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Analytics(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/analytics/")
	{
		group.Use(middlewares.GuardAdmin)
		{
			group.GET("/", services.GetAnalytics)
			//group.GET("/stats", services.GetAnalyticsStats)
			//group.GET("/status", services.GetAnalyticsStatus)
		}

		admin := group.Group("/admin/")
		{
			admin.POST("/generate/", services.ForceGenerateAnalytics)
			admin.PUT("/config/", services.UpdateAnalyticsConfig)
			admin.GET("/config/", services.GetAnalyticsConfig)
		}
	}
}
