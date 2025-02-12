package queries

// Vehicle Brand queries
var GetVehicleBrand = `
SELECT id, name, country, founded_year, deleted
FROM tbl_vehicle_brand 
WHERE deleted = 0
`

var CreateVehicleBrand = `
INSERT INTO tbl_vehicle_brand (name, country, founded_year)
VALUES ($1, $2, $3)
RETURNING id;
`

var UpdateVehicleBrand = `
UPDATE tbl_vehicle_brand
SET name = COALESCE($2, name),
    country = COALESCE($3, country),
    founded_year = COALESCE($4, founded_year)
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteVehicleBrand = `
UPDATE tbl_vehicle_brand
SET deleted = 1
WHERE id = $1;
`

// Vehicle Type queries
var GetVehicleType = `
SELECT id,
    title_en,
	desc_en,
	title_ru,
	desc_ru,
	title_tk,
	desc_tk,
	title_de,
	desc_de,
	title_ar,
	desc_ar,
	title_es,
	desc_es,
	title_fr,
	desc_fr,
	title_zh,
	desc_zh,
	title_ja,
	desc_ja,
	deleted
FROM tbl_vehicle_type 
WHERE deleted = 0
`

var CreateVehicleType = `
INSERT INTO tbl_vehicle_type (
    title_en,
    desc_en,
    title_ru,
    desc_ru,
    title_tk,
    desc_tk,
    title_de,
    desc_de,
    title_ar,
    desc_ar,
    title_es,
    desc_es,
    title_fr,
    desc_fr,
    title_zh,
    desc_zh,
    title_ja,
    desc_ja
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING id;
`

var UpdateVehicleType = `
UPDATE tbl_vehicle_type
SET 
    title_en = COALESCE($2, title_en),
    desc_en = COALESCE($3, desc_en),
    title_ru = COALESCE($4, title_ru),
    desc_ru = COALESCE($5, desc_ru),
    title_tk = COALESCE($6, title_tk),
    desc_tk = COALESCE($7, desc_tk),
    title_de = COALESCE($8, title_de),
    desc_de = COALESCE($9, desc_de),
    title_ar = COALESCE($10, title_ar),
    desc_ar = COALESCE($11, desc_ar),
    title_es = COALESCE($12, title_es),
    desc_es = COALESCE($13, desc_es),
    title_fr = COALESCE($14, title_fr),
    desc_fr = COALESCE($15, desc_fr),
    title_zh = COALESCE($16, title_zh),
    desc_zh = COALESCE($17, desc_zh),
    title_ja = COALESCE($18, title_ja),
    desc_ja = COALESCE($19, desc_ja)
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteVehicleType = `
UPDATE tbl_vehicle_type
SET deleted = 1
WHERE id = $1;
`

// Vehicle Model queries
var GetVehicleModel = `
SELECT m.id, m.name, m.year, m.vehicle_brand_id, m.vehicle_type_id,
       b.name AS vehicle_brand,
       t.title_en AS vehicle_type, m.feature, m.deleted
FROM tbl_vehicle_model m
LEFT JOIN tbl_vehicle_type t ON t.id = m.vehicle_type_id
LEFT JOIN tbl_vehicle_brand b ON b.id = m.vehicle_brand_id
WHERE m.deleted = 0

`

var CreateVehicleModel = `
INSERT INTO tbl_vehicle_model (name, year, vehicle_brand_id, vehicle_type_id, feature)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
`

var UpdateVehicleModel = `
UPDATE tbl_vehicle_model
SET name = COALESCE($2, name),
    year = COALESCE($3, year),
    vehicle_brand_id = COALESCE($4, vehicle_brand_id),
    vehicle_type_id = COALESCE($5, vehicle_type_id),
    feature = COALESCE($6, feature)
WHERE id = $1 AND deleted = 0
RETURNING id;
`

var DeleteVehicleModel = `
UPDATE tbl_vehicle_model
SET deleted = 1
WHERE id = $1;
`

const GetVehicleList = `
WITH vehicle_data AS (
    SELECT 
        v.*,
        COUNT(*) OVER() as total_count
    FROM tbl_vehicle v
    WHERE v.deleted = 0
    ORDER BY v.id
    LIMIT $1 OFFSET $2
)
SELECT 
    vd.id, vd.uuid, vd.company_id, vd.vehicle_type_id,
    vd.vehicle_brand_id, vd.vehicle_model_id, vd.year_of_issue,
    vd.mileage, vd.numberplate, vd.trailer_numberplate,
    vd.gps, vd.photo1_url, vd.photo2_url,
    vd.photo3_url, vd.docs1_url, vd.docs2_url,
    vd.docs3_url, vd.view_count, vd.created_at,
    vd.updated_at, vd.active, vd.deleted, vd.total_count,
    vd.meta, vd.meta2, vd.meta3, vd.available,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) AS company,
    json_build_object(
        'id', vb.id,
        'name', vb.name,
        'country', vb.country,
        'founded_year', vb.founded_year
    ) AS brand,
    json_build_object(
        'id', vm.id,
        'name', vm.name,
        'year', vm.year,
        'vehicle_type', t.title_en
    ) AS model
FROM vehicle_data vd
LEFT JOIN tbl_company c ON vd.company_id = c.id
LEFT JOIN tbl_vehicle_brand vb ON vd.vehicle_brand_id = vb.id
LEFT JOIN tbl_vehicle_model vm ON vd.vehicle_model_id = vm.id
LEFT JOIN tbl_vehicle_type t ON t.id = vm.vehicle_type_id
GROUP BY 
    vd.id, vd.uuid, vd.company_id, vd.vehicle_type_id,
    vd.vehicle_brand_id, vd.vehicle_model_id, vd.year_of_issue,
    vd.mileage, vd.numberplate, vd.trailer_numberplate,
    vd.gps, vd.photo1_url, vd.photo2_url,
    vd.photo3_url, vd.docs1_url, vd.docs2_url,
    vd.docs3_url, vd.view_count, vd.created_at,
    vd.updated_at, vd.active, vd.deleted, vd.total_count,
    vd.meta, vd.meta2, vd.meta3, vd.available,
    c.id, c.company_name, c.country,
    vb.id, vb.name, vb.country, vb.founded_year,
    vm.id, vm.name, vm.year, t.title_en;


`

const CreateVehicle = `
INSERT INTO tbl_vehicle (
    company_id, vehicle_type_id, vehicle_brand_id, vehicle_model_id,
    year_of_issue, mileage, numberplate, trailer_numberplate,
    gps, photo1_url, photo2_url, photo3_url,
    docs1_url, docs2_url, docs3_url,
    view_count, meta, meta2, meta3, available
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING id;
`

const UpdateVehicle = `
UPDATE tbl_vehicle
SET 
    vehicle_type_id = COALESCE($2, vehicle_type_id),
    vehicle_brand_id = COALESCE($3, vehicle_brand_id),
    vehicle_model_id = COALESCE($4, vehicle_model_id),
    year_of_issue = COALESCE($5, year_of_issue),
    mileage = COALESCE($6, mileage),
    numberplate = COALESCE($7, numberplate),
    trailer_numberplate = COALESCE($8, trailer_numberplate),
    gps = COALESCE($9, gps),
    photo1_url = COALESCE($10, photo1_url),
    photo2_url = COALESCE($11, photo2_url),
    photo3_url = COALESCE($12, photo3_url),
    docs1_url = COALESCE($13, docs1_url),
    docs2_url = COALESCE($14, docs2_url),
    docs3_url = COALESCE($15, docs3_url),
    active = COALESCE($16, active),
    company_id = COALESCE($17, company_id),
    deleted = COALESCE($18, deleted),
    view_count = COALESCE($19, view_count),
    meta = COALESCE($20, meta),
    meta2 = COALESCE($21, meta2),
    meta3 = COALESCE($22, meta3),
    available = COALESCE($23, available),
    updated_at = NOW()`

const DeleteVehicle = `
UPDATE tbl_vehicle
SET deleted = 1, updated_at = NOW(), active = 0
WHERE id = $1;
`

const GetVehicleByID = `
SELECT 
    v.*,
    json_build_object(
        'id', c.id,
        'company_name', c.company_name,
        'country', c.country
    ) as company,
    json_build_object(
        'id', vb.id,
        'name', vb.name,
        'country', vb.country
    ) as brand,
    json_build_object(
        'id', vm.id,
        'name', vm.name,
        'year', vm.year,
        'feature', vm.feature
    ) as model,
    json_build_object(
        'id', vt.id,
        'title_en', vt.title_en,
        'title_ru', vt.title_ru,
        'title_tk', vt.title_tk
    ) as type
FROM tbl_vehicle v
LEFT JOIN tbl_company c ON v.company_id = c.id
LEFT JOIN tbl_vehicle_brand vb ON v.vehicle_brand_id = vb.id
LEFT JOIN tbl_vehicle_model vm ON v.vehicle_model_id = vm.id
LEFT JOIN tbl_vehicle_type vt ON v.vehicle_type_id = vt.id
WHERE v.id = $1 AND v.deleted = 0;
`
