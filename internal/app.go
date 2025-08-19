package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"texApi/config"
	"texApi/database"
	"texApi/internal/chat"
	"texApi/internal/controllers"
	"texApi/internal/firebasePush"
	"texApi/pkg/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
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
	router.Use(gin.Recovery())

	store := cookie.NewStore([]byte(config.ENV.API_SECRET))
	router.Use(sessions.Sessions("google-auth-session", store))
	router.Use(sessions.Sessions(config.ENV.API_PREFIX, store))
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

	router.Use(middlewares.UpdateLastActive)
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
	controllers.OfferResponse(router)
	controllers.PackagingType(router)
	controllers.Cargo(router)
	controllers.Media(router)
	controllers.VerifyRequest(router)
	controllers.PlanMove(router)
	controllers.UserLog(router)
	controllers.Organization(router)
	controllers.GPS(router)
	controllers.News(router)
	controllers.Version(router)
	controllers.Plan(router)
	controllers.Analytics(router)
	controllers.Wiki(router)
	controllers.PriceQuote(router)
	controllers.Claim(router)
	chat.Chat(router)
	firebasePush.Controller(router)

	return router
}
