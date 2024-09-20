package queries

var GetAdmin = "SELECT id, phone, password FROM admins WHERE phone = $1"
var GetUserForLogin = `
    SELECT
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
    FROM users u WHERE u.phone = $1
`
var GetUserMe = `
    SELECT
        u.id,
        u.fullname,
        u.phone,
        u.address,
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
    FROM users u WHERE u.id = $1
`
var GetWorkerForLogin = `
    SELECT id, phone, password FROM workers WHERE phone = $1
`
