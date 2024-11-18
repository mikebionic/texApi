package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var SocketClients = make(map[string]*websocket.Conn)

type Config struct {
	API_HOST       string
	API_SERVER_URL string
	API_PREFIX     string
	API_PORT       string
	API_SECRET     string
	API_DEBUG      bool

	UPLOAD_PATH           string
	MAX_FILE_UPLOAD_COUNT int
	STATIC_URL            string

	ENCRYPT_PASSWORDS bool
	SESSION_MAX_AGE   int

	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	ACCESS_KEY   string
	ACCESS_TIME  time.Duration
	REFRESH_KEY  string
	REFRESH_TIME time.Duration

	GLE_KEY      string
	GLE_SECRET   string
	GLE_CALLBACK string

	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_MAIL     string
	SMTP_PASSWORD string
	APP_LOGO_URL  string
}

var ENV Config

func InitConfig() {
	godotenv.Load()
	ENV.API_HOST = os.Getenv("API_HOST")
	ENV.API_SERVER_URL = os.Getenv("API_SERVER_URL")
	ENV.API_PREFIX = os.Getenv("API_PREFIX")
	ENV.API_PORT = os.Getenv("API_PORT")
	ENV.API_DEBUG = os.Getenv("DEBUG") == "true"
	ENV.API_SECRET = os.Getenv("API_SECRET")
	ENV.SESSION_MAX_AGE = 86400 * 30 // TODO: WTF?

	ENV.UPLOAD_PATH = os.Getenv("UPLOAD_PATH")
	ENV.MAX_FILE_UPLOAD_COUNT, _ = strconv.Atoi(os.Getenv("MAX_FILE_UPLOAD_COUNT"))
	ENV.STATIC_URL = fmt.Sprintf("/%s/uploads/", ENV.API_PREFIX)

	ENV.ENCRYPT_PASSWORDS = os.Getenv("ENCRYPT_PASSWORDS") == "true"

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

	ENV.GLE_KEY = os.Getenv("GLE_KEY")
	ENV.GLE_SECRET = os.Getenv("GLE_SECRET")
	ENV.GLE_CALLBACK = os.Getenv("GLE_CALLBACK")

	ENV.SMTP_HOST = os.Getenv("SMTP_HOST")
	ENV.SMTP_PORT = os.Getenv("SMTP_PORT")
	ENV.SMTP_MAIL = os.Getenv("SMTP_MAIL")
	ENV.SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	ENV.APP_LOGO_URL = os.Getenv("APP_LOGO_URL")
}
