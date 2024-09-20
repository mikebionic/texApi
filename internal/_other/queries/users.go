package queries

var GetUsers = `
    SELECT
        u.id,
        u.fullname,
        u.phone,
        u.address,
        u.is_verified,
        u.subscriptions_id AS subscription_id
    FROM users u
    ORDER BY u.id DESC
    OFFSET $1
    LIMIT $2
`
var GetUser = `
    SELECT
        u.id,
        u.fullname,
        u.phone,
        u.address,
        u.is_verified,
        u.subscriptions_id AS subscription_id
    FROM users u
    WHERE u.id = $1
`
var CreateUser = `
    INSERT INTO users (
        fullname, phone, address, password, created_at, updated_at
    ) VALUES ($1, $2, $3, $4, $5, $6)
`
var CheckSubscription = `SELECT subscriptions_id FROM users WHERE id = $1`
var BuySubscription = "UPDATE users SET subscriptions_id = $1 WHERE id = $2"
var CheckUserExist = "SELECT phone FROM users WHERE phone = $1"
var CheckUserExistWithStatus = `
    SELECT phone, is_verified FROM users WHERE phone = $1
`
var UpdateUser = `
    UPDATE users SET fullname = $1, phone = $2, address = $3,
    password = $4, is_verified = $5, subscriptions_id = $6, updated_at = $7
    WHERE id = $8
`
var VerifyUser = "UPDATE users SET is_verified = true WHERE phone = $1"
var UpdateUserWithoutPassword = `
    UPDATE users SET fullname = $1, phone = $2, address = $3,
    is_verified = $4, subscriptions_id = $5, updated_at = $6
    WHERE id = $7
`
var DeleteUser = "DELETE FROM users WHERE id = $1"
var UpdateUserPassword = `
    UPDATE users u SET password = $1 WHERE u.phone = $2
    RETURNING
        u.id,
        u.fullname,
        u.phone,
        u.address,
        u.password,
        (
            SELECT
                JSON_BUILD_OBJECT(
                    'id', s.id,
                    'title', JSON_BUILD_OBJECT(
                        'tk', GET_SUBSCRIPTION_TITLE(s.id, 1),
                        'ru', GET_SUBSCRIPTION_TITLE(s.id, 2),
                        'en', GET_SUBSCRIPTION_TITLE(s.id, 3)
                    ),
                    'description', JSON_BUILD_OBJECT(
                        'tk', GET_SUBSCRIPTION_DESC(s.id, 1),
                        'ru', GET_SUBSCRIPTION_DESC(s.id, 2),
                        'en', GET_SUBSCRIPTION_DESC(s.id, 3)
                    ),
                    'start_at', s.start_at,
                    'end_at', s.end_at,
                    'days', s.days,
                    'count', s.count,
                    'price', s.price
                )
            FROM subscriptions s
            WHERE s.id = u.subscriptions_id
        ) AS subscription
`
var GetUserNotificationToken = `
    SELECT notification_token FROM users WHERE id = $1
`
var SetUserNotificationToken = `
    UPDATE users SET notification_token = $1 WHERE id = $2
`
