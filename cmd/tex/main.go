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

const (
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

func NewGoogleAuth() {
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd
	store.Options.Path = "/"
	gothic.Store = store

	goth.UseProviders(
		google.New(config.ENV.GoogleClientID, config.ENV.GoogleClientSecret, "http://localhost:7000/texapp/auth/oauth/google/callback/"))
}

func main() {
	config.InitConfig()
	database.InitDB()
	NewGoogleAuth()
	server := app.InitApp()
	address := fmt.Sprintf("%v:%v", config.ENV.API_HOST, config.ENV.API_PORT)
	server.Run(address)
}
