package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Cargo(router *gin.Engine) {
	group := router.Group("texapp/cargo/")

	group.GET("/", middlewares.Guard, services.GetCargoList)
	group.GET("/:id", middlewares.Guard, services.GetCargo)
	group.POST("/", middlewares.Guard, services.CreateCargo)
	group.PUT("/:id", middlewares.Guard, services.UpdateCargo)
	group.DELETE("/:id", middlewares.Guard, services.DeleteCargo)
}
