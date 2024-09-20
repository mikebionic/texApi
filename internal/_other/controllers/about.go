package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/_other/services"
	"texApi/pkg/middlewares"
)

func AboutUs(router *gin.Engine) {
	group := router.Group("texapp/about")

	group.GET("", middlewares.Guard, services.GetAboutUsAll)
	group.GET("/by-user", services.GetAboutUsForUser)
	group.GET("/:id", middlewares.Guard, services.GetAboutUs)
	group.POST("", middlewares.Guard, services.CreateAboutUs)
	group.PUT("", middlewares.Guard, services.UpdateAboutUs)
	group.DELETE("/:id", middlewares.Guard, services.DeleteAboutUs)
}
