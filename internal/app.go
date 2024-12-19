package internal

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"texApi/config"
	"texApi/database"
	"texApi/internal/controllers"
	"texApi/pkg/middlewares"
	"time"
)

func InitApp() *gin.Engine {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		file, _ := os.Create("gin.log")
		gin.DefaultWriter = io.MultiWriter(file)
	}

	middlewares.InitializeViewTracker(database.DB, 10)

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger())

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("texsession", store))

	//router.Use(middlewares.Cors)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"*"},
		MaxAge:          12 * time.Hour,
	}))
	router.Use(func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	})
	router.Static(fmt.Sprintf("/%s/uploads/", config.ENV.API_PREFIX), config.ENV.UPLOAD_PATH)
	router.Static(fmt.Sprintf("/%s/assets/", config.ENV.API_PREFIX), "assets/")

	log.SetOutput(gin.DefaultWriter)
	controllers.Content(router)
	controllers.ContentType(router)
	controllers.Auth(router)
	controllers.Company(router)
	controllers.Driver(router)
	controllers.Vehicle(router)
	controllers.Offer(router)
	controllers.Bid(router)
	controllers.PackagingType(router)
	controllers.Cargo(router)
	controllers.Media(router)

	return router
}
