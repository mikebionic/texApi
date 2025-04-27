package services

import (
	"fmt"
	"net/http"
	"texApi/internal/repo"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetContentTypes(ctx *gin.Context) {
	withContent, err := utils.HandleHeaderInt(ctx.GetHeader("WithContent"))
	langID, err := utils.HandleHeaderInt(ctx.GetHeader("LangID"))
	ctID, err := utils.HandleHeaderInt(ctx.GetHeader("ContentTypeID"))
	if err != nil {
		response := utils.FormatErrorResponse("Invalid header value", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	contentTypes, err := repo.GetContentTypes(withContent, langID, ctID)
	if err != nil {
		response := utils.FormatErrorResponse("Failed to retrieve content types", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if ctID > 0 {
		response := utils.FormatResponse(fmt.Sprintf("Content type %d retrieved successfully", ctID), contentTypes[0])
		ctx.JSON(http.StatusOK, response)
	} else {
		response := utils.FormatResponse("Content types retrieved successfully", contentTypes)
		ctx.JSON(http.StatusOK, response)
	}
}
