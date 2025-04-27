package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

// VerifyRequest sets up routes for verification request management
func VerifyRequest(router *gin.Engine) {
	verifyGroup := router.Group(config.ENV.API_PREFIX + "/verify-request/")
	{
		verifyGroup.GET("/", middlewares.GuardAdmin, services.GetVerifyRequestList)
		verifyGroup.GET("/:id", middlewares.GuardAdmin, services.GetVerifyRequest)
		verifyGroup.POST("/", middlewares.Guard, services.CreateVerifyRequest)
		verifyGroup.PUT("/:id", middlewares.GuardAdmin, services.UpdateVerifyRequest)
		verifyGroup.DELETE("/:id", middlewares.GuardAdmin, services.DeleteVerifyRequest)
	}
}

// PlanMove sets up routes for plan movement management
func PlanMove(router *gin.Engine) {
	planGroup := router.Group(config.ENV.API_PREFIX + "/plan-move/")
	{
		planGroup.GET("/", middlewares.GuardAdmin, services.GetPlanMovesList)
		planGroup.GET("/:id", middlewares.GuardAdmin, services.GetPlanMove)
		planGroup.POST("/", middlewares.Guard, services.CreatePlanMove)
		planGroup.PUT("/:id", middlewares.GuardAdmin, services.UpdatePlanMove)
		planGroup.DELETE("/:id", middlewares.GuardAdmin, services.DeletePlanMove)

		// Special endpoint for checking expired plans
		planGroup.POST("/check-expired", middlewares.GuardAdmin, services.CheckExpiredPlans)
	}
}

// UserLog sets up routes for user log management
func UserLog(router *gin.Engine) {
	logGroup := router.Group(config.ENV.API_PREFIX + "/user-log/")
	{
		logGroup.GET("/", middlewares.GuardAdmin, services.GetUserLogsList)
		logGroup.GET("/:id", middlewares.GuardAdmin, services.GetUserLog)
		logGroup.DELETE("/:id", middlewares.GuardAdmin, services.DeleteUserLog)
	}
}
