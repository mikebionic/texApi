package main

import (
	"fmt"
	"log"
	"texApi/config"
	"texApi/database"
	app "texApi/internal"
	"texApi/pkg/smtp"
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

	server := app.InitApp()
	address := fmt.Sprintf("%v:%v", config.ENV.API_HOST, config.ENV.API_PORT)
	if err := server.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
