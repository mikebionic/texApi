package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Company(router *gin.Engine) {
	group := router.Group("texapp/company/")

	//GET /texapp/company?page=1&per_page=10
	group.GET("/", services.GetCompanyList)
	group.GET("/:id", services.GetCompany)
	group.POST("/", services.CreateCompany)
	group.PUT("/:id", services.UpdateCompany)
	group.DELETE("/:id", services.DeleteCompany)

}
