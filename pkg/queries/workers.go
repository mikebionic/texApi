package queries

var GetWorkers = `
    SELECT
        w.id,
        w.fullname,
        w.phone,
        w.address,
        w.photo,
        w.about_self,
        ARRAY(
            SELECT JSON_BUILD_OBJECT(
                'id', ws.services_id,
                'tk', GET_SERVICE_TITLE(ws.services_id, 1),
                'ru', GET_SERVICE_TITLE(ws.services_id, 2),
                'en', GET_SERVICE_TITLE(ws.services_id, 3)
            )
            FROM workers_services ws
            WHERE ws.workers_id = w.id
        ) AS services,
        w.created_at::VARCHAR,
        w.updated_at::VARCHAR
    FROM workers w
    ORDER BY w.id DESC
    OFFSET $1
    LIMIT $2
`
var GetWorker = `
    SELECT
        w.id,
        w.fullname,
        w.phone,
        w.address,
        w.photo,
        w.about_self,
        ARRAY(
            SELECT JSON_BUILD_OBJECT(
                'id', ws.services_id,
                'tk', GET_SERVICE_TITLE(ws.services_id, 1),
                'ru', GET_SERVICE_TITLE(ws.services_id, 2),
                'en', GET_SERVICE_TITLE(ws.services_id, 3)
            )
            FROM workers_services ws
            WHERE ws.workers_id = w.id
        ) AS services,
        w.created_at::VARCHAR,
        w.updated_at::VARCHAR
    FROM workers w
    WHERE w.id = $1
`
var CreateWorker = `
    INSERT INTO workers (
        fullname, phone, address, photo, about_self, password,
        created_at, updated_at
    ) VALUES ($1, $2, $3, NULL, $4, $5, $6, $7) RETURNING id
`
var CreateWorkerService = `
    INSERT INTO workers_services (workers_id, services_id) VALUES ($1, $2)
`
var CheckWorkerExist = "SELECT phone FROM workers WHERE phone = $1"
var UpdateWorker = `
    UPDATE workers SET fullname = $1, phone = $2, address = $3, about_self = $4,
    password = $5, updated_at = $6
    WHERE id = $7
`
var UpdateWorkerWithoutPassword = `
    UPDATE workers SET fullname = $1, phone = $2, address = $3, about_self = $4,
    updated_at = $5 WHERE id = $6
`
var UpdateWorkerService = `
    UPDATE workers_services SET services_id = $1
    WHERE workers_id = $2 AND services_id = $3
`
var DeleteWorker = "DELETE FROM workers WHERE id = $1"
var SetWorkerImage = "UPDATE workers SET photo = $1 WHERE id = $2"
