package src

import (
	"io"
	"log"
	"os"
	"texApi/config"
	"texApi/pkg/controllers"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func InitApp() *gin.Engine {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		file, _ := os.Create("gin.log")
		gin.DefaultWriter = io.MultiWriter(file)
	}

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger())

	router.Use(middlewares.Cors)
	router.Static("/texapp/uploads", config.ENV.UPLOAD_PATH)

	log.SetOutput(gin.DefaultWriter)

	controllers.WebSocket(router)
	controllers.Auth(router)
	controllers.Services(router)
	controllers.Users(router)
	controllers.Workers(router)
	controllers.Subscriptions(router)
	controllers.Statuses(router)
	controllers.AboutUs(router)
	controllers.Content(router)

	// // Initialize the repositories, services, and controllers
	// contentRepo := repositories.NewContentRepository(DB)
	// contentService := services.NewContentService(contentRepo)
	// contentController := controllers.NewContentController(contentService)
	// // Routes for content
	// router.GET("/content", contentController.GetAll)
	// router.GET("/content/:id", contentController.GetByID)
	// router.GET("/content/title", contentController.GetByTitle)
	// router.GET("/content/uuid/:uuid", contentController.GetByUUID)
	// router.GET("/content/type/:content_type_id", contentController.GetByContentTypeID)

	// router.POST("/content", contentController.Create)
	// router.PATCH("/content/:id", contentController.Update)
	// router.DELETE("/content/:id", contentController.Delete) // Add delete route

	return router
}
