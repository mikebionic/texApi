package internal

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("texsession", store))

	router.Use(middlewares.Cors)
	router.Static("/texapp/uploads", config.ENV.UPLOAD_PATH)

	log.SetOutput(gin.DefaultWriter)
	controllers.Content(router)
	controllers.ContentType(router)
	controllers.Auth(router)
	controllers.Company(router)
	controllers.Driver(router)
	controllers.Vehicle(router)

	return router
}
