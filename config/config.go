package config

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	API_HOST       string
	API_SERVER_URL string
	API_PREFIX     string
	API_PORT       string
	API_SECRET     string
	SYSTEM_HEADER  string
	SYSTEM_SECRET  string

	API_DEBUG bool

	UPLOAD_PATH      string
	MAX_FILE_SIZE    int64
	MAX_FILES_UPLOAD int
	STATIC_URL       string
	COMPRESS_IMAGES  int
	COMPRESS_SIZE    int
	COMPRESS_QUALITY int

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

	AppTZ      *time.Location
	TZAddHours time.Duration

	GLE_KEY           string
	GLE_SECRET        string
	GLE_CALLBACK      string
	GoogleOAuthConfig oauth2.Config

	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_MAIL     string
	SMTP_PASSWORD string
	APP_LOGO_URL  string

	FileUpload FileUpload
}

type FileUpload struct {
	MaxFileSize      int64 // in bytes (e.g., 10MB = 10 * 1024 * 1024)
	MaxFiles         int
	AllowedMimeTypes map[string]bool
	StorageBasePath  string
}

var ENV Config

func InitConfig() {
	godotenv.Load()
	ENV.API_HOST = os.Getenv("API_HOST")
	ENV.API_SERVER_URL = os.Getenv("API_SERVER_URL")
	ENV.API_PREFIX = os.Getenv("API_PREFIX")
	ENV.API_PORT = os.Getenv("API_PORT")
	ENV.API_DEBUG = os.Getenv("API_DEBUG") == "true"
	ENV.API_SECRET = os.Getenv("API_SECRET")
	ENV.SYSTEM_HEADER = os.Getenv("SYSTEM_HEADER")
	ENV.SYSTEM_SECRET = os.Getenv("SYSTEM_SECRET")
	ENV.SESSION_MAX_AGE = 86400 * 30 // TODO: WTF?

	ENV.UPLOAD_PATH = os.Getenv("UPLOAD_PATH")
	ENV.MAX_FILES_UPLOAD, _ = strconv.Atoi(os.Getenv("MAX_FILES_UPLOAD"))
	ENV.STATIC_URL = fmt.Sprintf("/%s/uploads/", ENV.API_PREFIX)
	ENV.COMPRESS_IMAGES, _ = strconv.Atoi(os.Getenv("COMPRESS_IMAGES"))
	ENV.COMPRESS_SIZE, _ = strconv.Atoi(os.Getenv("COMPRESS_SIZE"))
	ENV.COMPRESS_QUALITY, _ = strconv.Atoi(os.Getenv("COMPRESS_QUALITY"))

	ENV.ENCRYPT_PASSWORDS = os.Getenv("ENCRYPT_PASSWORDS") == "true"

	ENV.DB_HOST = os.Getenv("DB_HOST")
	ENV.DB_PORT = os.Getenv("DB_PORT")
	ENV.DB_USER = os.Getenv("DB_USER")
	ENV.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	ENV.DB_NAME = os.Getenv("DB_NAME")

	ENV.ACCESS_KEY = os.Getenv("ACCESS_KEY")
	accessTime, err := time.ParseDuration(os.Getenv("ACCESS_TIME"))
	if err != nil {
		ENV.ACCESS_TIME = 15 * time.Minute
		fmt.Printf("Warning: Invalid ACCESS_TIME, using default: %v\n", ENV.ACCESS_TIME)
	} else {
		ENV.ACCESS_TIME = accessTime
	}

	ENV.REFRESH_KEY = os.Getenv("REFRESH_KEY")
	refreshTime, err := time.ParseDuration(os.Getenv("REFRESH_TIME"))
	if err != nil {
		ENV.REFRESH_TIME = 7 * 24 * time.Hour
		fmt.Printf("Warning: Invalid REFRESH_TIME, using default: %v\n", ENV.REFRESH_TIME)
	} else {
		ENV.REFRESH_TIME = refreshTime
	}

	ENV.AppTZ, err = time.LoadLocation(os.Getenv("APP_TZ"))
	if err != nil {
		ENV.AppTZ, _ = time.LoadLocation("Asia/Ashgabat")
	}
	ENV.TZAddHours, err = time.ParseDuration(os.Getenv("TZ_ADD_HOURS"))
	if err != nil {
		ENV.TZAddHours = 0 * time.Second
	}

	ENV.GLE_KEY = os.Getenv("GLE_KEY")
	ENV.GLE_SECRET = os.Getenv("GLE_SECRET")
	ENV.GLE_CALLBACK = os.Getenv("GLE_CALLBACK")

	ENV.SMTP_HOST = os.Getenv("SMTP_HOST")
	ENV.SMTP_PORT = os.Getenv("SMTP_PORT")
	ENV.SMTP_MAIL = os.Getenv("SMTP_MAIL")
	ENV.SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	ENV.APP_LOGO_URL = fmt.Sprintf("%s/%s/assets/logo.svg", ENV.API_SERVER_URL, ENV.API_PREFIX)
	if len(os.Getenv("APP_LOGO_URL")) > 8 {
		ENV.APP_LOGO_URL = os.Getenv("APP_LOGO_URL")
	}

	ENV.FileUpload = FileUpload{
		MaxFileSize:      ENV.MAX_FILE_SIZE * 1024 * 1024,
		MaxFiles:         ENV.MAX_FILES_UPLOAD,
		AllowedMimeTypes: mergeAllowedTypes(),
		StorageBasePath:  ENV.UPLOAD_PATH,
	}

	ENV.GoogleOAuthConfig = oauth2.Config{
		ClientID:     ENV.GLE_KEY,
		ClientSecret: ENV.GLE_SECRET,
		RedirectURL:  ENV.GLE_CALLBACK,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

var (
	AllowedImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}

	AllowedVideoTypes = map[string]bool{
		"video/mp4":  true,
		"video/webm": true,
		"video/mpeg": true,
	}

	AllowedDocumentTypes = map[string]bool{
		"application/pdf":  true,
		"application/doc":  true,
		"application/docx": true,
		"text/plain":       true,
	}
	AllowedAudioTypes = map[string]bool{
		"audio/mpeg":     true, // .mp3
		"audio/wav":      true, // .wav
		"audio/x-wav":    true, // alternate .wav
		"audio/ogg":      true, // .ogg
		"audio/webm":     true, // .webm
		"audio/aac":      true, // .aac
		"audio/flac":     true, // .flac
		"audio/mp4":      true, // .m4a
		"audio/3gpp":     true, // .3gp
		"audio/x-ms-wma": true, // .wma
	}
	AllowedArchiveTypes = map[string]bool{
		"application/zip":              true,
		"application/x-rar-compressed": false,
		"application/x-7z-compressed":  false,
		"application/gzip":             false,
	}
	AllowedSpreadsheetTypes = map[string]bool{
		"application/vnd.ms-excel": false, // .xls
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": false, // .xlsx
	}
	AllowedPresentationTypes = map[string]bool{
		"application/vnd.ms-powerpoint":                                             false, // .ppt
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": false, // .pptx
	}
	AllowedCodeTypes = map[string]bool{
		"text/markdown":          false,
		"text/html":              false,
		"text/css":               false,
		"application/javascript": false,
		"application/json":       false,
		"text/x-python":          false,
	}
)

func mergeAllowedTypes() map[string]bool {
	allowedTypes := make(map[string]bool)
	for k, v := range AllowedImageTypes {
		allowedTypes[k] = v
	}
	for k, v := range AllowedVideoTypes {
		allowedTypes[k] = v
	}
	for k, v := range AllowedDocumentTypes {
		allowedTypes[k] = v
	}
	for k, v := range AllowedAudioTypes {
		allowedTypes[k] = v
	}
	for k, v := range AllowedArchiveTypes {
		allowedTypes[k] = v
	}
	for k, v := range AllowedSpreadsheetTypes {
		allowedTypes[k] = v
	}
	for k, v := range AllowedPresentationTypes {
		allowedTypes[k] = v
	}
	for k, v := range AllowedCodeTypes {
		allowedTypes[k] = v
	}
	return allowedTypes
}
