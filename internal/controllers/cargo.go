package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Cargo(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/cargo/")

	group.GET("/detailed/", middlewares.Guard, services.GetDetailedCargoList)
	group.GET("/", middlewares.Guard, services.GetCargoList)
	group.GET("/:id", middlewares.Guard, services.GetCargo)
	group.POST("/", middlewares.Guard, services.CreateCargo)
	group.PUT("/:id", middlewares.Guard, services.UpdateCargo)
	group.DELETE("/:id", middlewares.Guard, services.DeleteCargo)
}
