package services

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func CreateRequest(ctx *gin.Context) {

}

func UpdateRequest(ctx *gin.Context) {

}

// User Company specific request
// TODO: take userID and validate in query
func GetCompanyRequests(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.GetHeader("UserID"))
	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))

	stmt := queries.GetMyRequest + " AND company_id = $1 AND user_id = $2;"
	var myRequests []dto.RequestCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&myRequests, stmt,
		companyID,
		userID,
	)

	if err != nil || len(myRequests) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Requests not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("My Requests", myRequests))
}

func GetRequests(ctx *gin.Context) {
	stmt := queries.GetMyRequest
	var allRequests []dto.RequestCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&allRequests, stmt,
	)

	if err != nil || len(allRequests) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Requests not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Requests", allRequests))
}

func DeleteRequest(ctx *gin.Context) {

}
