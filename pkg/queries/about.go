package queries

var GetAboutUsAll = `
    SELECT
        au.id,
        JSON_BUILD_OBJECT(
            'tk', GET_ABOUT_US_TEXT(au.id, 1),
            'ru', GET_ABOUT_US_TEXT(au.id, 2),
            'en', GET_ABOUT_US_TEXT(au.id, 3)
        ) AS text,
        au.is_active
    FROM about_us au
    ORDER BY au.id DESC
`
var GetAboutUs = `
    SELECT
        au.id,
        JSON_BUILD_OBJECT(
            'tk', GET_ABOUT_US_TEXT(au.id, 1),
            'ru', GET_ABOUT_US_TEXT(au.id, 2),
            'en', GET_ABOUT_US_TEXT(au.id, 3)
        ) AS text,
        au.is_active
    FROM about_us au
    WHERE au.id = $1
`
var GetAboutUsForUser = `
    SELECT
        au.id,
        JSON_BUILD_OBJECT(
            'tk', GET_ABOUT_US_TEXT(au.id, 1),
            'ru', GET_ABOUT_US_TEXT(au.id, 2),
            'en', GET_ABOUT_US_TEXT(au.id, 3)
        ) AS text,
        au.is_active
    FROM about_us au
    WHERE au.is_active = true
    LIMIT 1
`
var CreateAboutUs = "INSERT INTO about_us DEFAULT VALUES RETURNING id"
var CreateAboutUsTranslates = `
    INSERT INTO about_us_translates ("text", languages_id, about_us_id)
    VALUES ($1, $2, $7), ($3, $4, $7), ($5, $6, $7)
`
var UpdateAboutUsText = `
    UPDATE about_us_translates SET "text" = $1
    WHERE languages_id = $2 AND about_us_id = $3
`
var UpdateAboutUsStatus = "UPDATE about_us SET is_active = $1 WHERE id = $2"
var DeleteAboutUs = "DELETE FROM about_us WHERE id = $1"
