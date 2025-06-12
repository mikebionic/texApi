package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
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
		&cargo.Info, &cargo.Qty, &cargo.Weight, &cargo.WeightType, &cargo.Meta, &cargo.Meta2, &cargo.Meta3,
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

func GetDetailedCargoList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	validOrderColumns := map[string]bool{
		"id": true, "name": true, "qty": true, "weight": true,
		"created_at": true, "updated_at": true,
	}

	orderBy := ctx.DefaultQuery("order_by", "id")
	if !validOrderColumns[orderBy] {
		orderBy = "id"
	}
	orderDir := strings.ToUpper(ctx.DefaultQuery("order_dir", "DESC"))
	if orderDir != "ASC" && orderDir != "DESC" {
		orderDir = "DESC"
	}

	baseQuery := `
        SELECT 
            c.*,
            COUNT(*) OVER() as total_count,
            
            -- Company fields
            json_build_object(
                'id', comp.id,
                'company_name', comp.company_name,
                'email', comp.email,
                'phone', comp.phone,
                'address', comp.address,
                'country', comp.country
            ) as company,
            
            -- Vehicle type fields
            json_build_object(
                'id', vt.id,
                'title_en', vt.title_en,
                'title_ru', vt.title_ru,
                'title_tk', vt.title_tk
            ) as vehicle_type,
            
            -- Packaging type fields
            json_build_object(
                'id', pt.id,
                'name_en', pt.name_en,
                'name_ru', pt.name_ru,
                'name_tk', pt.name_tk,
                'material', pt.material,
                'dimensions', pt.dimensions
            ) as packaging_type
            
        FROM tbl_cargo c
        LEFT JOIN tbl_company comp ON c.company_id = comp.id
        LEFT JOIN tbl_vehicle_type vt ON c.vehicle_type_id = vt.id
        LEFT JOIN tbl_packaging_type pt ON c.packaging_type_id = pt.id
    `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		whereClauses = append(whereClauses, "c.deleted = 0")
	}

	filters := map[string]string{
		"cargo_id":          ctx.Query("cargo_id"),
		"company_id":        ctx.Query("company_id"),
		"vehicle_type_id":   ctx.Query("vehicle_type_id"),
		"packaging_type_id": ctx.Query("packaging_type_id"),
		"weight_type":       ctx.Query("weight_type"),
		"gps":               ctx.Query("gps"),
	}

	for key, value := range filters {
		if value != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("c.%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	numericRanges := map[string]struct {
		min string
		max string
	}{
		"weight": {ctx.Query("min_weight"), ctx.Query("max_weight")},
		"qty":    {ctx.Query("min_qty"), ctx.Query("max_qty")},
	}

	for field, ranges := range numericRanges {
		if ranges.min != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("c.%s >= $%d", field, argCounter))
			minVal, _ := strconv.Atoi(ranges.min)
			args = append(args, minVal)
			argCounter++
		}
		if ranges.max != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("c.%s <= $%d", field, argCounter))
			maxVal, _ := strconv.Atoi(ranges.max)
			args = append(args, maxVal)
			argCounter++
		}
	}

	searchTerm := ctx.Query("search")
	if searchTerm != "" {
		searchClause := fmt.Sprintf(`(
            c.name ILIKE $%d OR 
            c.description ILIKE $%d OR 
            c.info ILIKE $%d
        )`, argCounter, argCounter, argCounter)
		whereClauses = append(whereClauses, searchClause)
		args = append(args, "%"+searchTerm+"%")
		argCounter++
	}

	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY c.%s %s LIMIT $%d OFFSET $%d",
		orderBy, orderDir, argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var cargos []dto.CargoDetailed
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&cargos,
		query,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Couldn't retrieve data", err.Error()))
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

	ctx.JSON(http.StatusOK, utils.FormatResponse("Cargo list detailed", response))
}
