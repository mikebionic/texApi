package config

import (
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var SocketClients = make(map[string]*websocket.Conn)

type Config struct {
	API_HOST    string
	API_PORT    string
	UPLOAD_PATH string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	ACCESS_KEY   string
	ACCESS_TIME  time.Duration
	REFRESH_KEY  string
	REFRESH_TIME time.Duration
}

var ENV Config

func InitConfig() {
	godotenv.Load()

	ENV.API_HOST = os.Getenv("API_HOST")
	ENV.API_PORT = os.Getenv("API_PORT")
	ENV.UPLOAD_PATH = os.Getenv("UPLOAD_PATH")

	ENV.DB_HOST = os.Getenv("DB_HOST")
	ENV.DB_PORT = os.Getenv("DB_PORT")
	ENV.DB_USER = os.Getenv("DB_USER")
	ENV.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	ENV.DB_NAME = os.Getenv("DB_NAME")

	ENV.ACCESS_KEY = os.Getenv("ACCESS_KEY")
	AT, _ := time.ParseDuration(os.Getenv(("ACCESS_TIME")))
	ENV.ACCESS_TIME = AT

	ENV.REFRESH_KEY = os.Getenv("REFRESH_KEY")
	RT, _ := time.ParseDuration(os.Getenv(("REFRESH_TIME")))
	ENV.REFRESH_TIME = RT
}
