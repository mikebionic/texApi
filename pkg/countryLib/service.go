package countryLib

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/pkg/utils"
)

func GetCountryList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		GetCountriesPaginatedQuery,
		perPage,
		offset,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var countries []Country
	var totalCount int

	for rows.Next() {
		var country Country
		err := rows.Scan(
			&country.ID,
			&country.Name,
			&country.ISO3,
			&country.NumericCode,
			&country.ISO2,
			&country.Phonecode,
			&country.Capital,
			&country.Currency,
			&country.CurrencyName,
			&country.CurrencySymbol,
			&country.TLD,
			&country.Native,
			&country.Region,
			&country.RegionID,
			&country.Subregion,
			&country.SubregionID,
			&country.Nationality,
			&country.Timezones,
			&country.Translations,
			&country.Latitude,
			&country.Longitude,
			&country.Emoji,
			&country.EmojiU,
			&country.CreatedAt,
			&country.UpdatedAt,
			&country.Flag,
			&country.WikiDataID,
			&totalCount,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}
		countries = append(countries, country)
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    countries,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Country list", response))
}

func GetCountry(ctx *gin.Context) {
	id := ctx.Param("id") // Get the country ID from the request URL

	var countryDetail CountryDetail
	var citiesJSON []byte

	// Execute the query
	err := db.DB.QueryRow(
		context.Background(),
		GetCountryWithCitiesQuery, // This is the updated query
		id,
	).Scan(
		&countryDetail.ID,
		&countryDetail.Name,
		&countryDetail.ISO3,
		&countryDetail.NumericCode,
		&countryDetail.ISO2,
		&countryDetail.Phonecode,
		&countryDetail.Capital,
		&countryDetail.Currency,
		&countryDetail.CurrencyName,
		&countryDetail.CurrencySymbol,
		&countryDetail.TLD,
		&countryDetail.Native,
		&countryDetail.Region,
		&countryDetail.RegionID,
		&countryDetail.Subregion,
		&countryDetail.SubregionID,
		&countryDetail.Nationality,
		&countryDetail.Timezones,
		&countryDetail.Translations,
		&countryDetail.Latitude,
		&countryDetail.Longitude,
		&countryDetail.Emoji,
		&countryDetail.EmojiU,
		&countryDetail.CreatedAt,
		&countryDetail.UpdatedAt,
		&countryDetail.Flag,
		&countryDetail.WikiDataID,
		&citiesJSON, // Store the cities as JSON
	)

	// Handle errors (e.g., country not found)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Country not found", err.Error()))
		return
	}

	// Unmarshal the cities from JSON and assign to countryDetail
	json.Unmarshal(citiesJSON, &countryDetail.Cities)

	// Return the country details along with its cities
	ctx.JSON(http.StatusOK, utils.FormatResponse("Country details", countryDetail))
}

func GetCityList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		GetCitiesPaginatedQuery,
		perPage,
		offset,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var cities []City
	var totalCount int

	for rows.Next() {
		var city City
		err := rows.Scan(
			&city.ID, &city.Name, &city.StateID, &city.StateCode,
			&city.CountryID, &city.CountryCode, &city.Latitude,
			&city.Longitude, &city.CreatedAt, &city.UpdatedAt,
			&city.Flag, &city.WikiDataID, &totalCount,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}
		cities = append(cities, city)
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    cities,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("City list", response))
}

func SearchCountries(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	search := ctx.Query("q")
	if search == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Search query is required", "Please provide a search term using the 'q' parameter"))
		return
	}

	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		SearchCountriesQuery,
		perPage,
		offset,
		search,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var countries []Country
	var totalCount int

	for rows.Next() {
		var country Country
		err := rows.Scan(
			&country.ID,
			&country.Name,
			&country.ISO3,
			&country.NumericCode,
			&country.ISO2,
			&country.Phonecode,
			&country.Capital,
			&country.Currency,
			&country.CurrencyName,
			&country.CurrencySymbol,
			&country.TLD,
			&country.Native,
			&country.Region,
			&country.RegionID,
			&country.Subregion,
			&country.SubregionID,
			&country.Nationality,
			&country.Timezones,
			&country.Translations,
			&country.Latitude,
			&country.Longitude,
			&country.Emoji,
			&country.EmojiU,
			&country.CreatedAt,
			&country.UpdatedAt,
			&country.Flag,
			&country.WikiDataID,
			&totalCount,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}
		countries = append(countries, country)
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    countries,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Country search results", response))
}

func SearchCities(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	search := ctx.Query("q")
	if search == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Search query is required", "Please provide a search term using the 'q' parameter"))
		return
	}

	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		SearchCitiesQuery,
		perPage,
		offset,
		search,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var cities []CitySearchResult
	var totalCount int

	for rows.Next() {
		var city CitySearchResult
		err := rows.Scan(
			&city.ID,
			&city.Name,
			&city.StateID,
			&city.StateCode,
			&city.CountryID,
			&city.CountryCode,
			&city.CountryName,
			&city.Latitude,
			&city.Longitude,
			&city.CreatedAt,
			&city.UpdatedAt,
			&city.Flag,
			&city.WikiDataID,
			&totalCount,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}
		cities = append(cities, city)
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    cities,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("City search results", response))
}
