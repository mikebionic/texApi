package dto

import (
	"github.com/google/uuid"
	"time"
)

type Plan struct {
	ID                 int        `json:"id"`
	UUID               uuid.UUID  `json:"uuid"`
	Name               string     `json:"name"`
	Code               string     `json:"code"`
	Provider           string     `json:"provider"`
	Region             string     `json:"region"`
	PriceUSD           float64    `json:"price_usd"`
	PriceLocal         *float64   `json:"price_local"`
	LocalCurrency      *string    `json:"local_currency"`
	BillingCycle       string     `json:"billing_cycle"`
	LoadPostsLimit     *int       `json:"load_posts_limit"`
	LoadPostsUnlimited bool       `json:"load_posts_unlimited"`
	GPSTrackingLevel   string     `json:"gps_tracking_level"`
	GPSHasETA          bool       `json:"gps_has_eta"`
	RateToolsLevel     string     `json:"rate_tools_level"`
	RateToolsFeatures  *[]string  `json:"rate_tools_features"`
	EdocsAvailable     bool       `json:"edocs_available"`
	EdocsLimit         *int       `json:"edocs_limit"`
	EdocsHasArchiving  bool       `json:"edocs_has_archiving"`
	SupportLevel       string     `json:"support_level"`
	PaymentGuarantee   bool       `json:"payment_guarantee"`
	APIAccess          bool       `json:"api_access"`
	DisplayOrder       int        `json:"display_order"`
	IsPopular          bool       `json:"is_popular"`
	IsRecommended      bool       `json:"is_recommended"`
	Description        *string    `json:"description"`
	FeaturesSummary    *string    `json:"features_summary"`
	AvailableFrom      *time.Time `json:"available_from"`
	AvailableUntil     *time.Time `json:"available_until"`
	Meta               *string    `json:"meta"`
	Meta2              *string    `json:"meta2"`
	Meta3              *string    `json:"meta3"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type CreatePlanRequest struct {
	Name               string     `json:"name" binding:"required,max=100"`
	Code               string     `json:"code" binding:"required,max=50"`
	Provider           string     `json:"provider" binding:"required,max=50"`
	Region             string     `json:"region" binding:"required,max=50"`
	PriceUSD           float64    `json:"price_usd" binding:"required,min=0"`
	PriceLocal         *float64   `json:"price_local" binding:"omitempty,min=0"`
	LocalCurrency      *string    `json:"local_currency" binding:"omitempty,max=10"`
	BillingCycle       string     `json:"billing_cycle" binding:"required,oneof=monthly yearly quarterly"`
	LoadPostsLimit     *int       `json:"load_posts_limit" binding:"omitempty,min=0"`
	LoadPostsUnlimited bool       `json:"load_posts_unlimited"`
	GPSTrackingLevel   string     `json:"gps_tracking_level" binding:"required,oneof=none basic advanced full"`
	GPSHasETA          bool       `json:"gps_has_eta"`
	RateToolsLevel     string     `json:"rate_tools_level" binding:"required,oneof=none basic advanced full"`
	RateToolsFeatures  *[]string  `json:"rate_tools_features"`
	EdocsAvailable     bool       `json:"edocs_available"`
	EdocsLimit         *int       `json:"edocs_limit" binding:"omitempty,min=0"`
	EdocsHasArchiving  bool       `json:"edocs_has_archiving"`
	SupportLevel       string     `json:"support_level" binding:"required,oneof=none email phone priority dedicated"`
	PaymentGuarantee   bool       `json:"payment_guarantee"`
	APIAccess          bool       `json:"api_access"`
	DisplayOrder       int        `json:"display_order"`
	IsPopular          bool       `json:"is_popular"`
	IsRecommended      bool       `json:"is_recommended"`
	Description        *string    `json:"description" binding:"omitempty,max=1000"`
	FeaturesSummary    *string    `json:"features_summary"`
	AvailableFrom      *time.Time `json:"available_from"`
	AvailableUntil     *time.Time `json:"available_until"`
	Meta               *string    `json:"meta"`
	Meta2              *string    `json:"meta2"`
	Meta3              *string    `json:"meta3"`
}

type UpdatePlanRequest struct {
	Name               *string    `json:"name" binding:"omitempty,max=100"`
	Code               *string    `json:"code" binding:"omitempty,max=50"`
	Provider           *string    `json:"provider" binding:"omitempty,max=50"`
	Region             *string    `json:"region" binding:"omitempty,max=50"`
	PriceUSD           *float64   `json:"price_usd" binding:"omitempty,min=0"`
	PriceLocal         *float64   `json:"price_local" binding:"omitempty,min=0"`
	LocalCurrency      *string    `json:"local_currency" binding:"omitempty,max=10"`
	BillingCycle       *string    `json:"billing_cycle" binding:"omitempty,oneof=monthly yearly quarterly"`
	LoadPostsLimit     *int       `json:"load_posts_limit" binding:"omitempty,min=0"`
	LoadPostsUnlimited *bool      `json:"load_posts_unlimited"`
	GPSTrackingLevel   *string    `json:"gps_tracking_level" binding:"omitempty,oneof=none basic advanced full"`
	GPSHasETA          *bool      `json:"gps_has_eta"`
	RateToolsLevel     *string    `json:"rate_tools_level" binding:"omitempty,oneof=none basic advanced full"`
	RateToolsFeatures  *[]string  `json:"rate_tools_features"`
	EdocsAvailable     *bool      `json:"edocs_available"`
	EdocsLimit         *int       `json:"edocs_limit" binding:"omitempty,min=0"`
	EdocsHasArchiving  *bool      `json:"edocs_has_archiving"`
	SupportLevel       *string    `json:"support_level" binding:"omitempty,oneof=none email phone priority dedicated"`
	PaymentGuarantee   *bool      `json:"payment_guarantee"`
	APIAccess          *bool      `json:"api_access"`
	DisplayOrder       *int       `json:"display_order"`
	IsPopular          *bool      `json:"is_popular"`
	IsRecommended      *bool      `json:"is_recommended"`
	Description        *string    `json:"description" binding:"omitempty,max=1000"`
	FeaturesSummary    *string    `json:"features_summary"`
	AvailableFrom      *time.Time `json:"available_from"`
	AvailableUntil     *time.Time `json:"available_until"`
	Meta               *string    `json:"meta"`
	Meta2              *string    `json:"meta2"`
	Meta3              *string    `json:"meta3"`
}

type PlanFilter struct {
	Provider      *string `form:"provider"`
	Region        *string `form:"region"`
	BillingCycle  *string `form:"billing_cycle" binding:"omitempty,oneof=monthly yearly quarterly"`
	IsPopular     *bool   `form:"is_popular"`
	IsRecommended *bool   `form:"is_recommended"`
	Active        *bool   `form:"active"`
	Search        *string `form:"search"`
	Page          int     `form:"page,default=1" binding:"min=1"`
	PerPage       int     `form:"per_page,default=10" binding:"min=1,max=100"`
}
