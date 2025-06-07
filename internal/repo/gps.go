package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"texApi/pkg/utils"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
)

const (
	DefaultLimit    = 100
	DefaultOrderBy  = "id"
	DefaultOrderDir = "DESC"
)

func CreateTrip(input dto.StartTripInput) (int64, error) {
	ctx := context.Background()
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var offerData struct {
		FromAddress *string `db:"from_address"`
		ToAddress   *string `db:"to_address"`
		FromCountry *string `db:"from_country"`
		ToCountry   *string `db:"to_country"`
		DriverID    int     `db:"driver_id"`
		VehicleID   int     `db:"vehicle_id"`
		DistanceKM  *int    `db:"distance"`
	}

	mainOfferID := input.Offers[0].OfferID
	err = pgxscan.Get(ctx, tx, &offerData,
		`SELECT from_address, to_address, from_country, to_country, 
		        driver_id, vehicle_id, distance
		 FROM tbl_offer WHERE id = $1 AND deleted = 0`,
		mainOfferID)
	if err != nil {
		return 0, fmt.Errorf("failed to get offer data: %w", err)
	}

	driverID := getIntValue(input.DriverID, offerData.DriverID)
	vehicleID := getIntValue(input.VehicleID, offerData.VehicleID)
	fromAddress := getStringValue(input.FromAddress, offerData.FromAddress)
	toAddress := getStringValue(input.ToAddress, offerData.ToAddress)
	fromCountry := getStringValue(input.FromCountry, offerData.FromCountry)
	toCountry := getStringValue(input.ToCountry, offerData.ToCountry)

	var distanceKM *float64
	if input.DistanceKM != nil {
		distanceKM = input.DistanceKM
	} else if offerData.DistanceKM != nil {
		distance := float64(*offerData.DistanceKM)
		distanceKM = &distance
	}

	startDate := input.StartDate
	if startDate == nil {
		now := time.Now()
		startDate = &now
	}

	var tripID int64
	err = pgxscan.Get(ctx, tx, &tripID,
		`INSERT INTO tbl_trip (
			driver_id, vehicle_id, from_address, to_address, 
			from_country, to_country, start_date, end_date,
			from_location, to_location, distance_km, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'active')
		RETURNING id`,
		driverID, vehicleID, fromAddress, toAddress,
		fromCountry, toCountry, startDate, input.EndDate,
		input.FromLocation, input.ToLocation, distanceKM)
	if err != nil {
		return 0, fmt.Errorf("failed to create trip: %w", err)
	}

	for _, offer := range input.Offers {
		_, err = tx.Exec(ctx,
			`INSERT INTO tbl_offer_trip (trip_id, offer_id, is_main, status)
			 VALUES ($1, $2, $3, 'active')`,
			tripID, offer.OfferID, offer.IsMain)
		if err != nil {
			return 0, fmt.Errorf("failed to link offer to trip: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tripID, nil
}

func EndTrip(input dto.EndTripInput) error {
	ctx := context.Background()
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var logs []dto.GPSLog
	err = pgxscan.Select(ctx, tx, &logs,
		`SELECT * FROM tbl_gps_log 
		 WHERE trip_id = $1 
		 ORDER BY log_dt ASC`,
		input.ID)
	if err != nil {
		return fmt.Errorf("failed to get GPS logs: %w", err)
	}

	// Convert logs to JSON
	logsJSON, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to marshal GPS logs: %w", err)
	}

	// Update trip
	result, err := tx.Exec(ctx,
		`UPDATE tbl_trip 
		 SET status = 'completed', end_date = CURRENT_TIMESTAMP, 
		     gps_logs = $1, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $2 AND driver_id = $3 AND deleted = 0`,
		string(logsJSON), input.ID, input.DriverID)
	if err != nil {
		return fmt.Errorf("failed to end trip: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("trip not found or access denied")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type TripScan struct {
	ID              int64      `db:"id"`
	DriverID        int        `db:"driver_id"`
	VehicleID       int        `db:"vehicle_id"`
	FromAddress     *string    `db:"from_address"`
	ToAddress       *string    `db:"to_address"`
	FromCountry     *string    `db:"from_country"`
	ToCountry       *string    `db:"to_country"`
	StartDate       *time.Time `db:"start_date"`
	EndDate         *time.Time `db:"end_date"`
	FromLocationTxt *string    `db:"from_location_txt"` // ST_AsText result
	ToLocationTxt   *string    `db:"to_location_txt"`   // ST_AsText result
	DistanceKM      *float64   `db:"distance_km"`
	Status          string     `db:"status"`
	Meta            string     `db:"meta"`
	Meta2           string     `db:"meta2"`
	Meta3           string     `db:"meta3"`
	GPSLogs         string     `db:"gps_logs"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	Deleted         int        `db:"deleted"`
}

// TripScan to dto.Trip
func (ts *TripScan) ToTrip() dto.Trip {
	trip := dto.Trip{
		ID:          ts.ID,
		DriverID:    ts.DriverID,
		VehicleID:   ts.VehicleID,
		FromAddress: ts.FromAddress,
		ToAddress:   ts.ToAddress,
		FromCountry: ts.FromCountry,
		ToCountry:   ts.ToCountry,
		StartDate:   ts.StartDate,
		EndDate:     ts.EndDate,
		DistanceKM:  ts.DistanceKM,
		Status:      ts.Status,
		Meta:        ts.Meta,
		Meta2:       ts.Meta2,
		Meta3:       ts.Meta3,
		GPSLogs:     ts.GPSLogs,
		CreatedAt:   ts.CreatedAt,
		UpdatedAt:   ts.UpdatedAt,
		Deleted:     ts.Deleted,
	}

	// Parse from_location
	if ts.FromLocationTxt != nil && *ts.FromLocationTxt != "" {
		var point dto.Point
		if err := point.Scan(*ts.FromLocationTxt); err == nil {
			trip.FromLocation = &point
		}
	}

	// Parse to_location
	if ts.ToLocationTxt != nil && *ts.ToLocationTxt != "" {
		var point dto.Point
		if err := point.Scan(*ts.ToLocationTxt); err == nil {
			trip.ToLocation = &point
		}
	}

	return trip
}

// GPSLog struct for scanning
type GPSLogScan struct {
	ID             int64     `db:"id"`
	CompanyID      *int      `db:"company_id"`
	VehicleID      int       `db:"vehicle_id"`
	DriverID       int       `db:"driver_id"`
	OfferID        *int      `db:"offer_id"`
	TripID         *int      `db:"trip_id"`
	BatteryLevel   *int      `db:"battery_level"`
	Speed          *float64  `db:"speed"`
	Heading        *float64  `db:"heading"`
	Accuracy       *float64  `db:"accuracy"`
	CoordinatesTxt string    `db:"coordinates_txt"` // ST_AsText result
	Status         string    `db:"status"`
	LogDt          time.Time `db:"log_dt"`
	CreatedAt      time.Time `db:"created_at"`
}

// GPSLogScan to dto.GPSLog
func (gs *GPSLogScan) ToGPSLog() dto.GPSLog {
	log := dto.GPSLog{
		ID:           gs.ID,
		CompanyID:    gs.CompanyID,
		VehicleID:    gs.VehicleID,
		DriverID:     gs.DriverID,
		OfferID:      gs.OfferID,
		TripID:       gs.TripID,
		BatteryLevel: gs.BatteryLevel,
		Speed:        gs.Speed,
		Heading:      gs.Heading,
		Accuracy:     gs.Accuracy,
		Status:       gs.Status,
		LogDt:        gs.LogDt,
		CreatedAt:    gs.CreatedAt,
	}

	// Parse coordinates
	if gs.CoordinatesTxt != "" {
		if err := log.Coordinates.Scan(gs.CoordinatesTxt); err != nil {
			// If parsing fails, set to zero coordinates
			log.Coordinates = dto.Point{Lat: 0, Lng: 0}
		}
	}

	return log
}

func GetTrips(query dto.TripQuery) ([]dto.Trip, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "deleted = 0")

	if query.DriverID != nil {
		conditions = append(conditions, fmt.Sprintf("driver_id = $%d", argIndex))
		args = append(args, *query.DriverID)
		argIndex++
	}

	if query.VehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("vehicle_id = $%d", argIndex))
		args = append(args, *query.VehicleID)
		argIndex++
	}

	if query.FromAddress != nil {
		conditions = append(conditions, fmt.Sprintf("from_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.FromAddress+"%")
		argIndex++
	}

	if query.ToAddress != nil {
		conditions = append(conditions, fmt.Sprintf("to_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.ToAddress+"%")
		argIndex++
	}

	if query.FromCountry != nil {
		conditions = append(conditions, fmt.Sprintf("from_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.FromCountry+"%")
		argIndex++
	}

	if query.ToCountry != nil {
		conditions = append(conditions, fmt.Sprintf("to_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.ToCountry+"%")
		argIndex++
	}

	if query.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIndex))
		args = append(args, *query.StartDate)
		argIndex++
	}

	if query.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("end_date <= $%d", argIndex))
		args = append(args, *query.EndDate)
		argIndex++
	}

	if query.DistanceKM != nil {
		conditions = append(conditions, fmt.Sprintf("distance_km >= $%d", argIndex))
		args = append(args, *query.DistanceKM)
		argIndex++
	}

	if query.TripOfferID != nil {
		conditions = append(conditions, fmt.Sprintf(`id IN (
			SELECT trip_id FROM tbl_offer_trip 
			WHERE offer_id = $%d AND deleted = 0
		)`, argIndex))
		args = append(args, *query.TripOfferID)
		argIndex++
	}

	// Set defaults
	limit := query.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	orderBy := utils.SafeString(query.OrderBy)
	if orderBy == "" {
		orderBy = DefaultOrderBy
	}

	orderDir := utils.SafeString(query.OrderDir)
	if orderDir == "" {
		orderDir = DefaultOrderDir
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Updated query with ST_AsText to convert geometry to text
	queryStr := fmt.Sprintf(`
		SELECT id, driver_id, vehicle_id, from_address, to_address, 
		       from_country, to_country, start_date, end_date,
		       ST_AsText(from_location) as from_location_txt,
		       ST_AsText(to_location) as to_location_txt,
		       distance_km, status, meta, meta2, meta3, gps_logs,
		       created_at, updated_at, deleted
		FROM tbl_trip 
		%s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, orderDir, argIndex, argIndex+1)

	args = append(args, limit, query.Offset)

	var tripScans []TripScan
	err := pgxscan.Select(context.Background(), db.DB, &tripScans, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips: %w", err)
	}

	trips := make([]dto.Trip, len(tripScans))
	for i, scan := range tripScans {
		trips[i] = scan.ToTrip()
	}

	return trips, nil
}

func GetGPSLogs(query dto.GPSLogQuery) ([]dto.GPSLog, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if query.TripID != nil {
		conditions = append(conditions, fmt.Sprintf("trip_id = $%d", argIndex))
		args = append(args, *query.TripID)
		argIndex++
	}

	if query.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("company_id = $%d", argIndex))
		args = append(args, *query.CompanyID)
		argIndex++
	}

	if query.OfferID != nil {
		conditions = append(conditions, fmt.Sprintf("offer_id = $%d", argIndex))
		args = append(args, *query.OfferID)
		argIndex++
	}

	if query.DriverID != nil {
		conditions = append(conditions, fmt.Sprintf("driver_id = $%d", argIndex))
		args = append(args, *query.DriverID)
		argIndex++
	}

	if query.VehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("vehicle_id = $%d", argIndex))
		args = append(args, *query.VehicleID)
		argIndex++
	}

	if query.From != nil {
		conditions = append(conditions, fmt.Sprintf("log_dt >= $%d", argIndex))
		args = append(args, *query.From)
		argIndex++
	}

	if query.To != nil {
		conditions = append(conditions, fmt.Sprintf("log_dt <= $%d", argIndex))
		args = append(args, *query.To)
		argIndex++
	}

	if query.TripOfferID != nil {
		conditions = append(conditions, fmt.Sprintf(`trip_id IN (
			SELECT trip_id FROM tbl_offer_trip 
			WHERE offer_id = $%d AND deleted = 0
		)`, argIndex))
		args = append(args, *query.TripOfferID)
		argIndex++
	}

	limit := query.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	orderBy := utils.SafeString(query.OrderBy)
	if orderBy == "" {
		orderBy = DefaultOrderBy
	}

	orderDir := utils.SafeString(query.OrderDir)
	if orderDir == "" {
		orderDir = DefaultOrderDir
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	queryStr := fmt.Sprintf(`
		SELECT id, company_id, vehicle_id, driver_id, offer_id, trip_id,
		       battery_level, speed, heading, accuracy,
		       ST_AsText(coordinates) as coordinates_txt,
		       status, log_dt, created_at
		FROM tbl_gps_log 
		%s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, orderDir, argIndex, argIndex+1)

	args = append(args, limit, query.Offset)

	var logScans []GPSLogScan
	err := pgxscan.Select(context.Background(), db.DB, &logScans, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get GPS logs: %w", err)
	}

	// Convert scanned results to dto.GPSLog
	logs := make([]dto.GPSLog, len(logScans))
	for i, scan := range logScans {
		logs[i] = scan.ToGPSLog()
	}

	return logs, nil
}

func GetLastPositions(query dto.PositionQuery) ([]dto.GPSLog, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if len(query.TripIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("trip_id = ANY($%d)", argIndex))
		args = append(args, query.TripIDs)
		argIndex++
	}

	if len(query.CompanyIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("company_id = ANY($%d)", argIndex))
		args = append(args, query.CompanyIDs)
		argIndex++
	}

	if len(query.OfferIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("offer_id = ANY($%d)", argIndex))
		args = append(args, query.OfferIDs)
		argIndex++
	}

	if len(query.DriverIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("driver_id = ANY($%d)", argIndex))
		args = append(args, query.DriverIDs)
		argIndex++
	}

	if len(query.VehicleIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("vehicle_id = ANY($%d)", argIndex))
		args = append(args, query.VehicleIDs)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Updated query with ST_AsText
	queryStr := fmt.Sprintf(`
		SELECT DISTINCT ON (COALESCE(trip_id, 0), COALESCE(driver_id, 0), COALESCE(vehicle_id, 0)) 
		       id, company_id, vehicle_id, driver_id, offer_id, trip_id,
		       battery_level, speed, heading, accuracy,
		       ST_AsText(coordinates) as coordinates_txt,
		       status, log_dt, created_at
		FROM tbl_gps_log 
		%s 
		ORDER BY COALESCE(trip_id, 0), COALESCE(driver_id, 0), COALESCE(vehicle_id, 0), log_dt DESC`,
		whereClause)

	var logScans []GPSLogScan
	err := pgxscan.Select(context.Background(), db.DB, &logScans, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get last positions: %w", err)
	}

	// Convert scanned results to dto.GPSLog
	logs := make([]dto.GPSLog, len(logScans))
	for i, scan := range logScans {
		logs[i] = scan.ToGPSLog()
	}

	return logs, nil
}

func CreateGPSLogs(logs []dto.GPSLogInput) error {
	if len(logs) == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, log := range logs {
		_, err = tx.Exec(ctx,
			`INSERT INTO tbl_gps_log (
				company_id, vehicle_id, driver_id, offer_id, trip_id,
				battery_level, speed, heading, accuracy, coordinates,
				status, log_dt, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'active', $11, CURRENT_TIMESTAMP)`,
			log.CompanyID, log.VehicleID, log.DriverID, log.OfferID, log.TripID,
			log.BatteryLevel, log.Speed, log.Heading, log.Accuracy, log.Coordinates,
			log.LogDt)
		if err != nil {
			return fmt.Errorf("failed to create GPS log: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func getIntValue(input *int, fallback int) int {
	if input != nil {
		return *input
	}
	return fallback
}

func getStringValue(input *string, fallback *string) *string {
	if input != nil {
		return input
	}
	return fallback
}
