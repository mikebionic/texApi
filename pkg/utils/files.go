package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"strings"
	"texApi/config"
)

var extensions map[string]bool = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"svg":  true,
	"webp": true,
	"mp4":  true,
	"webm": true,
	"pdf":  true,
	"docx": true,
	"pptx": true,
}

func WriteImage(ctx *gin.Context, dir string) string {
	image, header, _ := ctx.Request.FormFile("image")

	if image == nil {
		return ""
	}

	splitedFileName := strings.Split(header.Filename, ".")
	extension := splitedFileName[len(splitedFileName)-1]

	if extension == "webp" || extension == "svg" || extension == "jpeg" ||
		extension == "jpg" || extension == "png" {

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, image)
		os.WriteFile(
			config.ENV.UPLOAD_PATH+dir+header.Filename,
			buf.Bytes(), os.ModePerm,
		)

		return header.Filename
	}

	return ""
}

func SaveFiles(ctx *gin.Context) ([]string, error) {
	form, _ := ctx.MultipartForm()

	if form == nil {
		ctx.JSON(http.StatusBadRequest, FormatErrorResponse("Must load minimum 1 file", ""))
		return nil, errors.New("Didn't upload the files")
	}

	files := form.File["files"]

	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, FormatErrorResponse("Must load minimum 1 file", ""))
		return nil, errors.New("Must load minimum 1 file")
	}

	var filePaths []string
	var fileNames []string
	var fileCount = 0

	for _, file := range files {
		splitedFileName := strings.Split(file.Filename, ".")
		extension := splitedFileName[len(splitedFileName)-1]

		extensionExists := extensions[extension]

		if extensionExists == false {
			ctx.JSON(http.StatusBadRequest, FormatErrorResponse("This file is forbidden", ""))
			return nil, errors.New(
				fmt.Sprintf("Trying to upload %v file", extension),
			)
		}

		fileCount += 1

		if fileCount > config.ENV.MAX_FILE_UPLOAD_COUNT {
			ctx.JSON(http.StatusBadRequest, FormatErrorResponse("Trying to upload too many files", ""))
			return nil, errors.New(fmt.Sprintf("Trying to upload %v files", fileCount))
		}

		fileNames = append(fileNames, uuid.NewString()+"."+extension)
	}

	for index, file := range files {
		readerFile, _ := file.Open()

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, readerFile)
		os.WriteFile(
			config.ENV.UPLOAD_PATH+"files/"+fileNames[index],
			buf.Bytes(),
			os.ModePerm,
		)

		filePaths = append(filePaths, "/uploads/files/"+fileNames[index])
	}

	return filePaths, nil
}
