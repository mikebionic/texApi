package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"texApi/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var extensions map[string]bool = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"svg":  true,
	"webp": true,
	"mp4":  true,
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
		ctx.JSON(400, gin.H{"message": "Must load minimum 1 file"})
		return nil, errors.New("Didn't upload the files")
	}

	files := form.File["files"]

	if len(files) == 0 {
		ctx.JSON(400, gin.H{"message": "Must load minimum 1 file"})
		return nil, errors.New("Must load minimum 1 file")
	}

	var filePaths []string
	var fileNames []string
	var video = 0
	var images = 0

	for _, file := range files {
		splitedFileName := strings.Split(file.Filename, ".")
		extension := splitedFileName[len(splitedFileName)-1]

		extensionExists := extensions[extension]

		if extensionExists == false {
			ctx.JSON(400, gin.H{"message": "This file is forbidden"})
			return nil, errors.New(
				fmt.Sprintf("Trying to upload %v file", extension),
			)
		}

		if extension == "mp4" {
			video += 1
		} else {
			images += 1
		}

		if video > 1 || images > 5 {
			ctx.JSON(400, gin.H{"message": "Only 5 images and 1 video"})
			return nil, errors.New(
				fmt.Sprintf(
					"Trying to upload %v video and %v images", video, images,
				),
			)
		}

		fileNames = append(fileNames, uuid.NewString()+"."+extension)
	}

	for index, file := range files {
		readerFile, _ := file.Open()

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, readerFile)
		os.WriteFile(
			config.ENV.UPLOAD_PATH+"orders/"+fileNames[index],
			buf.Bytes(),
			os.ModePerm,
		)

		filePaths = append(filePaths, "/uploads/orders/"+fileNames[index])
	}

	return filePaths, nil
}
