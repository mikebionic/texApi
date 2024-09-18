package src

import (
	"io"
	"log"
	"os"
	"texApi/config"
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
	router.Static("/texApi/uploads", config.ENV.UPLOAD_PATH)

	log.SetOutput(gin.DefaultWriter)

	return router
}
