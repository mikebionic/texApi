package internal

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"texApi/config"
	"texApi/internal/controllers"
	"texApi/pkg/middlewares"
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

	// controllers.WebSocket(router)
	// controllers.Auth(router)
	// controllers.Services(router)
	// controllers.Users(router)
	// controllers.Workers(router)
	// controllers.Subscriptions(router)
	// controllers.Statuses(router)
	// controllers.AboutUs(router)
	controllers.Content(router)

	return router
}
