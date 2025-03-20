package dto

import (
	"database/sql"
	"encoding/json"
	"time"
)

// GPS Location DTOs
type GPSLocationDetails struct {
	ID           int             `json:"id" db:"id"`
	UUID         string          `json:"uuid" db:"uuid"`
	VehicleID    int             `json:"vehicle_id" db:"vehicle_id"`
	DriverID     int             `json:"driver_id" db:"driver_id"`
	Latitude     float64         `json:"latitude" db:"latitude"`
	Longitude    float64         `json:"longitude" db:"longitude"`
	Altitude     sql.NullFloat64 `json:"altitude,omitempty" db:"altitude"`
	Speed        sql.NullFloat64 `json:"speed,omitempty" db:"speed"`
	Direction    sql.NullFloat64 `json:"direction,omitempty" db:"direction"`
	Accuracy     sql.NullFloat64 `json:"accuracy,omitempty" db:"accuracy"`
	LocationTime time.Time       `json:"location_time" db:"location_time"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	Meta         json.RawMessage `json:"meta,omitempty" db:"meta"`

	// Additional fields for responses
	DriverName   string `json:"driver_name,omitempty" db:"driver_name"`
	VehiclePlate string `json:"vehicle_plate,omitempty" db:"vehicle_plate"`
	LocationAge  string `json:"location_age,omitempty"`
	TotalCount   int    `json:"total_count,omitempty" db:"total_count"`
}

type GPSLocationCreate struct {
	VehicleID    int             `json:"vehicle_id" binding:"required"`
	DriverID     int             `json:"driver_id" binding:"required"`
	Latitude     float64         `json:"latitude" binding:"required"`
	Longitude    float64         `json:"longitude" binding:"required"`
	Altitude     *float64        `json:"altitude,omitempty"`
	Speed        *float64        `json:"speed,omitempty"`
	Direction    *float64        `json:"direction,omitempty"`
	Accuracy     *float64        `json:"accuracy,omitempty"`
	LocationTime time.Time       `json:"location_time" binding:"required"`
	DeviceID     string          `json:"device_id" binding:"required"`
	Meta         json.RawMessage `json:"meta,omitempty"`
}

type GPSBatchLocationCreate struct {
	DeviceID  string              `json:"device_id" binding:"required"`
	Locations []GPSLocationCreate `json:"locations" binding:"required,min=1"`
}

// GPS Trip DTOs
type GPSTripDetails struct {
	ID            int             `json:"id" db:"id"`
	UUID          string          `json:"uuid" db:"uuid"`
	VehicleID     int             `json:"vehicle_id" db:"vehicle_id"`
	DriverID      int             `json:"driver_id" db:"driver_id"`
	StartTime     time.Time       `json:"start_time" db:"start_time"`
	EndTime       sql.NullTime    `json:"end_time,omitempty" db:"end_time"`
	Distance      sql.NullFloat64 `json:"distance,omitempty" db:"distance"`
	AvgSpeed      sql.NullFloat64 `json:"avg_speed,omitempty" db:"avg_speed"`
	MaxSpeed      sql.NullFloat64 `json:"max_speed,omitempty" db:"max_speed"`
	Status        string          `json:"status" db:"status"`
	StartLocation json.RawMessage `json:"start_location,omitempty" db:"start_location"`
	EndLocation   json.RawMessage `json:"end_location,omitempty" db:"end_location"`
	Meta          json.RawMessage `json:"meta,omitempty" db:"meta"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`

	// Additional fields for responses
	DriverName   string `json:"driver_name,omitempty" db:"driver_name"`
	VehiclePlate string `json:"vehicle_plate,omitempty" db:"vehicle_plate"`
	Duration     string `json:"duration,omitempty"`
	TotalCount   int    `json:"total_count,omitempty" db:"total_count"`
}

type GPSTripCreate struct {
	VehicleID     int             `json:"vehicle_id" binding:"required"`
	DriverID      int             `json:"driver_id" binding:"required"`
	StartTime     time.Time       `json:"start_time" binding:"required"`
	EndTime       *time.Time      `json:"end_time,omitempty"`
	Distance      *float64        `json:"distance,omitempty"`
	AvgSpeed      *float64        `json:"avg_speed,omitempty"`
	MaxSpeed      *float64        `json:"max_speed,omitempty"`
	Status        string          `json:"status" binding:"required"`
	StartLocation json.RawMessage `json:"start_location,omitempty"`
	EndLocation   json.RawMessage `json:"end_location,omitempty"`
	Meta          json.RawMessage `json:"meta,omitempty"`
}

type GPSTripUpdate struct {
	EndTime     *time.Time      `json:"end_time,omitempty"`
	Distance    *float64        `json:"distance,omitempty"`
	AvgSpeed    *float64        `json:"avg_speed,omitempty"`
	MaxSpeed    *float64        `json:"max_speed,omitempty"`
	Status      *string         `json:"status,omitempty"`
	EndLocation json.RawMessage `json:"end_location,omitempty"`
	Meta        json.RawMessage `json:"meta,omitempty"`
}

// Analytics DTOs
type VehicleAnalytics struct {
	VehicleID        int            `json:"vehicle_id" db:"vehicle_id"`
	TotalDistance    float64        `json:"total_distance" db:"total_distance"`
	TotalDriveTime   string         `json:"total_drive_time"`
	AvgDailyDistance float64        `json:"avg_daily_distance" db:"avg_daily_distance"`
	MaxDailyDistance float64        `json:"max_daily_distance" db:"max_daily_distance"`
	AvgSpeed         float64        `json:"avg_speed" db:"avg_speed"`
	MaxSpeed         float64        `json:"max_speed" db:"max_speed"`
	IdleTime         string         `json:"idle_time"`
	StopCount        int            `json:"stop_count" db:"stop_count"`
	GeofenceVisits   map[string]int `json:"geofence_visits,omitempty"`
}

// Map Data DTOs
type MapDataPoint struct {
	VehicleID    int       `json:"vehicle_id"`
	DriverID     int       `json:"driver_id"`
	DriverName   string    `json:"driver_name"`
	VehiclePlate string    `json:"vehicle_plate"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Speed        float64   `json:"speed"`
	Direction    float64   `json:"direction"`
	Timestamp    time.Time `json:"timestamp"`
}

type GeoJSONFeature struct {
	Type       string          `json:"type"`
	Properties json.RawMessage `json:"properties"`
	Geometry   struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

type GeoJSONResponse struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}
