package controllers

import (
	"texApi/config"
	"texApi/internal/services"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Company(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/company/")
	{
		group.GET("/:id", services.GetCompany)
		//group.GET("/followers/:id/", services.GetCompanyFollowers)
		//group.GET("/following/:id/", services.GetCompanyFollowing)

		group.GET("/", services.GetCompanyList)
		//group.GET("/all/", services.GetCompanyAllList)
		protected := group.Use(middlewares.Guard)
		{
			//protected.GET("/", services.GetPublicCompanyList)

			protected.POST("/", services.CreateCompany)
			protected.PUT("/:id", services.UpdateCompany)
			protected.DELETE("/:id", services.DeleteCompany)

			//protected.POST("/follow/:id", services.CompanyFollow)
			//protected.DELETE("/follow/:id", services.CompanyUnfollow)
		}
		//group.POST("/block/", services.ChangeCompanyBlocked)
	}
}
