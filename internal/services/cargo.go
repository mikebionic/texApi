package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func GetCargoList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	stmt := queries.GetCargoList

	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))
	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		stmt += ` WHERE (c.company_id = $3 OR $3 = 0) AND c.deleted = 0`
	} else {
		stmt += ` WHERE (c.company_id = $3 OR $3 = 0)`
	}
	stmt += ` ORDER BY c.id DESC LIMIT $1 OFFSET $2;`

	var cargos []dto.Cargo

	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&cargos,
		stmt,
		perPage,
		offset,
		companyID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Couldn't retrieve data", err.Error()))
		return
	}

	var totalCount int
	if len(cargos) > 0 {
		totalCount = cargos[0].TotalCount
	}
	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    cargos,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Cargo list", response))
}

func GetCargo(ctx *gin.Context) {
	id := ctx.Param("id")

	var cargo dto.Cargo

	stmt := queries.GetCargoByID
	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		stmt += ` AND c.deleted = 0;`
	}

	err := db.DB.QueryRow(
		context.Background(),
		queries.GetCargoByID,
		id,
	).Scan(
		&cargo.ID, &cargo.UUID, &cargo.CompanyID, &cargo.Name, &cargo.Description,
		&cargo.Info, &cargo.Qty, &cargo.Weight, &cargo.Meta, &cargo.Meta2, &cargo.Meta3,
		&cargo.VehicleTypeID, &cargo.PackagingTypeID, &cargo.GPS, &cargo.Photo1URL,
		&cargo.Photo2URL, &cargo.Photo3URL, &cargo.Docs1URL, &cargo.Docs2URL,
		&cargo.Docs3URL, &cargo.Note, &cargo.CreatedAt, &cargo.UpdatedAt,
		&cargo.Active, &cargo.Deleted,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Cargo not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Cargo details", cargo))
}

func CreateCargo(ctx *gin.Context) {
	var cargo dto.Cargo

	if err := ctx.ShouldBindJSON(&cargo); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		cargo.CompanyID = ctx.MustGet("companyID").(int)
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateCargo,
		cargo.CompanyID, cargo.Name, cargo.Description, cargo.Info, cargo.Qty,
		cargo.Weight, cargo.Meta, cargo.Meta2, cargo.Meta3, cargo.VehicleTypeID,
		cargo.PackagingTypeID, cargo.GPS, cargo.Photo1URL, cargo.Photo2URL,
		cargo.Photo3URL, cargo.Docs1URL, cargo.Docs2URL, cargo.Docs3URL, cargo.Note, cargo.WeightType,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating cargo", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created cargo!", gin.H{"id": id}))
}

func UpdateCargo(ctx *gin.Context) {
	id := ctx.Param("id")
	var cargo dto.CargoUpdate

	if err := ctx.ShouldBindJSON(&cargo); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	stmt := queries.UpdateCargo

	// TODO: only admin can make active or restore from deleted
	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		companyID := ctx.MustGet("companyID").(int)
		stmt += fmt.Sprintf(` AND company_id =%d AND deleted = 0`, companyID)
	}
	result, err := db.DB.Exec(
		context.Background(),
		stmt,
		id, cargo.Name, cargo.Description, cargo.Info, cargo.Qty,
		cargo.Weight, cargo.Meta, cargo.Meta2, cargo.Meta3, cargo.VehicleTypeID,
		cargo.PackagingTypeID, cargo.GPS, cargo.Photo1URL, cargo.Photo2URL,
		cargo.Photo3URL, cargo.Docs1URL, cargo.Docs2URL, cargo.Docs3URL, cargo.Note,
		cargo.Active, cargo.Deleted, cargo.WeightType,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating cargo", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Cargo not found or no changes were made", ""))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully updated cargo!", id))
}

func DeleteCargo(ctx *gin.Context) {
	id := ctx.Param("id")

	stmt := queries.DeleteCargo

	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		companyID := ctx.MustGet("companyID").(int)
		stmt += fmt.Sprintf(` AND company_id = %d`, companyID)
	}

	result, err := db.DB.Exec(
		context.Background(),
		stmt,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting cargo", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Cargo not found or no changes were made", ""))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully deleted cargo!", gin.H{"id": id}))
}
