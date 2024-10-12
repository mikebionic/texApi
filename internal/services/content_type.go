package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	repo "texApi/internal/repositories"
	"texApi/pkg/utils"
)

func GetContentTypes(ctx *gin.Context) {
	withContent, err := strconv.Atoi(ctx.GetHeader("WithContent"))
	langID, err := strconv.Atoi(ctx.GetHeader("LangID"))
	if err != nil {
		response := utils.FormatErrorResponse("Invalid header value", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	contentTypes, err := repo.GetContentTypes(withContent, langID)
	if err != nil {
		response := utils.FormatErrorResponse("Failed to retrieve content types", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.FormatResponse("Content types retrieved successfully", contentTypes)
	ctx.JSON(http.StatusOK, response)
}
