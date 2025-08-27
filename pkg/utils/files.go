package utils

import (
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"texApi/config"
	"time"

	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTodayDir(absPath string) (string, error) {
	currentDate := time.Now().Format("2006-01-02")
	directory := filepath.Join(absPath, currentDate)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return "", err
		}
	}
	return directory + string(os.PathSeparator), nil
}

const (
	TypeImage = "image"
	TypeVideo = "video"
	TypeDoc   = "document"
	TypeOther = "other"
)

type FileConfig struct {
	ImageExtensions map[string]bool
	VideoExtensions map[string]bool
	DocExtensions   map[string]bool
	OtherExtensions map[string]bool
	WebPQuality     float32
	MaxFileSize     int64 // in bytes
}

var fileConfig = FileConfig{
	ImageExtensions: map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"webp": true,
	},
	VideoExtensions: map[string]bool{
		"mp4":  true,
		"webm": true,
	},
	DocExtensions: map[string]bool{
		"pdf":  true,
		"docx": true,
		"pptx": true,
		"svg":  true,
		"gif":  true,
	},
	OtherExtensions: map[string]bool{
		"apk": true,
		"zip": true,
	},
	WebPQuality: 85,
	MaxFileSize: config.ENV.MAX_FILE_SIZE * 1024 * 1024,
}

func getFileType(extension string) string {
	ext := strings.ToLower(extension)

	if fileConfig.ImageExtensions[ext] {
		return TypeImage
	}
	if fileConfig.VideoExtensions[ext] {
		return TypeVideo
	}
	if fileConfig.DocExtensions[ext] {
		return TypeDoc
	}
	if fileConfig.OtherExtensions[ext] {
		return TypeOther
	}

	return ""
}

func isAllowedExtension(extension string) bool {
	return getFileType(extension) != ""
}

func generateUniqueFileName(originalName string) (string, string) {
	parts := strings.Split(originalName, ".")
	if len(parts) < 2 {
		return "", ""
	}

	baseName := strings.Join(parts[:len(parts)-1], ".")
	extension := strings.ToLower(parts[len(parts)-1])

	baseName = strings.ReplaceAll(baseName, " ", "-")
	uniqueName := fmt.Sprintf("%s-%s", baseName, uuid.NewString())

	return uniqueName, extension
}

func compressImageToWebP(file multipart.File, outputPath string) error {
	file.Seek(0, 0)

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	options := &webp.Options{Quality: fileConfig.WebPQuality}
	if err := webp.Encode(outFile, img, options); err != nil {
		return fmt.Errorf("failed to encode to WebP: %w", err)
	}

	return nil
}

func saveRegularFile(file multipart.File, outputPath string) error {
	file.Seek(0, 0)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

func processFile(fileHeader *multipart.FileHeader, targetDir string) (string, error) {
	uniqueName, extension := generateUniqueFileName(fileHeader.Filename)
	if extension == "" {
		return "", errors.New("invalid file name or extension")
	}

	if !isAllowedExtension(extension) {
		return "", fmt.Errorf("file extension '%s' is not allowed", extension)
	}

	if fileHeader.Size > fileConfig.MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", fileConfig.MaxFileSize)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	var finalFileName string
	var outputPath string

	fileType := getFileType(extension)

	if fileType == TypeImage {
		finalFileName = uniqueName + ".webp"
		outputPath = filepath.Join(targetDir, finalFileName)

		if err = compressImageToWebP(file, outputPath); err != nil {
			return "", fmt.Errorf("failed to compress image: %w", err)
		}
	} else {
		finalFileName = uniqueName + "." + extension
		outputPath = filepath.Join(targetDir, finalFileName)

		if err = saveRegularFile(file, outputPath); err != nil {
			return "", fmt.Errorf("failed to save file: %w", err)
		}
	}

	return finalFileName, nil
}

func SaveFiles(ctx *gin.Context) ([]string, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, errors.New("failed to parse multipart form")
	}

	if form == nil {
		return nil, errors.New("no files uploaded")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return nil, errors.New("must upload at least 1 file")
	}

	if len(files) > config.ENV.MAX_FILES_UPLOAD {
		return nil, fmt.Errorf("too many files: maximum %d allowed", config.ENV.MAX_FILES_UPLOAD)
	}

	targetDir, err := CreateTodayDir(config.ENV.UPLOAD_PATH + "UploadFile_route")
	if err != nil {
		return nil, fmt.Errorf("failed to create target directory: %w", err)
	}

	var filePaths []string

	for _, fileHeader := range files {
		fileName, err := processFile(fileHeader, targetDir)
		if err != nil {
			return nil, fmt.Errorf("failed to process file '%s': %w", fileHeader.Filename, err)
		}

		relativePath := strings.ReplaceAll(targetDir, config.ENV.UPLOAD_PATH, "")
		publicURL := config.ENV.STATIC_URL + relativePath + fileName
		filePaths = append(filePaths, publicURL)
	}

	return filePaths, nil
}

func WriteImage(ctx *gin.Context, dir string) (string, error) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		return "", errors.New("no image file provided")
	}
	defer file.Close()

	uniqueName, extension := generateUniqueFileName(header.Filename)
	if extension == "" {
		return "", errors.New("invalid image file name")
	}

	if getFileType(extension) != TypeImage {
		return "", errors.New("uploaded file is not a valid image")
	}

	targetDir := filepath.Join(config.ENV.UPLOAD_PATH, dir)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fileName := uniqueName + ".webp"
	outputPath := filepath.Join(targetDir, fileName)

	if err := compressImageToWebP(file, outputPath); err != nil {
		return "", fmt.Errorf("failed to process image: %w", err)
	}

	return fileName, nil
}
