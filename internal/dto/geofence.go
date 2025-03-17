package dto

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Geofence DTOs
type GeofenceDetails struct {
	ID          int             `json:"id" db:"id"`
	UUID        string          `json:"uuid" db:"uuid"`
	CompanyID   int             `json:"company_id" db:"company_id"`
	Name        string          `json:"name" db:"name"`
	Description sql.NullString  `json:"description,omitempty" db:"description"`
	FenceType   string          `json:"fence_type" db:"fence_type"`
	Coordinates json.RawMessage `json:"coordinates" db:"coordinates"`
	Radius      sql.NullFloat64 `json:"radius,omitempty" db:"radius"`
	IsActive    bool            `json:"is_active" db:"is_active"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	TotalCount  int             `json:"total_count,omitempty" db:"total_count"`
}

type GeofenceCreate struct {
	CompanyID   int             `json:"company_id" binding:"required"`
	Name        string          `json:"name" binding:"required"`
	Description *string         `json:"description,omitempty"`
	FenceType   string          `json:"fence_type" binding:"required"`
	Coordinates json.RawMessage `json:"coordinates" binding:"required"`
	Radius      *float64        `json:"radius,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
}

type GeofenceUpdate struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	FenceType   *string         `json:"fence_type,omitempty"`
	Coordinates json.RawMessage `json:"coordinates,omitempty"`
	Radius      *float64        `json:"radius,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
}

// Geofence Event DTOs
type GeofenceEventDetails struct {
	ID         int             `json:"id" db:"id"`
	GeofenceID int             `json:"geofence_id" db:"geofence_id"`
	VehicleID  int             `json:"vehicle_id" db:"vehicle_id"`
	DriverID   int             `json:"driver_id" db:"driver_id"`
	EventType  string          `json:"event_type" db:"event_type"`
	EventTime  time.Time       `json:"event_time" db:"event_time"`
	Location   json.RawMessage `json:"location" db:"location"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`

	// Additional fields for responses
	GeofenceName string `json:"geofence_name,omitempty" db:"geofence_name"`
	DriverName   string `json:"driver_name,omitempty" db:"driver_name"`
	VehiclePlate string `json:"vehicle_plate,omitempty" db:"vehicle_plate"`
	TotalCount   int    `json:"total_count,omitempty" db:"total_count"`
}
