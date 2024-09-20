package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/_other/services"
	"texApi/pkg/middlewares"
)

func Workers(router *gin.Engine) {
	group := router.Group("texapp/workers", middlewares.Guard)

	group.GET("", services.GetWorkers)
	group.GET("/:id", services.GetWorker)
	group.GET("/me", services.GetWorkerMe)
	group.POST("", services.CreateWorker)
	group.POST("/:id/image", services.SetWorkerImage)
	group.POST("/image", services.SetWorkerImage)
	// group.PUT("", services.UpdateWorker)
	group.DELETE("/:id", services.DeleteWorker)
}
