package services

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func GetRequestBids(ctx *gin.Context) {
	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))

	stmt := queries.GetMyBids + " AND company_id = $1"
	var bids []dto.BidCreate

	err := pgxscan.Select(
		context.Background(), database.DB,
		&bids, stmt,
		companyID,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Error retrieving bids", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company bids", bids))
}

func GetUserBids(ctx *gin.Context) {

}
func CreateBid(ctx *gin.Context) {

}
func ChangeBidState(ctx *gin.Context) {

}
func DeleteUserBid(ctx *gin.Context) {

}
