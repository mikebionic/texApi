package queries

var GetOrdersTotal = `
    SELECT COUNT(o.id) AS total FROM orders o
`
var GetOrders = `
    SELECT
        o.id,
        o.order_number,
        ARRAY(
            SELECT om.path FROM order_files om WHERE om.orders_id = o.id
        ) AS file_paths,
        o.users_id AS user_id,
        o.workers_id AS worker_id,
        o.address,
        o.date::VARCHAR,
        o.time::VARCHAR,
        o.time_duration,
        (
            SELECT JSON_BUILD_OBJECT(
                'id', s.id,
                'tk', GET_STATUS_TITLE(s.id, 1),
                'ru', GET_STATUS_TITLE(s.id, 2),
                'en', GET_STATUS_TITLE(s.id, 3)
            )
            FROM statuses s
            WHERE s.id = o.statuses_id
        ) AS status,
        o.description,
        ARRAY(
            SELECT
                JSON_BUILD_OBJECT( 
                    'id', os.services_id,
                    'image', (
                        SELECT s.image FROM services s 
                        WHERE s.id = os.services_id
                    ),
                    'tk', GET_SERVICE_TITLE(os.services_id, 1),
                    'ru', GET_SERVICE_TITLE(os.services_id, 2),
                    'en', GET_SERVICE_TITLE(os.services_id, 3)
                )
            FROM orders_services os
            WHERE os.orders_id = o.id
        ) AS services,
        o.read_by_admin
    FROM orders o
    ORDER BY o.id DESC
    OFFSET $1
    LIMIT $2
`
var GetNewOrders = `
    SELECT
        o.id,
        o.order_number,
        ARRAY(
            SELECT om.path FROM order_files om WHERE om.orders_id = o.id
        ) AS file_paths,
        o.users_id AS user_id,
        o.workers_id AS worker_id,
        o.address,
        o.date::VARCHAR,
        o.time::VARCHAR,
        o.time_duration,
        (
            SELECT JSON_BUILD_OBJECT(
                'id', s.id,
                'tk', GET_STATUS_TITLE(s.id, 1),
                'ru', GET_STATUS_TITLE(s.id, 2),
                'en', GET_STATUS_TITLE(s.id, 3)
            )
            FROM statuses s
            WHERE s.id = o.statuses_id
        ) AS status,
        o.description,
        ARRAY(
            SELECT
                JSON_BUILD_OBJECT( 
                    'id', os.services_id,
                    'image', (
                        SELECT s.image FROM services s 
                        WHERE s.id = os.services_id
                    ),
                    'tk', GET_SERVICE_TITLE(os.services_id, 1),
                    'ru', GET_SERVICE_TITLE(os.services_id, 2),
                    'en', GET_SERVICE_TITLE(os.services_id, 3)
                )
            FROM orders_services os
            WHERE os.orders_id = o.id
        ) AS services,
        o.read_by_admin
    FROM orders o
    WHERE o.read_by_admin = false
    ORDER BY o.id DESC
`
var GetOrdersByStatusTotal = `
    SELECT COUNT(o.id) AS total
    FROM orders o
    WHERE o.statuses_id = $1
`
var GetOrdersByStatus = `
    SELECT
        o.id,
        o.order_number,
        ARRAY(
            SELECT om.path FROM order_files om WHERE om.orders_id = o.id
        ) AS file_paths,
        o.users_id AS user_id,
        o.workers_id AS worker_id,
        o.address,
        o.date::VARCHAR,
        o.time::VARCHAR,
        o.time_duration,
        (
            SELECT JSON_BUILD_OBJECT(
                'id', s.id,
                'tk', GET_STATUS_TITLE(s.id, 1),
                'ru', GET_STATUS_TITLE(s.id, 2),
                'en', GET_STATUS_TITLE(s.id, 3)
            )
            FROM statuses s
            WHERE s.id = o.statuses_id
        ) AS status,
        o.description,
        ARRAY(
            SELECT
                JSON_BUILD_OBJECT( 
                    'id', os.services_id,
                    'image', (
                        SELECT s.image FROM services s 
                        WHERE s.id = os.services_id
                    ),
                    'tk', GET_SERVICE_TITLE(os.services_id, 1),
                    'ru', GET_SERVICE_TITLE(os.services_id, 2),
                    'en', GET_SERVICE_TITLE(os.services_id, 3)
                )
            FROM orders_services os
            WHERE os.orders_id = o.id
        ) AS services
    FROM orders o
    WHERE o.statuses_id = $1
    ORDER BY o.id DESC
    OFFSET $2
    LIMIT $3
`
var GetOrdersByWorker = `
    SELECT
        o.id,
        o.order_number,
        ARRAY(
            SELECT om.path FROM order_files om WHERE om.orders_id = o.id
        ) AS file_paths,
        (
            SELECT JSON_BUILD_OBJECT(
                'fullname', u.fullname,
                'phone', u.phone
            )
            FROM users u 
            WHERE u.id = o.users_id
        ) AS user,
        o.workers_id AS worker_id,
        o.address,
        o.date::VARCHAR,
        o.time::VARCHAR,
        o.time_duration,
        (
            SELECT JSON_BUILD_OBJECT(
                'id', s.id,
                'tk', GET_STATUS_TITLE(s.id, 1),
                'ru', GET_STATUS_TITLE(s.id, 2),
                'en', GET_STATUS_TITLE(s.id, 3)
            )
            FROM statuses s
            WHERE s.id = o.statuses_id
        ) AS status,
        o.description,
        ARRAY(
            SELECT
                JSON_BUILD_OBJECT( 
                    'id', os.services_id,
                    'tk', GET_SERVICE_TITLE(os.services_id, 1),
                    'ru', GET_SERVICE_TITLE(os.services_id, 2),
                    'en', GET_SERVICE_TITLE(os.services_id, 3)
                )
            FROM orders_services os
            WHERE os.orders_id = o.id
        ) AS services
    FROM orders o
    WHERE o.workers_id = $1
    ORDER BY o.id DESC
`
var GetOrdersByUser = `
    SELECT
        o.id,
        o.order_number,
        ARRAY(
            SELECT om.path FROM order_files om WHERE om.orders_id = o.id
        ) AS file_paths,
        (
            SELECT JSON_BUILD_OBJECT(
                'fullname', w.fullname,
                'phone', w.phone,
                'photo', w.photo,
                'about_self', w.about_self
            )
            FROM workers w 
            WHERE w.id = o.workers_id
        ) AS worker,
        o.address,
        o.date::VARCHAR,
        o.time::VARCHAR,
        o.time_duration,
        o.secret_word,
        (
            SELECT JSON_BUILD_OBJECT(
                'id', s.id,
                'tk', GET_STATUS_TITLE(s.id, 1),
                'ru', GET_STATUS_TITLE(s.id, 2),
                'en', GET_STATUS_TITLE(s.id, 3)
            )
            FROM statuses s
            WHERE s.id = o.statuses_id
        ) AS status,
        o.description,
        ARRAY(
            SELECT
                JSON_BUILD_OBJECT( 
                    'id', os.services_id,
                    'tk', GET_SERVICE_TITLE(os.services_id, 1),
                    'ru', GET_SERVICE_TITLE(os.services_id, 2),
                    'en', GET_SERVICE_TITLE(os.services_id, 3)
                )
            FROM orders_services os
            WHERE os.orders_id = o.id
        ) AS services
    FROM orders o
    WHERE o.users_id = $1
    ORDER BY o.id DESC
`
var GetOrder = `
    SELECT
        o.id,
        o.order_number,
        ARRAY(
            SELECT om.path FROM order_files om WHERE om.orders_id = o.id
        ) AS file_paths,
        o.users_id AS user_id,
        o.workers_id AS worker_id,
        o.address,
        o.date::VARCHAR,
        o.time::VARCHAR,
        o.time_duration,
        (
            SELECT JSON_BUILD_OBJECT(
                'id', s.id,
                'tk', GET_STATUS_TITLE(s.id, 1),
                'ru', GET_STATUS_TITLE(s.id, 2),
                'en', GET_STATUS_TITLE(s.id, 3)
            )
            FROM statuses s
            WHERE s.id = o.statuses_id
        ) AS status,
        o.description,
        ARRAY(
            SELECT
                JSON_BUILD_OBJECT( 
                    'id', os.services_id,
                    'image', (
                        SELECT s.image FROM services s 
                        WHERE s.id = os.services_id
                    ),
                    'tk', GET_SERVICE_TITLE(os.services_id, 1),
                    'ru', GET_SERVICE_TITLE(os.services_id, 2),
                    'en', GET_SERVICE_TITLE(os.services_id, 3)
                )
            FROM orders_services os
            WHERE os.orders_id = o.id
        ) AS services
    FROM orders o
    WHERE o.id = $1
`
var CreateOrder = `
    INSERT INTO orders (
        users_id, address, "date", "time", secret_word, description
    )
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
`
var CreateOrderServices = `
    INSERT INTO orders_services (orders_id, services_id) VALUES ($1, $2)
`
var UpdateOrder = `
    UPDATE orders SET workers_id = $1, statuses_id = $2 WHERE id = $3
    RETURNING users_id
`
var UpdateOrderStatusStart = `
    UPDATE orders SET statuses_id = 3 WHERE id = $1 AND workers_id = $2
`
var UpdateOrderTimeDuration = `
    UPDATE orders SET time_duration = $1, statuses_id = 4
    WHERE id = $2 AND workers_id = $3
`
var UpdateOrderRead = "UPDATE orders SET read_by_admin = true WHERE id = $1"
var CheckOrderExist = "SELECT id FROM orders WHERE id = $1"
var SaveOrderFile = "INSERT INTO order_files (orders_id, path) VALUES ($1, $2)"
var CheckOrderStatus = "SELECT statuses_id FROM orders WHERE id = $1"
var AbortOrder = `
    UPDATE orders SET statuses_id = 5 WHERE id = $1 AND users_id = $2
`
var DeleteOrder = `
    DELETE FROM orders o WHERE o.id = $1
    RETURNING ARRAY(
        SELECT of.path FROM order_files of WHERE of.orders_id = o.id
    )
`
