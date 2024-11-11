package countryLib

import "time"

type Country struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	ISO3           *string    `json:"iso3"`
	NumericCode    *string    `json:"numeric_code"`
	ISO2           *string    `json:"iso2"`
	Phonecode      *string    `json:"phonecode"`
	Capital        *string    `json:"capital"`
	Currency       *string    `json:"currency"`
	CurrencyName   *string    `json:"currency_name"`
	CurrencySymbol *string    `json:"currency_symbol"`
	TLD            *string    `json:"tld"`
	Native         *string    `json:"native"`
	Region         *string    `json:"region"`
	RegionID       *int64     `json:"region_id"`
	Subregion      *string    `json:"subregion"`
	SubregionID    *int64     `json:"subregion_id"`
	Nationality    *string    `json:"nationality"`
	Timezones      *string    `json:"timezones"`
	Translations   *string    `json:"translations"`
	Latitude       *float64   `json:"latitude"`
	Longitude      *float64   `json:"longitude"`
	Emoji          *string    `json:"emoji"`
	EmojiU         *string    `json:"emojiU"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	Flag           int        `json:"flag"`
	WikiDataID     *string    `json:"wikiDataId"`
}

type CountryDetail struct {
	Country
	Cities []City `json:"cities,omitempty"`
}

type City struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	StateID     int64      `json:"state_id"`
	StateCode   string     `json:"state_code"`
	CountryID   int64      `json:"country_id"`
	CountryCode string     `json:"country_code"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Flag        int        `json:"flag"`
	WikiDataID  string     `json:"wikiDataId"`
}

type CitySearchResult struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	StateID     int64      `json:"state_id"`
	StateCode   string     `json:"state_code"`
	CountryID   int64      `json:"country_id"`
	CountryCode string     `json:"country_code"`
	CountryName string     `json:"country_name"`
	Latitude    *float64   `json:"latitude"`
	Longitude   *float64   `json:"longitude"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Flag        int        `json:"flag"`
	WikiDataID  *string    `json:"wikiDataId"`
}
