package queries

const GetCargoList = `
    SELECT 
        c.*,
        COUNT(*) OVER() as total_count
    FROM tbl_cargo c
`

const GetCargoByID = `
SELECT 
    c.*
FROM tbl_cargo c
WHERE c.id = $1
`

const CreateCargo = `
INSERT INTO tbl_cargo (
    company_id, name, description, info, qty, weight, meta, meta2, meta3, 
    vehicle_type_id, packaging_type_id, gps, photo1_url, photo2_url, photo3_url, 
    docs1_url, docs2_url, docs3_url, note
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
RETURNING id;
`

const UpdateCargo = `
UPDATE tbl_cargo
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    info = COALESCE($4, info),
    qty = COALESCE($5, qty),
    weight = COALESCE($6, weight),
    meta = COALESCE($7, meta),
    meta2 = COALESCE($8, meta2),
    meta3 = COALESCE($9, meta3),
    vehicle_type_id = COALESCE($10, vehicle_type_id),
    packaging_type_id = COALESCE($11, packaging_type_id),
    gps = COALESCE($12, gps),
    photo1_url = COALESCE($13, photo1_url),
    photo2_url = COALESCE($14, photo2_url),
    photo3_url = COALESCE($15, photo3_url),
    docs1_url = COALESCE($16, docs1_url),
    docs2_url = COALESCE($17, docs2_url),
    docs3_url = COALESCE($18, docs3_url),
    note = COALESCE($19, note),
    active = COALESCE($20, active),
    deleted = COALESCE($21, deleted),
    updated_at = NOW()
WHERE id = $1;
`

const DeleteCargo = `
UPDATE tbl_cargo
SET deleted = 1, updated_at = NOW()
WHERE id = $1
`
