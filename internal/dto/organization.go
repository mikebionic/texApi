package dto

import "time"

type Organization struct {
	ID   int    `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
	OrganizationBasic
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    int       `json:"active"`
	Deleted   int       `json:"deleted"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	OrganizationBasic
}

type UpdateOrganizationRequest struct {
	Name   *string `json:"name"`
	Active *int    `json:"active"`
	OrganizationBasic
}

type OrganizationBasic struct {
	DescriptionEN   *string `json:"description_en"`
	DescriptionRU   *string `json:"description_ru"`
	DescriptionTK   *string `json:"description_tk"`
	Email           *string `json:"email"`
	ImageUrl        *string `json:"image_url"`
	LogoUrl         *string `json:"logo_url"`
	IconUrl         *string `json:"icon_url"`
	BannerUrl       *string `json:"banner_url"`
	WebsiteUrl      *string `json:"website_url"`
	AboutText       *string `json:"about_text"`
	RefundText      *string `json:"refund_text"`
	DeliveryText    *string `json:"delivery_text"`
	ContactText     *string `json:"contact_text"`
	TermsConditions *string `json:"terms_conditions"`
	PrivacyPolicy   *string `json:"privacy_policy"`
	Address1        *string `json:"address1"`
	Address2        *string `json:"address2"`
	Address3        *string `json:"address3"`
	Address4        *string `json:"address4"`
	AddressTitle1   *string `json:"address_title1"`
	AddressTitle2   *string `json:"address_title2"`
	AddressTitle3   *string `json:"address_title3"`
	AddressTitle4   *string `json:"address_title4"`
	ContactPhone1   *string `json:"contact_phone1"`
	ContactPhone2   *string `json:"contact_phone2"`
	ContactPhone3   *string `json:"contact_phone3"`
	ContactPhone4   *string `json:"contact_phone4"`
	ContactTitle1   *string `json:"contact_title1"`
	ContactTitle2   *string `json:"contact_title2"`
	ContactTitle3   *string `json:"contact_title3"`
	ContactTitle4   *string `json:"contact_title4"`
	Meta            *string `json:"meta"`
	Meta2           *string `json:"meta2"`
	Meta3           *string `json:"meta3"`
}
