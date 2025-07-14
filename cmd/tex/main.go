package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"texApi/config"
	"texApi/database"
	app "texApi/internal"
	"texApi/internal/scheduler"
	"texApi/pkg/smtp"
	"time"
)

func setupSMTPConfig() {
	smtp.DefaultConfig.SMTPHost = config.ENV.SMTP_HOST
	smtp.DefaultConfig.SMTPPort = config.ENV.SMTP_PORT
	smtp.DefaultConfig.SenderEmail = config.ENV.SMTP_MAIL
	smtp.DefaultConfig.Password = config.ENV.SMTP_PASSWORD
	smtp.DefaultConfig.LogoURL = config.ENV.APP_LOGO_URL
}

func main() {
	config.InitConfig()
	database.InitDB()
	setupSMTPConfig()

	analyticsScheduler := scheduler.NewAnalyticsScheduler()
	if err := analyticsScheduler.Start(); err != nil {
		log.Fatalf("Failed to start analytics scheduler: %v", err)
	}

	router := app.InitApp()
	address := fmt.Sprintf("%v:%v", config.ENV.API_HOST, config.ENV.API_PORT)

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server running at %s\n", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down gracefully...")

	// Stop background jobs
	analyticsScheduler.Stop()

	// Gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped properly")
}
