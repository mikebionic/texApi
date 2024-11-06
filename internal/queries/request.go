package queries

var GetMyRequest = `
SELECT 
    id, 
    user_id, 
    company_id, 
    driver_id, 
    vehicle_id, 
    cost_per_km, 
    from_country, 
    from_region, 
    to_country, 
    to_region,
    view_count,
    validity_start::varchar, 
    validity_end::varchar, 
    note, 
    created_at::varchar, 
    updated_at::varchar,
    deleted 
FROM tbl_request WHERE deleted = 0
`

var CreateMyRequest = `
INSERT INTO tbl_request (
    user_id, company_id, driver_id, vehicle_id, cost_per_km, 
    from_country, from_region, to_country, to_region, validity_start, 
    validity_end, note
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
RETURNING id;
`

var UpdateMyRequest = `
UPDATE tbl_request
SET driver_id = COALESCE($2, driver_id),
    vehicle_id = COALESCE($3, vehicle_id),
    cost_per_km = COALESCE($4, cost_per_km),
    from_country = COALESCE($5, from_country),
    from_region = COALESCE($6, from_region),
    to_country = COALESCE($7, to_country),
    to_region = COALESCE($8, to_region),
    validity_start = COALESCE($9, validity_start),
    validity_end = COALESCE($10, validity_end),
    note = COALESCE($11, note),
    updated_at = NOW()
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteMyRequest = `
UPDATE tbl_request
SET deleted = 1, updated_at = NOW()
WHERE id = $1;
`
