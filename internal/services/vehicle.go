package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func CreateVehicle(ctx *gin.Context) {
	var vehicle dto.VehicleCreate

	if err := ctx.ShouldBindJSON(&vehicle); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		vehicle.CompanyID = companyID
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicle,
		vehicle.CompanyID, vehicle.VehicleTypeID, vehicle.VehicleBrandID,
		vehicle.VehicleModelID, vehicle.YearOfIssue, vehicle.Mileage,
		vehicle.Numberplate, vehicle.TrailerNumberplate, vehicle.Gps,
		vehicle.Photo1URL, vehicle.Photo2URL, vehicle.Photo3URL,
		vehicle.Docs1URL, vehicle.Docs2URL, vehicle.Docs3URL,
		vehicle.ViewCount,
		vehicle.Meta,
		vehicle.Meta2,
		vehicle.Meta3,
		vehicle.Available,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created vehicle!", gin.H{"id": id}))
}

func UpdateVehicle(ctx *gin.Context) {
	id := ctx.Param("id")
	var vehicle dto.VehicleUpdate

	stmt := queries.UpdateVehicle

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		vehicle.CompanyID = &companyID
		stmt += ` WHERE (id = $1 AND company_id = $17) AND (active = 1 AND deleted = 0)`
	} else {
		stmt += ` WHERE id = $1`
	}
	stmt += ` RETURNING id;`

	if err := ctx.ShouldBindJSON(&vehicle); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		stmt,
		id, vehicle.VehicleTypeID, vehicle.VehicleBrandID,
		vehicle.VehicleModelID, vehicle.YearOfIssue, vehicle.Mileage,
		vehicle.Numberplate, vehicle.TrailerNumberplate, vehicle.Gps,
		vehicle.Photo1URL, vehicle.Photo2URL, vehicle.Photo3URL,
		vehicle.Docs1URL, vehicle.Docs2URL, vehicle.Docs3URL,
		vehicle.Active, vehicle.CompanyID, vehicle.Deleted,
		vehicle.ViewCount,
		vehicle.Meta,
		vehicle.Meta2,
		vehicle.Meta3,
		vehicle.Available,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated vehicle!", gin.H{"id": updatedID}))
}

func DeleteVehicle(ctx *gin.Context) {
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Operation can't be done by user", ""))
		return
	}

	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicle,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted vehicle!", gin.H{"id": id}))
}

func GetVehicleList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	var vehicles []dto.VehicleDetails
	var totalCount int
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&vehicles,
		queries.GetVehicleList,
		perPage,
		offset,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    vehicles,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle list", response))
}

func GetVehicle(ctx *gin.Context) {
	id := ctx.Param("id")

	var vehicle dto.VehicleDetails

	query := queries.GetVehicleByID
	err := pgxscan.Get(
		context.Background(),
		db.DB,
		&vehicle,
		query,
		id,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle not found", err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle details", vehicle))
}

func GetFilteredVehicles(ctx *gin.Context) {
	available := ctx.DefaultQuery("available", "")
	brandID := ctx.DefaultQuery("brand_id", "")
	typeID := ctx.DefaultQuery("type_id", "")
	modelID := ctx.DefaultQuery("model_id", "")
	mileageLess := ctx.DefaultQuery("mileage_less", "")
	companyID := ctx.DefaultQuery("company_id", "")
	numberplate := ctx.DefaultQuery("numberplate", "")

	query := `
		SELECT 
			v.*, 
			json_build_object(
				'id', c.id, 
				'company_name', c.company_name, 
				'country', c.country
			) as company,
			json_build_object(
				'id', vb.id, 
				'name', vb.name
			) as brand,
			json_build_object(
				'id', vm.id, 
				'name', vm.name
			) as model
		FROM tbl_vehicle v
		LEFT JOIN tbl_company c ON v.company_id = c.id
		LEFT JOIN tbl_vehicle_brand vb ON v.vehicle_brand_id = vb.id
		LEFT JOIN tbl_vehicle_model vm ON v.vehicle_model_id = vm.id
		WHERE v.deleted = 0
	`

	var conditions []string
	var args []interface{}
	argIdx := 1

	if available != "" {
		conditions = append(conditions, fmt.Sprintf("v.available = $%d", argIdx))
		args = append(args, available)
		argIdx++
	}
	if brandID != "" {
		conditions = append(conditions, fmt.Sprintf("v.vehicle_brand_id = $%d", argIdx))
		args = append(args, brandID)
		argIdx++
	}
	if typeID != "" {
		conditions = append(conditions, fmt.Sprintf("v.vehicle_type_id = $%d", argIdx))
		args = append(args, typeID)
		argIdx++
	}
	if modelID != "" {
		conditions = append(conditions, fmt.Sprintf("v.vehicle_model_id = $%d", argIdx))
		args = append(args, modelID)
		argIdx++
	}
	if mileageLess != "" {
		conditions = append(conditions, fmt.Sprintf("v.mileage < $%d", argIdx))
		args = append(args, mileageLess)
		argIdx++
	}
	if companyID != "" {
		conditions = append(conditions, fmt.Sprintf("v.company_id = $%d", argIdx))
		args = append(args, companyID)
		argIdx++
	}
	if numberplate != "" {
		conditions = append(conditions, fmt.Sprintf("v.numberplate LIKE $%d", argIdx))
		args = append(args, "%"+numberplate+"%") // Use wildcard for partial matching
		argIdx++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var vehicles []dto.VehicleDetails
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&vehicles,
		query,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Filtered vehicles", vehicles))
}

func SingleVehicleBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.GetVehicleBrand + " AND id = $1;"
	var brand []dto.VehicleBrand

	err := pgxscan.Select(
		context.Background(), db.DB,
		&brand, stmt, id,
	)

	if err != nil || len(brand) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle brand not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle brand", brand[0]))
}

func CreateVehicleBrand(ctx *gin.Context) {
	var brand dto.VehicleBrand

	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicleBrand,
		brand.Name,
		brand.Country,
		brand.FoundedYear,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle brand", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateVehicleBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	var brand dto.VehicleBrandUpdate

	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateVehicleBrand,
		id,
		brand.Name,
		brand.Country,
		brand.FoundedYear,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle brand", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

func DeleteVehicleBrand(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicleBrand,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle brand", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}

func SingleVehicleType(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.GetVehicleType + " AND id = $1;"
	var vehicleType []dto.VehicleType

	err := pgxscan.Select(
		context.Background(), db.DB,
		&vehicleType, stmt, id,
	)

	if err != nil || len(vehicleType) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle type not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle type", vehicleType[0]))
}

func CreateVehicleType(ctx *gin.Context) {
	var vehicleType dto.VehicleType

	if err := ctx.ShouldBindJSON(&vehicleType); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicleType,
		vehicleType.TitleEn,
		vehicleType.DescEn,
		vehicleType.TitleRu,
		vehicleType.DescRu,
		vehicleType.TitleTk,
		vehicleType.DescTk,
		vehicleType.TitleDe,
		vehicleType.DescDe,
		vehicleType.TitleAr,
		vehicleType.DescAr,
		vehicleType.TitleEs,
		vehicleType.DescEs,
		vehicleType.TitleFr,
		vehicleType.DescFr,
		vehicleType.TitleZh,
		vehicleType.DescZh,
		vehicleType.TitleJa,
		vehicleType.DescJa,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateVehicleType(ctx *gin.Context) {
	id := ctx.Param("id")
	var vehicleType dto.VehicleTypeUpdate

	if err := ctx.ShouldBindJSON(&vehicleType); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateVehicleType,
		id,
		vehicleType.TitleEn,
		vehicleType.DescEn,
		vehicleType.TitleRu,
		vehicleType.DescRu,
		vehicleType.TitleTk,
		vehicleType.DescTk,
		vehicleType.TitleDe,
		vehicleType.DescDe,
		vehicleType.TitleAr,
		vehicleType.DescAr,
		vehicleType.TitleEs,
		vehicleType.DescEs,
		vehicleType.TitleFr,
		vehicleType.DescFr,
		vehicleType.TitleZh,
		vehicleType.DescZh,
		vehicleType.TitleJa,
		vehicleType.DescJa,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

func DeleteVehicleType(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicleType,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}

func SingleVehicleModel(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.GetVehicleModel + " AND m.id = $1;"
	var model []dto.VehicleModel

	err := pgxscan.Select(
		context.Background(), db.DB,
		&model, stmt, id,
	)

	if err != nil || len(model) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle model not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle model", model[0]))
}

func CreateVehicleModel(ctx *gin.Context) {
	var model dto.VehicleModel

	if err := ctx.ShouldBindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicleModel,
		model.Name,
		model.Year,
		model.VehicleBrandID,
		model.VehicleTypeID,
		model.Feature,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle model", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateVehicleModel(ctx *gin.Context) {
	id := ctx.Param("id")
	var model dto.VehicleModelUpdate

	if err := ctx.ShouldBindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateVehicleModel,
		id,
		model.Name,
		model.Year,
		model.VehicleBrandID,
		model.VehicleTypeID,
		model.Feature,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle model", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}
func DeleteVehicleModel(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicleModel,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle model", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}

func GetVehicleModels(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	baseQuery := `
        SELECT 
            vm.*,
            COUNT(*) OVER() as total_count,
            json_build_object(
                'id', vb.id,
                'name', vb.name,
                'country', vb.country,
                'founded_year', vb.founded_year
            ) as brand,
            json_build_object(
                'id', vt.id,
                'title_en', vt.title_en,
                'title_ru', vt.title_ru,
                'title_tk', vt.title_tk
            ) as vehicle_type
        FROM tbl_vehicle_model vm
        LEFT JOIN tbl_vehicle_brand vb ON vm.vehicle_brand_id = vb.id
        LEFT JOIN tbl_vehicle_type vt ON vm.vehicle_type_id = vt.id
        WHERE vm.deleted = 0
    `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	// Add filters
	if brandID := ctx.Query("brand_id"); brandID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vm.vehicle_brand_id = $%d", argCounter))
		args = append(args, brandID)
		argCounter++
	}

	if typeID := ctx.Query("type_id"); typeID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vm.vehicle_type_id = $%d", argCounter))
		args = append(args, typeID)
		argCounter++
	}

	if year := ctx.Query("year"); year != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vm.year = $%d", argCounter))
		args = append(args, year)
		argCounter++
	}

	if searchTerm := ctx.Query("search"); searchTerm != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(vm.name ILIKE $%d OR vm.feature ILIKE $%d)",
			argCounter, argCounter))
		args = append(args, "%"+searchTerm+"%")
		argCounter++
	}

	if len(whereClauses) > 0 {
		baseQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY vm.id DESC LIMIT $%d OFFSET $%d",
		argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var models []dto.VehicleModelDetailed
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&models,
		baseQuery,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	var totalCount int
	if len(models) > 0 {
		totalCount = models[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    models,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle models", response))
}

func GetVehicleTypes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	baseQuery := `
        SELECT 
            vt.*,
            COUNT(*) OVER() as total_count,
            (
                SELECT json_agg(
                    json_build_object(
                        'id', vm.id,
                        'name', vm.name,
                        'year', vm.year
                    )
                )
                FROM tbl_vehicle_model vm 
                WHERE vm.vehicle_type_id = vt.id AND vm.deleted = 0
            ) as models
        FROM tbl_vehicle_type vt
        WHERE vt.deleted = 0
    `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	if searchTerm := ctx.Query("search"); searchTerm != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(
			"(vt.title_en ILIKE $%d OR vt.title_ru ILIKE $%d OR vt.title_tk ILIKE $%d)",
			argCounter, argCounter, argCounter))
		args = append(args, "%"+searchTerm+"%")
		argCounter++
	}

	if len(whereClauses) > 0 {
		baseQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY vt.id DESC LIMIT $%d OFFSET $%d",
		argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var types []dto.VehicleTypeDetailed
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&types,
		baseQuery,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	var totalCount int
	if len(types) > 0 {
		totalCount = types[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    types,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle types", response))
}

func GetVehicleBrands(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	baseQuery := `
        SELECT 
            vb.*,
            COUNT(*) OVER() as total_count,
            (
                SELECT json_agg(
                    json_build_object(
                        'id', vm.id,
                        'name', vm.name,
                        'year', vm.year
                    )
                )
                FROM tbl_vehicle_model vm 
                WHERE vm.vehicle_brand_id = vb.id AND vm.deleted = 0
            ) as models
        FROM tbl_vehicle_brand vb
        WHERE vb.deleted = 0
    `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	if country := ctx.Query("country"); country != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vb.country = $%d", argCounter))
		args = append(args, country)
		argCounter++
	}

	if yearFrom := ctx.Query("founded_year_from"); yearFrom != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vb.founded_year >= $%d", argCounter))
		args = append(args, yearFrom)
		argCounter++
	}

	if yearTo := ctx.Query("founded_year_to"); yearTo != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vb.founded_year <= $%d", argCounter))
		args = append(args, yearTo)
		argCounter++
	}

	if searchTerm := ctx.Query("search"); searchTerm != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vb.name ILIKE $%d", argCounter))
		args = append(args, "%"+searchTerm+"%")
		argCounter++
	}

	if len(whereClauses) > 0 {
		baseQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY vb.id DESC LIMIT $%d OFFSET $%d",
		argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var brands []dto.VehicleBrandDetailed
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&brands,
		baseQuery,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	var totalCount int
	if len(brands) > 0 {
		totalCount = brands[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    brands,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle brands", response))
}
