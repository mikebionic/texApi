package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
	"texApi/pkg/middlewares"
)

func Company(router *gin.Engine) {
	group := router.Group("texapp/company/")

	group.GET("/", services.GetCompanyList)
	group.GET("/:id", services.GetCompany)
	group.POST("/", middlewares.Guard, services.CreateCompany)
	group.PUT("/:id", middlewares.Guard, services.UpdateCompany)
	group.DELETE("/:id", middlewares.Guard, services.DeleteCompany)

}
