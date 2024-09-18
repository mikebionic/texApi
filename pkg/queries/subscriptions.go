package queries

var GetSubscriptions = `
    SELECT
        s.id,
        JSON_BUILD_OBJECT(
            'tk', GET_SUBSCRIPTION_TITLE(s.id, 1),
            'ru', GET_SUBSCRIPTION_TITLE(s.id, 2),
            'en', GET_SUBSCRIPTION_TITLE(s.id, 3)
        ) AS title,
        JSON_BUILD_OBJECT(
            'tk', GET_SUBSCRIPTION_DESC(s.id, 1),
            'ru', GET_SUBSCRIPTION_DESC(s.id, 2),
            'en', GET_SUBSCRIPTION_DESC(s.id, 3)
        ) AS description,
        s.start_at::VARCHAR,
        s.end_at::VARCHAR,
        s.days,
        s.count,
        s.price
    FROM subscriptions s
`
var GetSubscription = `
    SELECT
        s.id,
        JSON_BUILD_OBJECT(
            'tk', GET_SUBSCRIPTION_TITLE(s.id, 1),
            'ru', GET_SUBSCRIPTION_TITLE(s.id, 2),
            'en', GET_SUBSCRIPTION_TITLE(s.id, 3)
        ) AS title,
        JSON_BUILD_OBJECT(
            'tk', GET_SUBSCRIPTION_DESC(s.id, 1),
            'ru', GET_SUBSCRIPTION_DESC(s.id, 2),
            'en', GET_SUBSCRIPTION_DESC(s.id, 3)
        ) AS description,
        s.start_at::VARCHAR,
        s.end_at::VARCHAR,
        s.days,
        s.count,
        s.price
    FROM subscriptions s
    WHERE s.id = $1
`
var CreateSubscription = `
    INSERT INTO subscriptions (start_at, end_at, days, count, price)
    VALUES ($1, $2, $3, $4, $5) RETURNING id
`
var CreateSubscriptionTranslates = `
    INSERT INTO subscription_translates (
        title, description, languages_id, subscriptions_id
    ) VALUES ($1, $2, $3, $10), ($4, $5, $6, $10), ($7, $8, $9, $10)
`
var UpdateSubscription = `
    UPDATE subscriptions SET start_at = $1, end_at = $2, days = $3, count = $4,
    price = $5 WHERE id = $6
`
var UpdateSubscriptionTranslates = `
    UPDATE subscription_translates SET title = $1, description = $2
    WHERE languages_id = $3 AND subscriptions_id = $4
`
var DeleteSubscription = "DELETE FROM subscriptions WHERE id = $1"
var GetSubscriptionPrice = "SELECT price FROM subscriptions WHERE id = $1"
