package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Driver(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/driver/")

	group.GET("/", middlewares.Guard, services.GetDriverList)
	group.GET("/:id", middlewares.Guard, middlewares.ViewCounterMiddleware("tbl_driver"), services.GetDriver)
	group.POST("/", middlewares.Guard, services.CreateDriver)
	group.PUT("/:id", middlewares.Guard, services.UpdateDriver)
	group.DELETE("/:id", middlewares.Guard, services.DeleteDriver)
	group.GET("/filter/", services.GetFilteredDriverList)

}
