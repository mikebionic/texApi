package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Company(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/company/")

	group.GET("/", services.GetCompanyList)
	group.GET("/:id", services.GetCompany)
	group.POST("/", middlewares.Guard, services.CreateCompany)
	group.PUT("/:id", middlewares.Guard, services.UpdateCompany)
	group.DELETE("/:id", middlewares.Guard, services.DeleteCompany)

}
