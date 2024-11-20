package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"os"
	"strings"
	"texApi/config"
	"time"
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
		return nil, errors.New("didn't upload the files")
	}

	files := form.File["files"]

	if len(files) == 0 {
		return nil, errors.New("must load minimum 1 file")
	}

	var filePaths []string
	var fileNames []string
	var fileCount = 0

	for _, file := range files {
		splitedFileName := strings.Split(file.Filename, ".")
		extension := splitedFileName[len(splitedFileName)-1]

		extensionExists := extensions[extension]
		if extensionExists == false {
			return nil, errors.New(fmt.Sprintf("This file is forbidden: %s", extension))
		}

		fileCount += 1
		if fileCount > config.ENV.MAX_FILE_UPLOAD_COUNT {
			return nil, errors.New("trying to upload too many files")
		}

		fileNames = append(fileNames, strings.ReplaceAll(splitedFileName[0]+"-"+uuid.NewString()+"."+extension, " ", "-"))
	}

	for index, file := range files {
		readerFile, _ := file.Open()

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, readerFile)
		dir, err := CreateTodayDir(config.ENV.UPLOAD_PATH)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(
			dir+fileNames[index],
			buf.Bytes(),
			os.ModePerm,
		)
		if err != nil {
			return nil, err
		}

		filePaths = append(filePaths, config.ENV.API_SERVER_URL+config.ENV.STATIC_URL+strings.ReplaceAll(dir, config.ENV.UPLOAD_PATH, "")+fileNames[index])
	}

	return filePaths, nil
}

func CreateTodayDir(absPath string) (directory string, err error) {
	currentDate := time.Now().Format("2006-01-02")
	directory = absPath + currentDate + "/"
	if _, err = os.Stat(directory); os.IsNotExist(err) {
		err = os.Mkdir(directory, os.ModePerm)
		if err != nil {
			return
		}
	}
	return
}
