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
)

func NewGoogleAuth() {
	store := sessions.NewCookieStore([]byte(config.ENV.API_SECRET))
	store.MaxAge(config.ENV.SESSION_MAX_AGE)
	store.Options.HttpOnly = true
	store.Options.Secure = !config.ENV.API_DEBUG
	store.Options.Path = "/"
	gothic.Store = store

	goth.UseProviders(
		google.New(config.ENV.GoogleClientID, config.ENV.GoogleClientSecret, fmt.Sprintf("%s/texapp/auth/oauth/google/callback/", config.ENV.API_SERVER_URL)))
}

func main() {
	config.InitConfig()
	database.InitDB()
	NewGoogleAuth()
	server := app.InitApp()
	address := fmt.Sprintf("%v:%v", config.ENV.API_HOST, config.ENV.API_PORT)
	server.Run(address)
}
