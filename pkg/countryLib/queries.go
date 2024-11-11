package countryLib

const GetCountriesPaginatedQuery = `
WITH country_data AS (
    SELECT 
        c.*,
        COUNT(*) OVER() as total_count
    FROM tbl_country c
    WHERE c.flag = 1
    ORDER BY c.name
    LIMIT $1 OFFSET $2
)
SELECT 
    id, 
    name,
    iso3,
    numeric_code,
    iso2,
    phonecode,
    capital,
    currency,
    currency_name,
    currency_symbol,
    tld,
    native,
    region,
    region_id,
    subregion,
    subregion_id,
    nationality,
    timezones,
    translations,
    latitude,
    longitude,
    emoji,
    "emojiU",
    created_at,
    updated_at,
    flag,
    "wikiDataId",
    total_count
FROM country_data;
`

const GetCountryWithCitiesQuery = `
SELECT 
    c.id,
    c.name,
    c.iso3,
    c.numeric_code,
    c.iso2,
    c.phonecode,
    c.capital,
    c.currency,
    c.currency_name,
    c.currency_symbol,
    c.tld,
    c.native,
    c.region,
    c.region_id,
    c.subregion,
    c.subregion_id,
    c.nationality,
    c.timezones,
    c.translations,
    c.latitude,
    c.longitude,
    c.emoji,
    c."emojiU",
    c.created_at,
    c.updated_at,
    c.flag,
    c."wikiDataId",
    json_agg(
        json_build_object(
            'id', ci.id,
            'name', ci.name,
            'state_id', ci.state_id,
            'state_code', ci.state_code,
            'country_id', ci.country_id,
            'country_code', ci.country_code,
            'latitude', ci.latitude,
            'longitude', ci.longitude,
            'created_at', ci.created_at,
            'updated_at', ci.updated_at,
            'flag', ci.flag,
            'wikiDataId', ci."wikiDataId"
        )
    ) AS cities
FROM tbl_country c
LEFT JOIN tbl_city ci ON ci.country_id = c.id
WHERE c.id = $1
GROUP BY c.id;
`

const GetCitiesPaginatedQuery = `
WITH city_data AS (
    SELECT 
        c.*,
        COUNT(*) OVER() as total_count
    FROM tbl_city c
    WHERE c.flag = 1
    ORDER BY c.name
    LIMIT $1 OFFSET $2
)
SELECT * FROM city_data;
`

const SearchCountriesQuery = `
WITH country_data AS (
    SELECT 
        c.*,
        COUNT(*) OVER() as total_count
    FROM tbl_country c
    WHERE c.flag = 1
    AND LOWER(c.name) LIKE LOWER($3 || '%')
    ORDER BY c.name
    LIMIT $1 OFFSET $2
)
SELECT 
    id, 
    name,
    iso3,
    numeric_code,
    iso2,
    phonecode,
    capital,
    currency,
    currency_name,
    currency_symbol,
    tld,
    native,
    region,
    region_id,
    subregion,
    subregion_id,
    nationality,
    timezones,
    translations,
    latitude,
    longitude,
    emoji,
    "emojiU",
    created_at,
    updated_at,
    flag,
    "wikiDataId",
    total_count
FROM country_data;
`

const SearchCitiesQuery = `
WITH city_data AS (
    SELECT 
        ci.*,
        co.name as country_name,
        COUNT(*) OVER() as total_count
    FROM tbl_city ci
    JOIN tbl_country co ON ci.country_id = co.id
    WHERE ci.flag = 1
    AND LOWER(ci.name) LIKE LOWER($3 || '%')
    ORDER BY ci.name
    LIMIT $1 OFFSET $2
)
SELECT 
    id,
    name,
    state_id,
    state_code,
    country_id,
    country_code,
    country_name,
    latitude,
    longitude,
    created_at,
    updated_at,
    flag,
    "wikiDataId",
    total_count
FROM city_data;

`
