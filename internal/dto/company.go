package dto

import "time"

//	type CompanyDetails struct {
//		CompanyCreate
//		Drivers    []DriverShort  `json:"drivers"`
//		Vehicles   []VehicleShort `json:"vehicles"`
//		TotalCount int            `json:"total_count"`
//	}
type CompanyDetails struct {
	CompanyMain
	Drivers        []DriverShort  `json:"drivers"`
	Vehicles       []VehicleShort `json:"vehicles"`
	TotalCount     int            `json:"total_count"`
	FollowersCount int            `json:"followers_count"`
	FollowingCount int            `json:"following_count"`
	LastActive     time.Time      `json:"last_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type CompanyMain struct {
	ID                  int    `json:"id"`
	UUID                string `json:"uuid"`
	UserID              int    `json:"user_id"`
	Role                string `json:"role"`
	RoleID              int    `json:"role_id"`
	Plan                string `json:"plan"`
	PlanActive          int    `json:"plan_active"`
	CompanyName         string `json:"company_name"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	PatronymicName      string `json:"patronymic_name"`
	About               string `json:"about"`
	Phone               string `json:"phone"`
	Phone2              string `json:"phone2"`
	Phone3              string `json:"phone3"`
	Email               string `json:"email"`
	Email2              string `json:"email2"`
	Email3              string `json:"email3"`
	Meta                string `json:"meta"`
	Meta2               string `json:"meta2"`
	Meta3               string `json:"meta3"`
	Address             string `json:"address"`
	Country             string `json:"country"`
	CountryID           int    `json:"country_id"`
	CityID              int    `json:"city_id"`
	ImageURL            string `json:"image_url"`
	Verified            int    `json:"verified"`
	ConfirmationRequest int    `json:"confirmation_request"`
	Entity              string `json:"entity"`
	Featured            int    `json:"featured"`
	Rating              int    `json:"rating"`
	Partner             int    `json:"partner"`
	ViewCount           int    `json:"view_count"`
	SuccessfulOps       int    `json:"successful_ops"`

	SelfDestructDuration int      `json:"self_destruct_duration"`
	Passkey              string   `json:"passkey"`
	Blacklist            []string `json:"blacklist"`
	LoginDevices         []string `json:"login_devices"`

	// Visibility settings
	ShowAvatar      string `json:"show_avatar"`
	ShowBio         string `json:"show_bio"`
	ShowLastSeen    string `json:"show_last_seen"`
	ShowPhoneNumber string `json:"show_phone_number"`
	ReceiveCalls    string `json:"receive_calls"`
	InviteGroup     string `json:"invite_group"`

	// Notification settings
	NotificationsChat      int `json:"notifications_chat"`
	NotificationsGroup     int `json:"notifications_group"`
	NotificationsStory     int `json:"notifications_story"`
	NotificationsReactions int `json:"notifications_reactions"`

	// Exceptions lists
	AvatarExceptions       []string `json:"avatar_exceptions"`
	BioExceptions          []string `json:"bio_exceptions"`
	LastSeenExceptions     []string `json:"last_seen_exceptions"`
	PhoneNumberExceptions  []string `json:"phone_number_exceptions"`
	ReceiveCallsExceptions []string `json:"receive_calls_exceptions"`
	InviteGroupExceptions  []string `json:"invite_group_exceptions"`

	Blocked int `json:"blocked"`
	Active  int `json:"active"`
	Deleted int `json:"deleted"`
}

type CompanyMainStringTime struct {
	CompanyMain
	LastActive string `json:"last_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CompanyCreateShort struct {
	UserID      int    `json:"user_id"`
	Role        string `json:"role"`
	RoleID      int    `json:"role_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CompanyName string `json:"company_name"`
	About       string `json:"about"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ImageURL    string `json:"image_url"`
}

type CompanyCreate struct {
	ID                  int    `json:"id"`
	UUID                string `json:"uuid"`
	UserID              int    `json:"user_id"`
	Role                string `json:"role"`
	RoleID              int    `json:"role_id"`
	Plan                string `json:"plan"`
	PlanActive          int    `json:"plan_active"`
	CompanyName         string `json:"company_name"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	PatronymicName      string `json:"patronymic_name"`
	About               string `json:"about"`
	Phone               string `json:"phone"`
	Phone2              string `json:"phone2"`
	Phone3              string `json:"phone3"`
	Email               string `json:"email"`
	Email2              string `json:"email2"`
	Email3              string `json:"email3"`
	Meta                string `json:"meta"`
	Meta2               string `json:"meta2"`
	Meta3               string `json:"meta3"`
	Address             string `json:"address"`
	Country             string `json:"country"`
	CountryID           int    `json:"country_id"`
	CityID              int    `json:"city_id"`
	ImageURL            string `json:"image_url"`
	Verified            int    `json:"verified"`
	ConfirmationRequest int    `json:"confirmation_request"`
	Entity              string `json:"entity"`
	Featured            int    `json:"featured"`
	Rating              int    `json:"rating"`
	Partner             int    `json:"partner"`
	ViewCount           int    `json:"view_count"`
	SuccessfulOps       int    `json:"successful_ops"`

	// Privacy settings
	ShowAvatar      string `json:"show_avatar"`
	ShowBio         string `json:"show_bio"`
	ShowLastSeen    string `json:"show_last_seen"`
	ShowPhoneNumber string `json:"show_phone_number"`
	ReceiveCalls    string `json:"receive_calls"`
	InviteGroup     string `json:"invite_group"`

	// Security settings
	SelfDestructDuration int      `json:"self_destruct_duration"`
	Passkey              string   `json:"passkey"`
	Blacklist            []string `json:"blacklist"`
	LoginDevices         []string `json:"login_devices"`

	// Notification settings
	NotificationsChat      int `json:"notifications_chat"`
	NotificationsGroup     int `json:"notifications_group"`
	NotificationsStory     int `json:"notifications_story"`
	NotificationsReactions int `json:"notifications_reactions"`

	// Exceptions lists
	AvatarExceptions       []string `json:"avatar_exceptions"`
	BioExceptions          []string `json:"bio_exceptions"`
	LastSeenExceptions     []string `json:"last_seen_exceptions"`
	PhoneNumberExceptions  []string `json:"phone_number_exceptions"`
	ReceiveCallsExceptions []string `json:"receive_calls_exceptions"`
	InviteGroupExceptions  []string `json:"invite_group_exceptions"`
}

type CompanyUpdate struct {
	CompanyID      *int    `json:"company_id,omitempty"`
	UserID         *int    `json:"user_id,omitempty"`
	Role           *string `json:"role,omitempty"`
	RoleID         *int    `json:"role_id,omitempty"`
	CompanyName    *string `json:"company_name,omitempty"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	PatronymicName *string `json:"patronymic_name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Phone2         *string `json:"phone2,omitempty"`
	Phone3         *string `json:"phone3,omitempty"`
	Email          *string `json:"email,omitempty"`
	Email2         *string `json:"email2,omitempty"`
	Email3         *string `json:"email3,omitempty"`
	Meta           *string `json:"meta,omitempty"`
	Meta2          *string `json:"meta2,omitempty"`
	Meta3          *string `json:"meta3,omitempty"`
	Address        *string `json:"address,omitempty"`
	Country        *string `json:"country,omitempty"`
	CountryID      *int    `json:"country_id,omitempty"`
	CityID         *int    `json:"city_id,omitempty"`
	ImageURL       *string `json:"image_url,omitempty"`
	Entity         *string `json:"entity,omitempty"`
	Featured       *int    `json:"featured,omitempty"`
	Rating         *int    `json:"rating,omitempty"`
	Partner        *int    `json:"partner,omitempty"`
	SuccessfulOps  *int    `json:"successful_ops,omitempty"`

	// Privacy settings
	ShowAvatar      *string `json:"show_avatar,omitempty"`
	ShowBio         *string `json:"show_bio,omitempty"`
	ShowLastSeen    *string `json:"show_last_seen,omitempty"`
	ShowPhoneNumber *string `json:"show_phone_number,omitempty"`
	ReceiveCalls    *string `json:"receive_calls,omitempty"`
	InviteGroup     *string `json:"invite_group,omitempty"`

	// Security settings
	SelfDestructDuration *int      `json:"self_destruct_duration,omitempty"`
	Passkey              *string   `json:"passkey,omitempty"`
	Blacklist            *[]string `json:"blacklist,omitempty"`
	LoginDevices         *[]string `json:"login_devices,omitempty"`

	// Notification settings
	NotificationsChat      *int `json:"notifications_chat,omitempty"`
	NotificationsGroup     *int `json:"notifications_group,omitempty"`
	NotificationsStory     *int `json:"notifications_story,omitempty"`
	NotificationsReactions *int `json:"notifications_reactions,omitempty"`

	// Exceptions lists
	AvatarExceptions       *[]string `json:"avatar_exceptions,omitempty"`
	BioExceptions          *[]string `json:"bio_exceptions,omitempty"`
	LastSeenExceptions     *[]string `json:"last_seen_exceptions,omitempty"`
	PhoneNumberExceptions  *[]string `json:"phone_number_exceptions,omitempty"`
	ReceiveCallsExceptions *[]string `json:"receive_calls_exceptions,omitempty"`
	InviteGroupExceptions  *[]string `json:"invite_group_exceptions,omitempty"`

	Blocked *int `json:"blocked,omitempty"`
	Active  *int `json:"active,omitempty"`
	Deleted *int `json:"deleted,omitempty"`
}

// Basic DTOs for nested responses
type CompanyBasic struct {
	ID          int    `json:"id"`
	CompanyName string `json:"company_name"`
	Country     string `json:"country"`
}
