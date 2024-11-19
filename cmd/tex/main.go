package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"texApi/config"
	"texApi/database"
	app "texApi/internal"
	"texApi/pkg/smtp"
)

func NewGoogleAuth() {
	store := sessions.NewCookieStore([]byte(config.ENV.API_SECRET))
	store.MaxAge(config.ENV.SESSION_MAX_AGE)
	store.Options.HttpOnly = true
	store.Options.Secure = !config.ENV.API_DEBUG
	store.Options.Path = "/"
	gothic.Store = store

	goth.UseProviders(
		google.New(config.ENV.GLE_KEY, config.ENV.GLE_SECRET, fmt.Sprintf("%s/%s/%s", config.ENV.API_SERVER_URL, config.ENV.API_PREFIX, config.ENV.GLE_CALLBACK)))
}

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
	NewGoogleAuth()
	server := app.InitApp()
	address := fmt.Sprintf("%v:%v", config.ENV.API_HOST, config.ENV.API_PORT)
	server.Run(address)
}
