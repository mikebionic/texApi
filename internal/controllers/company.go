package controllers

import (
	"github.com/gin-gonic/gin"
	"texApi/internal/services"
)

func Company(router *gin.Engine) {
	group := router.Group("texapp/company/")

	group.GET("/:id", services.GetCompany)
	group.POST("/", services.CreateCompany)
	group.PUT("/:id", services.UpdateCompany)
	group.DELETE("/:id", services.DeleteCompany)

}
