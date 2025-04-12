package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"texApi/config"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/fileUtils"
	"texApi/pkg/utils"
)

func UploadFile(ctx *gin.Context) {

	filePaths, err := utils.SaveFiles(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error saving file", err.Error()))
		return
	}

	for k, filePath := range filePaths {
		filePaths[k] = config.ENV.API_SERVER_URL + filePath
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully uploaded", filePaths))
}

func ValidateAndProcessFiles(ctx *gin.Context, categoryFN, fileForm string) ([]fileUtils.FileValidationResult, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, errors.New("error parsing multipart form")
	}

	files := form.File[fileForm]
	if len(files) == 0 {
		return nil, errors.New("no files uploaded")
	}

	if len(files) > config.ENV.FileUpload.MaxFiles {
		return nil, fmt.Errorf("too many files. Max allowed: %d", config.ENV.FileUpload.MaxFiles)
	}

	var validationResults []fileUtils.FileValidationResult
	var failedFiles []string

	for _, fileHeader := range files {
		result := fileUtils.ValidateSingleFile(fileHeader, categoryFN)

		if len(result.ValidationErrors) > 0 {
			failedFiles = append(failedFiles, fileHeader.Filename)
			validationResults = append(validationResults, result)
			continue
		}

		err = fileUtils.SaveFile(fileHeader, &result.ProcessedFile)
		if err != nil {
			failedFiles = append(failedFiles, fileHeader.Filename)
			result.ValidationErrors = append(result.ValidationErrors, fmt.Sprintf("failed to save file: %v", err))
			validationResults = append(validationResults, result)
			continue
		}

		if fileUtils.IsImageFile(result.ProcessedFile.StoragePath) {
			err = fileUtils.CompressImageIfNeeded(result.ProcessedFile.StoragePath)
			if err != nil {
				result.ValidationErrors = append(result.ValidationErrors, fmt.Sprintf("image compression warning: %v", err))
			}
		}

		validationResults = append(validationResults, result)
	}

	if len(failedFiles) > 0 {
		return validationResults, fmt.Errorf("some files failed validation or saving: %v", failedFiles)
	}

	return validationResults, nil
}

func SaveMediaToDatabase(
	ctx context.Context,
	tx pgx.Tx,
	processedFile fileUtils.ProcessedFile,
	userID int,
	companyID int,
	context string,
	contextID int,
	isPrimary bool,
) (dto.MediaCreate, error) {
	query := `
       INSERT INTO tbl_media (
          user_id, company_id, media_type, context, context_id,
          filename, file_path, thumb_path, thumb_fn, original_fn,
          mime_type, file_size, duration, width, height
       ) VALUES (
          $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
       ) RETURNING id, uuid, user_id, media_type, context, filename, original_fn
    `

	var media dto.MediaCreate
	err := pgxscan.Get(ctx, tx, &media, query,
		userID,
		companyID,
		processedFile.MediaType,
		context,
		contextID,
		processedFile.UniqueFileName,
		processedFile.FilePath,
		processedFile.ThumbPath,
		processedFile.ThumbFn,
		processedFile.OriginalFn,
		processedFile.MimeType,
		processedFile.FileSize,
		processedFile.Duration,
		processedFile.Width,
		processedFile.Height,
	)

	if err != nil {
		return media, err
	}

	if context == "message" {
		_, err = tx.Exec(ctx, `
            INSERT INTO tbl_message_media (message_id, media_id, is_primary, sort_order)
            VALUES ($1, $2, $3, $4)
        `, contextID, media.ID, isPrimary, 0)
	}

	return media, err
}

func ExtractMediaIDs(mediaList []dto.MediaCreate) []map[string]interface{} {
	ids := make([]map[string]interface{}, len(mediaList))
	for i, media := range mediaList {
		ids[i] = map[string]interface{}{
			"id":   media.ID,
			"uuid": media.UUID,
		}
	}
	return ids
}

func MediaFileHandler(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	filename := ctx.Param("filename")
	isThumb := ctx.Param("thumb") == "thumb"

	media, filePath, mimeType, err := retrieveMediaInfo(ctx, uuid, filename, isThumb)
	if err != nil {
		return
	}

	isStreamingMedia := strings.HasPrefix(mimeType, "video/") ||
		strings.HasPrefix(mimeType, "audio/")

	if isStreamingMedia && !isThumb {
		serveStreamingMedia(ctx, filePath, mimeType, media.Filename)
	} else {
		serveRegularFile(ctx, filePath, mimeType, media.Filename)
	}
}

func retrieveMediaInfo(ctx *gin.Context, uuid, filename string, isThumb bool) (*dto.MediaMain, string, string, error) {
	query := `
       SELECT filename, file_path, thumb_fn, thumb_path, mime_type, media_type
       FROM tbl_media 
       WHERE uuid = $1 AND active = 1 AND deleted = 0
    `

	var media dto.MediaMain
	err := pgxscan.Get(context.Background(), db.DB, &media, query, uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Media not found", err.Error()))
		return nil, "", "", err
	}

	if media.Filename != filename {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Media not found", "Filename mismatch"))
		return nil, "", "", fmt.Errorf("filename mismatch")
	}

	var filePath string
	mimeType := *media.MimeType

	if isThumb {
		filePath = filepath.Join(config.ENV.UPLOAD_PATH, *media.ThumbPath, *media.ThumbFn)
		if media.MediaType != "image" {
			mimeType = "image/jpg"
		}
	} else {
		filePath = filepath.Join(config.ENV.UPLOAD_PATH, *media.FilePath, media.Filename)
	}

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("File not found", "File does not exist: "+err.Error()))
		return nil, "", "", err
	}

	return &media, filePath, mimeType, nil
}

func serveRegularFile(ctx *gin.Context, filePath, mimeType, filename string) {
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	ctx.Header("Content-Type", mimeType)
	ctx.File(filePath)
}

func serveStreamingMedia(ctx *gin.Context, filePath, mimeType, filename string) {
	file, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("File error", err.Error()))
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("File error", err.Error()))
		return
	}
	fileSize := fileInfo.Size()

	rangeHeader := ctx.GetHeader("Range")
	if rangeHeader != "" {
		servePartialContent(ctx, file, fileSize, rangeHeader, mimeType, filename)
		return
	}

	ctx.Header("Accept-Ranges", "bytes")
	ctx.Header("Content-Length", fmt.Sprintf("%d", fileSize))
	serveRegularFile(ctx, filePath, mimeType, filename)
}

func servePartialContent(ctx *gin.Context, file *os.File, fileSize int64, rangeHeader, mimeType, filename string) {
	ranges, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		ctx.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		ctx.AbortWithStatus(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	if len(ranges) > 1 {
		// We don't support multiple ranges
		ctx.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		ctx.AbortWithStatus(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	r := ranges[0]
	if r.start >= fileSize {
		ctx.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		ctx.AbortWithStatus(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	if r.end >= fileSize {
		r.end = fileSize - 1
	}

	file.Seek(r.start, io.SeekStart)

	contentLength := r.end - r.start + 1
	ctx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", r.start, r.end, fileSize))
	ctx.Header("Accept-Ranges", "bytes")
	ctx.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	ctx.Header("Content-Type", mimeType)
	ctx.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Description", "File Transfer")
	ctx.Status(http.StatusPartialContent)
	io.CopyN(ctx.Writer, file, contentLength)
}

type httpRange struct {
	start, end int64
}

func parseRange(rangeHeader string, size int64) ([]httpRange, error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, fmt.Errorf("invalid range header format")
	}

	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
	var ranges []httpRange

	for _, r := range strings.Split(rangeHeader, ",") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}

		var start, end int64
		var err error

		if parts[0] == "" {
			// suffix range: -N
			end = size - 1
			val, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid range value")
			}
			if val > size {
				start = 0
			} else {
				start = size - val
			}
		} else {
			// standard range: N-M or N-
			start, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid range start value")
			}

			if parts[1] == "" {
				// open-ended range: N-
				end = size - 1
			} else {
				// standard range: N-M
				end, err = strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid range end value")
				}
			}
		}

		if start > end || start < 0 || end >= size {
			continue // Skip invalid ranges
		}

		ranges = append(ranges, httpRange{start, end})
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("no valid ranges")
	}

	return ranges, nil
}
