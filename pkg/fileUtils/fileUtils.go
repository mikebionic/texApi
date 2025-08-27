package fileUtils

import (
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"texApi/config"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type FileValidationResult struct {
	File             *multipart.FileHeader
	ProcessedFile    ProcessedFile
	ValidationErrors []string
}

type ProcessedFile struct {
	UniqueFileName string
	StoragePath    string
	MediaType      string

	FilePath   string
	ThumbPath  string
	ThumbFn    string
	OriginalFn string
	MimeType   string
	FileSize   int64
	Duration   *int
	Width      int
	Height     int
}

func IsImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true, ".bmp": true,
	}
	return imageExts[ext]
}

func GenerateMediaURL(uuid, filename string) map[string]string {
	return map[string]string{
		"url":       strings.Join([]string{config.ENV.API_SERVER_URL, config.ENV.API_PREFIX, "media", uuid, filename}, "/"),
		"thumb_url": strings.Join([]string{config.ENV.API_SERVER_URL, config.ENV.API_PREFIX, "media", uuid, filename, "thumb"}, "/"),
	}
}

func SaveFile(fileHeader *multipart.FileHeader, processedFile *ProcessedFile) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	outFile, err := os.Create(processedFile.StoragePath)
	if err != nil {
		return fmt.Errorf("cannot create file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func ValidateSingleFile(fileHeader *multipart.FileHeader, categoryFN string) FileValidationResult {
	result := FileValidationResult{File: fileHeader}
	if fileHeader.Size > config.ENV.FileUpload.MaxFileSize {
		result.ValidationErrors = append(result.ValidationErrors, fmt.Sprintf("File too large. Max size: %d bytes", config.ENV.FileUpload.MaxFileSize))
		return result
	}

	file, err := fileHeader.Open()
	if err != nil {
		result.ValidationErrors = append(result.ValidationErrors, "Cannot open file")
		return result
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		result.ValidationErrors = append(result.ValidationErrors, "Cannot read file")
		return result
	}
	mimeType := DetectMimeType(buffer)

	if !config.ENV.FileUpload.AllowedMimeTypes[mimeType] {
		result.ValidationErrors = append(result.ValidationErrors, fmt.Sprintf("Unsupported file type: %s", mimeType))
		return result
	}

	ext := filepath.Ext(fileHeader.Filename)
	mediaType := DetermineMediaType(mimeType)
	uniqueFileName := GenerateUniqueFileName(fileHeader.Filename, ext)
	storagePath, filePath, err := GenerateStoragePath(config.ENV.FileUpload.StorageBasePath, categoryFN, mediaType, uniqueFileName)
	if err != nil {
		result.ValidationErrors = append(result.ValidationErrors, "Cannot generate storage path")
	}

	result.ProcessedFile = ProcessedFile{
		OriginalFn:     fileHeader.Filename,
		UniqueFileName: uniqueFileName,
		StoragePath:    storagePath,
		MediaType:      mediaType,

		FilePath:  filePath,
		ThumbPath: filepath.Join(filePath, "thumbnails"),
		ThumbFn:   "thumb_" + uniqueFileName,
		MimeType:  mimeType,
		FileSize:  fileHeader.Size,
	}

	return result
}

func ProcessMediaFiles(fileResults []FileValidationResult) ([]ProcessedFile, error) {
	var processedFiles []ProcessedFile
	var errorMessages []string

	for _, result := range fileResults {
		if len(result.ValidationErrors) > 0 {
			continue
		}

		processedFile := result.ProcessedFile
		var tempFile ProcessedFile
		var err error

		switch processedFile.MediaType {
		case "image":
			tempFile, err = ProcessImageFile(processedFile)
		case "video":
			tempFile, err = ProcessVideoFile(processedFile)
		case "audio":
			tempFile, err = ProcessAudioFile(processedFile)
		case "document":
			tempFile, err = ProcessDocumentFile(processedFile)
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("Unsupported media type: %s", processedFile.MediaType))
			continue
		}

		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error processing %s: %v", processedFile.OriginalFn, err))
			continue
		}

		processedFiles = append(processedFiles, tempFile)
	}

	if len(errorMessages) > 0 {
		return processedFiles, fmt.Errorf("some files failed to process: %s", strings.Join(errorMessages, "; "))
	}

	return processedFiles, nil
}

func CompressImageIfNeeded(imagePath string) error {
	if config.ENV.COMPRESS_IMAGES != 1 {
		return nil
	}

	img, err := imaging.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image for compression: %w", err)
	}

	origWidth := img.Bounds().Dx()
	origHeight := img.Bounds().Dy()

	maxSize := config.ENV.COMPRESS_SIZE
	if origWidth > maxSize || origHeight > maxSize {
		img = imaging.Fit(img, maxSize, maxSize, imaging.Lanczos)
		err = imaging.Save(img, imagePath, imaging.JPEGQuality(config.ENV.COMPRESS_QUALITY))
		if err != nil {
			return fmt.Errorf("failed to save compressed image: %w", err)
		}
	}

	return nil
}

func ProcessImageFile(processedFile ProcessedFile) (ProcessedFile, error) {
	if _, err := os.Stat(processedFile.StoragePath); os.IsNotExist(err) {
		return processedFile, fmt.Errorf("file does not exist: %s", processedFile.StoragePath)
	}

	err := CompressImageIfNeeded(processedFile.StoragePath)
	if err != nil {
		log.Printf("Warning: Failed to compress image: %s, error: %v", processedFile.StoragePath, err)
		// Continue processing even if compression fails
	}

	img, err := imaging.Open(processedFile.StoragePath)
	if err != nil {
		log.Printf("Failed to open image: %s, error: %v", processedFile.StoragePath, err)
		return processedFile, err
	}

	processedFile.Width = img.Bounds().Dx()
	processedFile.Height = img.Bounds().Dy()

	thumbnailPath := GenerateThumbPath(processedFile.StoragePath)
	thumbDir := filepath.Dir(thumbnailPath)
	os.MkdirAll(thumbDir, os.ModePerm)

	thumbnail := imaging.Fit(img, 300, 300, imaging.Lanczos)
	err = imaging.Save(thumbnail, thumbnailPath, imaging.JPEGQuality(75))

	if err != nil {
		log.Printf("Failed to save thumbnail: %s, error: %v", thumbnailPath, err)
		return processedFile, err
	}
	processedFile.ThumbPath = strings.TrimPrefix(thumbDir, config.ENV.UPLOAD_PATH)

	return processedFile, nil
}

func GenerateUniqueFileName(originalName string, ext string) string {
	baseName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' {
			return r
		}
		return '_'
	}, filepath.Base(originalName))

	if len(baseName) > 100 {
		baseName = baseName[:100]
	}
	uniqueID := uuid.New().String()[:12]
	return fmt.Sprintf("%s_%s%s", strings.TrimSuffix(baseName, ext), uniqueID, ext)
}

func GenerateStoragePath(baseDir string, categoryFN string, mediaType string, fileName string) (fullPath, filePath string, err error) {
	currentDate := time.Now().Format("2006-01-02")
	filePath = filepath.Join(categoryFN, mediaType, currentDate)
	dirPath := filepath.Join(baseDir, filePath)

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create directory: %s, error: %v", dirPath, err)
		return fullPath, filePath, err
	}

	return filepath.Join(dirPath, fileName), filePath, nil
}

func DetectMimeType(buffer []byte) string {
	// Using standard library to detect mime type
	return http.DetectContentType(buffer)
}

func DetermineMediaType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case strings.HasPrefix(mimeType, "application/") || strings.HasPrefix(mimeType, "text/"):
		return "document"
	default:
		return "unknown"
	}
}

func GenerateThumbPath(originalPath string) string {
	dir := filepath.Dir(originalPath)
	filename := filepath.Base(originalPath)
	return filepath.Join(dir, "thumbnails", "thumb_"+filename)
}

func GenerateVideoThumbnail(videoPath string) (string, error) {
	thumbnailPath := strings.TrimSuffix(GenerateThumbPath(videoPath), filepath.Ext(videoPath)) + ".jpg"
	thumbDir := filepath.Dir(thumbnailPath)

	if err := os.MkdirAll(thumbDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create thumbnail directory: %v", err)
	}

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-q:v", "2",
		thumbnailPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg thumbnail generation failed: %v, output: %s", err, string(output))
	}

	return thumbnailPath, nil
}

func ProcessVideoFile(processedFile ProcessedFile) (ProcessedFile, error) {
	if _, err := os.Stat(processedFile.StoragePath); os.IsNotExist(err) {
		return processedFile, fmt.Errorf("video file does not exist: %s", processedFile.StoragePath)
	}

	thumbnailPath, err := GenerateVideoThumbnail(processedFile.StoragePath)
	if err != nil {
		log.Printf("Failed to generate video thumbnail: %v", err)
	} else {
		processedFile.ThumbFn = "thumb_" + strings.TrimSuffix(processedFile.UniqueFileName, filepath.Ext(processedFile.UniqueFileName)) + ".jpg"
		thumbnailPath = strings.TrimSuffix(strings.TrimPrefix(thumbnailPath, config.ENV.UPLOAD_PATH), processedFile.ThumbFn)
	}

	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,duration",
		"-of", "csv=p=0",
		processedFile.StoragePath,
	)

	output, err := cmd.Output()
	if err == nil {
		values := strings.Split(strings.TrimSpace(string(output)), ",")
		if len(values) >= 2 {
			width, _ := strconv.Atoi(values[0])
			height, _ := strconv.Atoi(values[1])
			processedFile.Width = width
			processedFile.Height = height
		}
		if len(values) == 3 {
			duration, _ := strconv.ParseFloat(values[2], 64)
			durationInt := int(math.Round(duration))
			processedFile.Duration = &durationInt
		}
	}

	return processedFile, nil
}

func ProcessAudioFile(processedFile ProcessedFile) (ProcessedFile, error) {
	if _, err := os.Stat(processedFile.StoragePath); os.IsNotExist(err) {
		return processedFile, fmt.Errorf("audio file does not exist: %s", processedFile.StoragePath)
	}

	// Generate a thumbnail for audio files (optional, could use a default audio icon)
	// For demonstration, we'll create a waveform image using ffmpeg
	thumbnailPath := strings.TrimSuffix(GenerateThumbPath(processedFile.StoragePath), filepath.Ext(processedFile.StoragePath)) + ".jpg"
	thumbDir := filepath.Dir(thumbnailPath)

	if err := os.MkdirAll(thumbDir, os.ModePerm); err != nil {
		return processedFile, fmt.Errorf("failed to create thumbnail directory: %v", err)
	}

	// Generate waveform image as thumbnail
	cmd := exec.Command("ffmpeg",
		"-i", processedFile.StoragePath,
		"-filter_complex", "showwavespic=s=640x120:colors=#3498db",
		"-frames:v", "1",
		thumbnailPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Failed to generate audio waveform thumbnail: %v, output: %s", err, string(output))
		// Continue even if thumbnail generation fails
	} else {
		processedFile.ThumbFn = "thumb_" + strings.TrimSuffix(processedFile.UniqueFileName, filepath.Ext(processedFile.UniqueFileName)) + ".jpg"
	}

	// Get audio metadata using ffprobe
	cmd = exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		processedFile.StoragePath,
	)

	output, err := cmd.Output()
	if err == nil {
		durationStr := strings.TrimSpace(string(output))
		if durationStr != "" {
			duration, err := strconv.ParseFloat(durationStr, 64)
			if err == nil {
				durationInt := int(math.Round(duration))
				processedFile.Duration = &durationInt
			}
		}
	}

	return processedFile, nil
}

func ProcessDocumentFile(processedFile ProcessedFile) (ProcessedFile, error) {
	if _, err := os.Stat(processedFile.StoragePath); os.IsNotExist(err) {
		return processedFile, fmt.Errorf("document file does not exist: %s", processedFile.StoragePath)
	}
	return processedFile, nil
}
