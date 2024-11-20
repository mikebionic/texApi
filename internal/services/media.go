package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"texApi/pkg/utils"
)

func UploadFile(ctx *gin.Context) {

	filePaths, err := utils.SaveFiles(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Error saving file", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully uploaded", filePaths))
}
