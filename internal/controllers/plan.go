package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Plan(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/plan/")
	{
		group.GET("/", services.GetPlans)
		group.GET("/:uuid", services.GetPlanByID)

		group.POST("/", middlewares.GuardAdmin, services.CreatePlan)
		group.PUT("/:uuid", middlewares.GuardAdmin, services.UpdatePlan)
		group.DELETE("/:uuid", middlewares.GuardAdmin, services.DeletePlan)
	}
}
