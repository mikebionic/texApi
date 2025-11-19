package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/joho/godotenv"
)

type Config struct {
	API_HOST       string
	API_SERVER_URL string
	API_PREFIX     string
	API_PORT       string
	API_SECRET     string
	JITSI_URL      string
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

	OTP_SERVICE_ROUTE      string
	OTP_SERVICE_TEXT       string
	OTP_ANDROID_HASH       string
	APP_NAME               string
	FIREBASE_ADMINSDK_FILE string

	FileUpload FileUpload
}

type FileUpload struct {
	MaxFileSize      int64 // in bytes (e.g., 10MB = 10 * 1024 * 1024)
	MaxFiles         int
	AllowedMimeTypes map[string]bool
	StorageBasePath  string
}

var ENV Config

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
		"application/pdf":    true,
		"application/msword": true, // .doc
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
		"text/plain": true,
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
	AllowedApkTypes = map[string]bool{
		"application/vnd.android.package-archive": true, // .apk
	}
)

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	ENV.API_HOST = getEnv("API_HOST", "localhost")
	ENV.API_SERVER_URL = getEnv("API_SERVER_URL", "http://localhost")
	ENV.API_PREFIX = getEnv("API_PREFIX", "api")
	ENV.API_PORT = getEnv("API_PORT", "8080")
	ENV.API_DEBUG = getEnvBool("API_DEBUG", false)
	ENV.API_SECRET = getEnv("API_SECRET", "")
	ENV.JITSI_URL = getEnv("JITSI_URL", "")
	ENV.SYSTEM_HEADER = getEnv("SYSTEM_HEADER", "")
	ENV.SYSTEM_SECRET = getEnv("SYSTEM_SECRET", "")

	if ENV.API_SECRET == "" {
		return fmt.Errorf("API_SECRET is required")
	}

	ENV.SESSION_MAX_AGE = getEnvInt("SESSION_MAX_AGE", 86400*30) // 30 days default

	ENV.UPLOAD_PATH = getEnv("UPLOAD_PATH", "./uploads")
	ENV.MAX_FILE_SIZE = int64(getEnvInt("MAX_FILE_SIZE", 10)) // MB
	ENV.MAX_FILES_UPLOAD = getEnvInt("MAX_FILES_UPLOAD", 5)
	ENV.STATIC_URL = fmt.Sprintf("/%s/uploads/", ENV.API_PREFIX)
	ENV.COMPRESS_IMAGES = getEnvInt("COMPRESS_IMAGES", 1)
	ENV.COMPRESS_SIZE = getEnvInt("COMPRESS_SIZE", 1920)
	ENV.COMPRESS_QUALITY = getEnvInt("COMPRESS_QUALITY", 85)

	if err := os.MkdirAll(ENV.UPLOAD_PATH, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	ENV.ENCRYPT_PASSWORDS = getEnvBool("ENCRYPT_PASSWORDS", true)

	ENV.DB_HOST = getEnv("DB_HOST", "localhost")
	ENV.DB_PORT = getEnv("DB_PORT", "5432")
	ENV.DB_USER = getEnv("DB_USER", "")
	ENV.DB_PASSWORD = getEnv("DB_PASSWORD", "")
	ENV.DB_NAME = getEnv("DB_NAME", "")

	if ENV.DB_USER == "" || ENV.DB_PASSWORD == "" || ENV.DB_NAME == "" {
		return fmt.Errorf("database credentials (DB_USER, DB_PASSWORD, DB_NAME) are required")
	}

	ENV.ACCESS_KEY = getEnv("ACCESS_KEY", "")
	ENV.REFRESH_KEY = getEnv("REFRESH_KEY", "")

	if ENV.ACCESS_KEY == "" || ENV.REFRESH_KEY == "" {
		return fmt.Errorf("ACCESS_KEY and REFRESH_KEY are required")
	}

	var err error
	ENV.ACCESS_TIME, err = time.ParseDuration(getEnv("ACCESS_TIME", "15m"))
	if err != nil {
		ENV.ACCESS_TIME = 15 * time.Minute
		fmt.Printf("Warning: Invalid ACCESS_TIME, using default: %v\n", ENV.ACCESS_TIME)
	}

	ENV.REFRESH_TIME, err = time.ParseDuration(getEnv("REFRESH_TIME", "168h")) // 7 дней
	if err != nil {
		ENV.REFRESH_TIME = 7 * 24 * time.Hour
		fmt.Printf("Warning: Invalid REFRESH_TIME, using default: %v\n", ENV.REFRESH_TIME)
	}

	ENV.AppTZ, err = time.LoadLocation(getEnv("APP_TZ", "Asia/Ashgabat"))
	if err != nil {
		ENV.AppTZ, _ = time.LoadLocation("Asia/Ashgabat")
		fmt.Printf("Warning: Invalid APP_TZ, using default: Asia/Ashgabat\n")
	}

	ENV.TZAddHours, err = time.ParseDuration(getEnv("TZ_ADD_HOURS", "0h"))
	if err != nil {
		ENV.TZAddHours = 0 * time.Second
	}

	ENV.GLE_KEY = getEnv("GLE_KEY", "")
	ENV.GLE_SECRET = getEnv("GLE_SECRET", "")
	ENV.GLE_CALLBACK = getEnv("GLE_CALLBACK", "")

	if ENV.GLE_KEY != "" && ENV.GLE_SECRET != "" {
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

	ENV.SMTP_HOST = getEnv("SMTP_HOST", "")
	ENV.SMTP_PORT = getEnv("SMTP_PORT", "587")
	ENV.SMTP_MAIL = getEnv("SMTP_MAIL", "")
	ENV.SMTP_PASSWORD = getEnv("SMTP_PASSWORD", "")

	customLogoURL := getEnv("APP_LOGO_URL", "")
	if customLogoURL != "" {
		ENV.APP_LOGO_URL = customLogoURL
	} else {
		ENV.APP_LOGO_URL = fmt.Sprintf("%s/%s/assets/logo.png", ENV.API_SERVER_URL, ENV.API_PREFIX)
	}

	ENV.OTP_SERVICE_ROUTE = getEnv("OTP_SERVICE_ROUTE", "")
	ENV.OTP_SERVICE_TEXT = getEnv("OTP_SERVICE_TEXT", "")
	ENV.OTP_ANDROID_HASH = getEnv("OTP_ANDROID_HASH", "")
	ENV.APP_NAME = getEnv("APP_NAME", "MyApp")
	ENV.FIREBASE_ADMINSDK_FILE = getEnv("FIREBASE_ADMINSDK_FILE", "")

	ENV.FileUpload = FileUpload{
		MaxFileSize:      ENV.MAX_FILE_SIZE * 1024 * 1024, // Convert MB to bytes
		MaxFiles:         ENV.MAX_FILES_UPLOAD,
		AllowedMimeTypes: mergeAllowedTypes(),
		StorageBasePath:  ENV.UPLOAD_PATH,
	}

	return nil
}

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
	for k, v := range AllowedApkTypes {
		allowedTypes[k] = v
	}

	return allowedTypes
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		fmt.Printf("Warning: Invalid integer value for %s, using default: %d\n", key, defaultValue)
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
