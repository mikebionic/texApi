package dto

import (
	"github.com/google/uuid"
	"time"
)

type Analytics struct {
	ID                   int       `json:"id"`
	UUID                 uuid.UUID `json:"uuid"`
	UserAll              int       `json:"user_all"`
	UserSender           int       `json:"user_sender"`
	UserCarrier          int       `json:"user_carrier"`
	LastUserID           int       `json:"last_user_id"`
	UserSenderNew        int       `json:"user_sender_new"`
	UserCarrierNew       int       `json:"user_carrier_new"`
	LastOfferID          int       `json:"last_offer_id"`
	OfferNewSender       int       `json:"offer_new_sender"`
	OfferNewCarrier      int       `json:"offer_new_carrier"`
	OfferAll             int       `json:"offer_all"`
	OfferActive          int       `json:"offer_active"`
	OfferPending         int       `json:"offer_pending"`
	OfferCompleted       int       `json:"offer_completed"`
	OfferNoResponse      int       `json:"offer_no_response"`
	LastCompletedOfferID int       `json:"last_completed_offer_id"`
	TotalRevenue         float64   `json:"total_revenue"`
	AverageCostPerKm     float64   `json:"average_cost_per_km"`
	TotalDistance        int       `json:"total_distance"`
	ActiveCompanies      int       `json:"active_companies"`
	PeriodStart          time.Time `json:"period_start"`
	PeriodEnd            time.Time `json:"period_end"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type AnalyticsFilter struct {
	DateFrom *time.Time `form:"date_from" binding:"omitempty"`
	DateTo   *time.Time `form:"date_to" binding:"omitempty"`

	UserAllMin  *int     `form:"user_all_min" binding:"omitempty,min=0"`
	UserAllMax  *int     `form:"user_all_max" binding:"omitempty,min=0"`
	OfferAllMin *int     `form:"offer_all_min" binding:"omitempty,min=0"`
	OfferAllMax *int     `form:"offer_all_max" binding:"omitempty,min=0"`
	RevenueMin  *float64 `form:"revenue_min" binding:"omitempty,min=0"`
	RevenueMax  *float64 `form:"revenue_max" binding:"omitempty,min=0"`

	UserSenderMin  *int `form:"user_sender_min" binding:"omitempty,min=0"`
	UserSenderMax  *int `form:"user_sender_max" binding:"omitempty,min=0"`
	UserCarrierMin *int `form:"user_carrier_min" binding:"omitempty,min=0"`
	UserCarrierMax *int `form:"user_carrier_max" binding:"omitempty,min=0"`

	OfferActiveMin    *int `form:"offer_active_min" binding:"omitempty,min=0"`
	OfferActiveMax    *int `form:"offer_active_max" binding:"omitempty,min=0"`
	OfferPendingMin   *int `form:"offer_pending_min" binding:"omitempty,min=0"`
	OfferPendingMax   *int `form:"offer_pending_max" binding:"omitempty,min=0"`
	OfferCompletedMin *int `form:"offer_completed_min" binding:"omitempty,min=0"`
	OfferCompletedMax *int `form:"offer_completed_max" binding:"omitempty,min=0"`

	PeriodStart *time.Time `form:"period_start" binding:"omitempty"`
	PeriodEnd   *time.Time `form:"period_end" binding:"omitempty"`

	OrderBy  string `form:"order_by" binding:"omitempty,oneof=created_at updated_at user_all user_sender user_carrier offer_all offer_active offer_pending offer_completed total_revenue period_start period_end"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`

	Page    int `form:"page,default=1" binding:"min=1"`
	PerPage int `form:"per_page,default=10" binding:"min=1,max=100"`
}

// summary
type AnalyticsStats struct {
	TotalRecords       int       `json:"total_records"`
	AvgUsersPerPeriod  float64   `json:"avg_users_per_period"`
	AvgOffersPerPeriod float64   `json:"avg_offers_per_period"`
	TotalRevenue       float64   `json:"total_revenue"`
	GrowthRate         float64   `json:"growth_rate"`
	LastUpdate         time.Time `json:"last_update"`
}

// full response with stats
type AnalyticsResponse struct {
	Stats AnalyticsStats `json:"stats"`
	Data  []Analytics    `json:"data"`
}
